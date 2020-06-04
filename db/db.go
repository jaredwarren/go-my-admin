package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Database ...
type Database struct {
	Schemas    map[string]Schema
	DSN        *DSN
	Connection *sql.DB
}

// Schema ...
type Schema struct {
	Name   string
	Tables []Table
}

// Table ...
type Table struct {
	Name string
}

// Close ...
func (d *Database) Close() error {
	return d.Connection.Close()
}

// Using current selected db
func (d *Database) Using() string {
	if d.Connection == nil {
		return ""
	}

	rs, _ := d.Query("SELECT DATABASE() as db;")
	if len(rs.Rows) > 0 && rs.Rows[0]["db"] == nil {
		return "??"
	}
	return fmt.Sprintf("%s", rs.Rows[0]["db"])
}

// Use set
func (d *Database) Use(selectedDB string) {
	if d.Connection == nil {
		return
	}
	if selectedDB == "" {
		return
	}
	// TODO: is there any harm in calling "use" if already using?
	d.Exec(fmt.Sprintf("USE %s;", selectedDB))
}

// Query run query string
func (d *Database) Query(query string, args ...interface{}) (*Result, error) {
	start := time.Now()
	rows, err := d.Connection.Query(query, args...)
	end := time.Since(start)
	if err != nil {
		return nil, err
	}
	rs := getResults(rows, 0)
	rs.Time = end

	return rs, nil
}

// Exec execute query
func (d *Database) Exec(q string) (sql.Result, error) {
	return d.Connection.Exec(q)
}

// populateStructure populate database structure (db -> schema -> tabel)
func (d *Database) populateStructure() {
	d.Schemas = map[string]Schema{}

	r, _ := d.Query("show databases;")
	for _, row := range r.Rows {
		dbName := fmt.Sprintf("%s", row["Database"])
		schema := Schema{
			Name:   dbName,
			Tables: []Table{},
		}

		// get tables
		r, err := d.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=? order by TABLE_NAME ASC;", dbName)
		if err != nil {
			fmt.Printf("%+v\n", err)
			continue
		}
		for _, t := range r.Rows {
			table := Table{fmt.Sprintf("%s", t["TABLE_NAME"])}
			schema.Tables = append(schema.Tables, table)
		}

		d.Schemas[schema.Name] = schema
	}
}

// getResults process rows in to map of interface
func getResults(rows *sql.Rows, limit int) *Result {
	if limit <= 0 {
		limit = 1000
	}
	ct := 0

	columns, _ := rows.Columns()
	colTypes, _ := rows.ColumnTypes()

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	result := []map[string]interface{}{}
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		row := make(map[string]interface{}, count)
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			var v interface{}
			if col == "password" {
				v = "**********"
			} else {
				if ok {
					v = string(b)
				} else {
					v = val
				}
			}

			row[col] = v
		}
		result = append(result, row)
		ct++
		if ct > limit {
			break
		}
	}

	return &Result{
		Rows:       result,
		Columns:    columns,
		ColumnType: colTypes,
	}
}

// FetchAll return all saved db connections
func FetchAll() (dbc map[string]*DSN, err error) {
	dbc = map[string]*DSN{}
	if !fileExists("connections.json") {
		return
	}
	// Load Connections
	data, err := ioutil.ReadFile("connections.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &dbc)
	if err != nil {
		return
	}

	return
}

// Fetch find db connection for given key
func Fetch(key string) (*Database, error) {
	dsns, err := FetchAll()
	if err != nil {
		return nil, err
	}
	for k, dsn := range dsns {
		if k == key {
			return Connect(dsn)
		}
	}
	return nil, fmt.Errorf("Unable to find connection: %s", key)
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// New return new instance of db, return connection
func New(username, password, host, port, path string) (*Database, error) {
	dsn := NewDSN(username, password, host, port, path)
	return Connect(dsn)
}

// Connect to db
func Connect(dsn *DSN) (*Database, error) {
	// check connection
	dbc, err := sql.Open("mysql", dsn.String())
	if err != nil {
		return nil, err
	}
	db := &Database{
		DSN:        dsn,
		Connection: dbc,
	}
	db.populateStructure()
	return db, nil
}

// saveConneciton save connection info to file
func saveConneciton(key string, connection *DSN) error {
	dsns, err := FetchAll()
	if err != nil {
		return err
	}

	// TODO: should I check for old connection?
	dsns[key] = connection
	file, err := json.Marshal(dsns)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("connections.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}
