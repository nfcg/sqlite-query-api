package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// ColumnInfo holds metadata about a table column.
type ColumnInfo struct {
	CID        int
	Name       string
	Type       string
	NotNull    int
	DfltValue  sql.NullString
	PK         int
}

var (
	db           *sql.DB
	dbPath       string
	serverPort   string
	tableName    string
	excludeCol   string
	sortCol      string
	sortOrder    string
	limitResults int
	showHelp     bool
)

func init() {
	// Flexible flag configuration (supports short and long options)
	flag.StringVar(&dbPath, "db", "./data.db", "Path to the SQLite database file")
	flag.StringVar(&dbPath, "d", "./data.db", "Path to the SQLite database file (short)")
	flag.StringVar(&serverPort, "port", "8080", "Port for the HTTP server")
	flag.StringVar(&serverPort, "p", "8080", "Port for the HTTP server (short)")
	flag.StringVar(&tableName, "table", "", "Name of the table to query (required)")
	flag.StringVar(&tableName, "t", "", "Name of the table to query (required, short)")
	flag.StringVar(&excludeCol, "exclude", "", "Columns to exclude (comma-separated)")
	flag.StringVar(&excludeCol, "e", "", "Columns to exclude (comma-separated, short)")
	flag.StringVar(&sortCol, "sort", "", "Column for sorting")
	flag.StringVar(&sortCol, "s", "", "Column for sorting (short)")
	flag.StringVar(&sortOrder, "order", "asc", "Sorting direction (asc|desc)")
	flag.StringVar(&sortOrder, "o", "asc", "Sorting direction (asc|desc, short)")
	flag.IntVar(&limitResults, "limit", 0, "Default number of results to return (0 for all)")
	flag.IntVar(&limitResults, "l", 0, "Default number of results to return (0 for all, short)")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Show help message (short)")
}

func main() {
	flag.Parse()

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	if tableName == "" {
		fmt.Println("Error: The table name must be specified (use --table or -t)\n")
		printHelp()
		os.Exit(1) // Exit with a non-zero status code to indicate an error
	}

	var err error
	// Resolve the absolute path for the database file to ensure cross-platform compatibility.
	absoluteDBPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatalf("Error resolving absolute path for database: %v", err)
	}

	db, err = sql.Open("sqlite3", absoluteDBPath) 
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if !tableExists(db, tableName) {
		log.Fatalf("Table '%s' not found in database '%s'", tableName, absoluteDBPath)
	}

	// Dynamic route based on the table name.
	route := "/" + tableName
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		handleDataRequest(w, r)
	})

	log.Printf("Server started on port %s", serverPort)
	log.Printf("Access at: http://localhost:%s%s", serverPort, route)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}

// tableExists checks if a table exists in the database.
func tableExists(db *sql.DB, tableName string) bool {
	var name string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&name)
	return err == nil
}

// getTableColumns retrieves metadata for all columns in a table.
func getTableColumns() ([]ColumnInfo, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return nil, fmt.Errorf("error getting table info: %v", err)
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		err := rows.Scan(&col.CID, &col.Name, &col.Type, &col.NotNull, &col.DfltValue, &col.PK)
		if err != nil {
			return nil, fmt.Errorf("error scanning column: %v", err)
		}
		columns = append(columns, col)
	}
	return columns, nil
}

// getColumnNames filters out columns based on an exclusion list.
func getColumnNames(columns []ColumnInfo, exclude string) []string {
	excludeList := parseColumns(exclude)
	var result []string

	for _, col := range columns {
		excluded := false
		for _, ex := range excludeList {
			if strings.EqualFold(col.Name, ex) {
				excluded = true
				break
			}
		}
		if !excluded {
			result = append(result, col.Name)
		}
	}

	return result
}

