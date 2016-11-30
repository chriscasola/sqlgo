package sqlgo

import (
	"database/sql"
	"testing"
)

type mockResult struct{}

func (r *mockResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (r *mockResult) RowsAffected() (int64, error) {
	return 1, nil
}

type mockDB struct {
	closeCount int
	lastQuery  string
}

func (db *mockDB) Close() error {
	db.closeCount++
	return nil
}

func (db *mockDB) Ping() error {
	return nil
}

func (db *mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	db.lastQuery = query
	return &mockResult{}, nil
}

func (db *mockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func TestExecutorClose(t *testing.T) {
	mockData := mockDB{}
	sut := &Executor{db: &mockData}
	sut.Close()

	if mockData.closeCount != 1 {
		t.Error("did not close")
	}
}

func TestExecutorExec(t *testing.T) {
	mockData := mockDB{}
	sut := &Executor{db: &mockData}

	sut.Exec("my query", 1, 2, 3)

	if mockData.lastQuery != "my query" {
		t.Error("did not run correct query")
	}
}
