
# SQLite Query API

The **SQLite Query API** is a lightweight and flexible Go application that exposes SQLite database tables via a simple HTTP API. It allows you to query data, apply filters, sort results, and limit the number of records returned, all through URL parameters. This application is designed for quick setup and easy integration, providing a convenient way to access your SQLite data without complex configurations.

---

## ✨ Features

* **Dynamic Table Routing**: Access data from any specified table directly via its name in the URL (e.g., `/clients`, `/products`).
* **Flexible Filtering**: Filter records based on column values using URL query parameters.
* **Sorting Capabilities**: Order results by any column in ascending or descending order.
* **Result Limiting**: Control the number of records returned per query, either via command-line default or URL parameter.
* **Column Exclusion**: Exclude specific columns from the query results.
* **Easy to Deploy**: A single Go binary, ideal for direct server deployment.

---

## 🚀 Getting Started

These instructions will get you a copy of the project up and running.

### Prerequisites

Ensure you have the following installed:

* **Go (1.16 or higher)**: [https://golang.org/doc/install](https://golang.org/doc/install)
* **SQLite3**: The database itself. Most Linux distributions have it pre-installed. For Windows/macOS, you might need to install it.

### Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/nfcg/sqlite-query-api.git
    cd sqlite-query-api
    ```
2.  **Build the application**:
    ```bash
    go build sqlite-query-api.go
    ```
    This will create an executable named `sqlite-query-api` in your current directory.

### Database Setup (Example)

You can import the sample file data.sql with `clients` and `products` tables using the following bash command. 

```bash
sqlite3 'data.db' < 'data.sql'
````

### Usage

Run the application from your terminal:

```bash
./sqlite-query-api -d data.db -t clients
````

This will start the server on `http://localhost:8080` and expose the `clients` table at `http://localhost:8080/clients`.

-----

## ⚙️ Configuration

The application can be configured using command-line flags:

  * `-d, --db PATH`: Path to the SQLite database file (default: `./data.db`)
  * `-p, --port PORT`: Port for the HTTP server (default: `8080`)
  * `-t, --table TABLE`: **(Required)** Name of the table to be queried.
  * `-e, --exclude COLS`: Comma-separated list of columns to exclude from results.
  * `-s, --sort COL`: Column name for default sorting.
  * `-o, --order DIR`: Default sorting direction (`asc` or `desc`, default: `asc`).
  * `-l, --limit N`: Default number of results to return (0 for all, default: `0`).
  * `-h, --help`: Show the help message.

**Examples:**

```bash
# Start server for 'clients' table, sorted by 'name' descending, with a default limit of 20
./sqlite-query-api -d data.db -t clients -s name -o desc -l 20

# Start server for 'products' table, excluding 'id' and 'price'
./sqlite-query-api --db data.db --table products --exclude id,price
```

-----

## 🌐 API Endpoints

The API exposes a single endpoint based on the table name you specify.

**Base URL**: `http://localhost:8080/<YOUR_TABLE_NAME>` (or your Nginx/Apache proxy URL)

### Query Parameters

You can customize your queries using URL parameters:

  * **Filtering**: Add any column name as a query parameter to filter results. The value will be matched using `%value%`.
      * Example: `/clients?city=New York`
  * **`sort`**: Column name to sort the results by. Overrides the command-line `--sort` flag.
      * Example: `/clients?sort=name`
  * **`order`**: Sorting direction (`asc` or `desc`). Overrides the command-line `--order` flag.
      * Example: `/clients?sort=name&order=desc`
  * **`limit`**: Maximum number of records to return. Overrides the command-line `--limit` flag.
      * Example: `/clients?limit=5`

### Combined Example

```
http://localhost:8080/clients?city=New%20York&sort=name&order=asc&limit=10
```

This query would:

1.  Access the `clients` table.
2.  Filter records where the `city` column contains "New York".
3.  Sort the results by `name` in ascending order.
4.  Limit the output to the first 10 records.

-----

## 🚀 Deployment with Reverse Proxy

For production environments, it's recommended to use a reverse proxy like Nginx or Apache to serve your application. This allows you to handle SSL/TLS, load balancing, and serve static files more efficiently.

### Nginx Configuration Example

Add this block to your Nginx configuration (e.g., in `/etc/nginx/sites-available/sqlite-query-api.conf`):

```nginx
server {
    listen 80;
    server_name example.com; # Replace with your domain

    location / {
        proxy_pass http://localhost:8080/; # Your Go app's address
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Apache Configuration Example

Ensure `mod_proxy`, `mod_proxy_http`, `mod_proxy_html`, and `mod_headers` are enabled:


Add this block to your Apache VirtualHost configuration (e.g., in `/etc/apache2/sites-available/sqlite-query-api.conf`):

```apache
<VirtualHost *:80>
    ServerName example.com # Replace with your domain

    ProxyRequests Off
    ProxyPreserveHost On

    ProxyPass / http://localhost:8080/
    ProxyPassReverse / http://localhost:8080/

    RequestHeader set X-Real-IP "%{REMOTE_ADDR}s"
    RequestHeader set X-Forwarded-For "%{REMOTE_ADDR}s"
    RequestHeader set X-Forwarded-Proto "http"
</VirtualHost>
```

-----

## 🤝 Contributing

Contributions are welcome\! If you have suggestions for improvements or new features, please open an issue or submit a pull request.

-----




