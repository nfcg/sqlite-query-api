BEGIN TRANSACTION;
DROP TABLE IF EXISTS "clients";
CREATE TABLE "clients" (
	"id"	INTEGER,
	"name"	TEXT NOT NULL,
	"email"	TEXT NOT NULL UNIQUE,
	"city"	TEXT,
	"created_at"	DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id" AUTOINCREMENT)
);
DROP TABLE IF EXISTS "products";
CREATE TABLE "products" (
	"id"	INTEGER,
	"name"	TEXT NOT NULL,
	"price"	REAL NOT NULL,
	"stock"	INTEGER DEFAULT 0,
	"is_available"	BOOLEAN DEFAULT TRUE,
	PRIMARY KEY("id" AUTOINCREMENT)
);
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (1,'Alice Johnson','alice.j@example.com','New York','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (2,'Bob Smith','bob.s@example.com','Los Angeles','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (3,'Charlie Brown','charlie.b@example.com','Chicago','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (4,'David Wilson','david.w@example.com','New York','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (5,'Eve Davis','eve.d@example.com','Houston','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (6,'Frank White','frank.w@example.com','Phoenix','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (7,'Grace Taylor','grace.t@example.com','Philadelphia','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (8,'Heidi Miller','heidi.m@example.com','San Antonio','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (9,'Ivan Martinez','ivan.m@example.com','San Diego','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (10,'Judy Garcia','judy.g@example.com','Dallas','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (11,'Kevin Rodriguez','kevin.r@example.com','San Jose','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (12,'Linda Hernandez','linda.h@example.com','New York','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (13,'Mike Lopez','mike.l@example.com','Los Angeles','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (14,'Nancy Gonzalez','nancy.g@example.com','Chicago','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (15,'Oscar Perez','oscar.p@example.com','Houston','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (16,'Pamela Hall','pamela.h@example.com','Phoenix','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (17,'Quentin Moore','quentin.m@example.com','Philadelphia','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (18,'Rachel Clark','rachel.c@example.com','San Antonio','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (19,'Steve King','steve.k@example.com','San Diego','2025-08-06 08:54:22');
INSERT INTO "clients" ("id","name","email","city","created_at") VALUES (20,'Tina Young','tina.y@example.com','Dallas','2025-08-06 08:54:22');
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (1,'Laptop',1200.0,50,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (2,'Mouse',25.5,150,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (3,'Keyboard',75.0,100,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (4,'Monitor',300.0,30,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (5,'Webcam',50.0,0,0);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (6,'Headphones',99.99,75,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (7,'Speaker',150.0,25,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (8,'USB Drive 128GB',20.0,200,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (9,'External SSD',89.99,45,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (10,'Graphics Card',550.0,10,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (11,'Power Supply',85.0,60,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (12,'Motherboard',170.0,20,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (13,'RAM 16GB',65.0,90,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (14,'CPU',250.0,35,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (15,'Desk Chair',120.0,15,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (16,'Desk Lamp',35.0,80,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (17,'Webcam Stand',15.0,120,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (18,'Mousepad',10.0,300,1);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (19,'Smartwatch',199.0,40,0);
INSERT INTO "products" ("id","name","price","stock","is_available") VALUES (20,'Tablet',399.0,20,1);
COMMIT;
