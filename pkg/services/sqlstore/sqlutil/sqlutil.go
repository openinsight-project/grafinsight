package sqlutil

import (
	"fmt"
	"os"
)

type TestDB struct {
	DriverName string
	ConnStr    string
}

func SQLite3TestDB() TestDB {
	// To run all tests in a local test database, set ConnStr to "grafinsight_test.db"
	return TestDB{
		DriverName: "sqlite3",
		// ConnStr specifies an In-memory database shared between connections.
		ConnStr: "file::memory:?cache=shared",
	}
}

func MySQLTestDB() TestDB {
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}
	return TestDB{
		DriverName: "mysql",
		ConnStr:    fmt.Sprintf("grafinsight:password@tcp(%s:%s)/grafinsight_tests?collation=utf8mb4_unicode_ci", host, port),
	}
}

func PostgresTestDB() TestDB {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	connStr := fmt.Sprintf("user=grafinsighttest password=grafinsighttest host=%s port=%s dbname=grafinsighttest sslmode=disable",
		host, port)
	return TestDB{
		DriverName: "postgres",
		ConnStr:    connStr,
	}
}

func MSSQLTestDB() TestDB {
	host := os.Getenv("MSSQL_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("MSSQL_PORT")
	if port == "" {
		port = "1433"
	}
	return TestDB{
		DriverName: "mssql",
		ConnStr:    fmt.Sprintf("server=%s;port=%s;database=grafinsighttest;user id=grafinsight;password=Password!", host, port),
	}
}
