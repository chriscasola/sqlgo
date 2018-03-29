package sqlgo

import (
	"database/sql"

	_ "github.com/lib/pq" // always use the postgres driver
)

// DB defines the interface of types that can be used as the database
// with Executor
type DB interface {
	Close() error
	Ping() error
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	Begin() (*sql.Tx, error)
}

// ScannerFunction defines a function that takes in pointers
// to the members of a struct to fill with the results of the
// database result scan.
type ScannerFunction func(...interface{}) error

// Deserializable defines the interface of types that can be read from
// rows in the database
type Deserializable interface {
	FromRow(ScannerFunction) error
}

// Result contains the results of a query and allows each item to
// be retrieved out into a struct
type Result struct {
	rows *sql.Rows
}

// Close closes the result object, preventing any further use
func (r *Result) Close() error {
	return r.rows.Close()
}

// Err returns the error, if any, that was encountered while enumerating
// the results
func (r *Result) Err() error {
	return r.rows.Err()
}

// Next prepares the result to allow the next result item to be read using
// the Scan method. It returns true on success, or false if there is no next
// row.
func (r *Result) Next() bool {
	return r.rows.Next()
}

// Read reads the data from the current row into the struct pointed at by
// dest.
func (r *Result) Read(dest Deserializable) error {
	return dest.FromRow(func(fields ...interface{}) error {
		return r.rows.Scan(fields...)
	})
}

// Executor connects to a SQL database and provides an API for executing
// queries
type Executor struct {
	db DB
}

// Close closes the database object
func (e *Executor) Close() error {
	return e.db.Close()
}

// Exec executes the given query that does not return a result
func (e *Executor) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := e.db.Exec(query, args...)
	return result, err
}

// Query executes the given query that is expected to return one or more
// rows and deserializes each row into the array of structs pointed to by dest
func (e *Executor) Query(query string, args ...interface{}) (*Result, error) {
	result, err := e.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return &Result{rows: result}, nil
}

// Transaction represents a database transaction
type Transaction struct {
	tx *sql.Tx
}

// Exec executes a query in the context of the transaction
func (t *Transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

// Query executes a query in the context of the transaction that returns results
func (t *Transaction) Query(query string, args ...interface{}) (*Result, error) {
	result, err := t.tx.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return &Result{rows: result}, nil
}

// Commit commits the transaction
func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

// Begin returns a new transaction
func (e *Executor) Begin() (*Transaction, error) {
	t, err := e.db.Begin()
	if err != nil {
		return nil, err
	}

	return &Transaction{tx: t}, nil
}

// NewExecutor creates a new Executor
func NewExecutor(dbURL string) (*Executor, error) {
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return &Executor{db: db}, nil
}