// parseColumns splits a comma-separated string into a slice of strings.
func parseColumns(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	var result []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// handleDataRequest processes HTTP requests for data from the table.
func handleDataRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get table columns
	columnInfo, err := getTableColumns()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get filtered column names
	columns := getColumnNames(columnInfo, excludeCol)
	if len(columns) == 0 {
		http.Error(w, "No columns available after filtering", http.StatusBadRequest)
		return
	}

	// Process sorting parameters (URL takes precedence over command line)
	urlQuery := r.URL.Query()
	sortParam := urlQuery.Get("sort")
	orderParam := urlQuery.Get("order")

	if sortParam == "" {
		sortParam = sortCol
	}
	if orderParam == "" {
		orderParam = sortOrder
	}

	if sortParam != "" {
		if !contains(columns, sortParam) {
			http.Error(w, fmt.Sprintf("Invalid sort column: %s", sortParam), http.StatusBadRequest)
			return
		}
		if orderParam != "asc" && orderParam != "desc" {
			http.Error(w, "Sorting direction must be 'asc' or 'desc'", http.StatusBadRequest)
			return
		}
	}

	// Process limit parameter (URL takes precedence over command line)
	limitParam := urlQuery.Get("limit")
	limit := limitResults // Use command line limit as default
	if limitParam != "" {
		limitValue, err := strconv.Atoi(limitParam)
		if err != nil {
			http.Error(w, "Invalid limit value. Must be an integer.", http.StatusBadRequest)
			return
		}
		limit = limitValue
	}

	// Build and execute the query
	query, args := buildQuery(urlQuery, columns, sortParam, orderParam, limit)
	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Process results
	results, err := scanRows(rows, columns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// buildQuery constructs a SQL query based on URL filters, sorting, and limit.
func buildQuery(filters map[string][]string, columns []string, sortColumn, sortDirection string, limit int) (string, []interface{}) {
	var whereClauses []string
	var args []interface{}

	// Build WHERE clauses for filters
	for col, values := range filters {
		// Ignore special parameters (sort, order, limit)
		if col == "sort" || col == "order" || col == "limit" {
			continue
		}

		if contains(columns, col) && len(values) > 0 && values[0] != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s LIKE ?", col))
			args = append(args, "%"+values[0]+"%")
		}
	}

	// Build the base query
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), tableName)

	// Add WHERE if there are filters
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Add ORDER BY if specified
	if sortColumn != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, strings.ToUpper(sortDirection))
	}

	// Add LIMIT if specified and is a positive value
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	return query, args
}

// scanRows reads all rows from the database result set and returns them as a slice of maps.
func scanRows(rows *sql.Rows, columns []string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			switch v := val.(type) {
			case []byte:
				entry[col] = string(v)
			case int64:
				entry[col] = v
			case float64:
				entry[col] = v
			case string:
				entry[col] = v
			case bool:
				entry[col] = v
			case nil:
				entry[col] = nil
			default:
				entry[col] = fmt.Sprintf("%v", v)
			}
		}
		results = append(results, entry)
	}

	return results, nil
}

// contains checks if a string is present in a slice of strings.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Println("Usage: sqlite-query-api [OPTIONS]")
	fmt.Println("\nAPI for secure querying of SQLite databases")
	fmt.Println("\nOptions:")
	fmt.Println("  -d, --db PATH       Path to the SQLite database file (default: ./data.db)")
	fmt.Println("  -p, --port PORT     Port for the HTTP server (default: 8080)")
	fmt.Println("  -t, --table TABLE   Name of the table to query (required)")
	fmt.Println("  -e, --exclude COLS  Columns to exclude (comma-separated)")
	fmt.Println("  -s, --sort COL      Column for sorting")
	fmt.Println("  -o, --order DIR     Sorting direction (asc|desc, default: asc)")
	fmt.Println("  -l, --limit N       Default number of results to return (0 for all, default: 0)")
	fmt.Println("  -h, --help          Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  sqlite-query-api -t clients -s name -o desc -l 20")
	fmt.Println("  sqlite-query-api --table products --exclude id,stock --sort price --limit 50")
	fmt.Println("\nURL access:")
	fmt.Println("  http://localhost:8080/clients?name=Pamela")
	fmt.Println("  http://localhost:8080/produtos?limit=10")
}

