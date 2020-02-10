package core

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestLogin_QueryError(t *testing.T) {
	// TODO: рассказать про разницу t.Error и t.Fatal
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = Login("", "", db)
	// errors.Is vs errors.As
	var typedErr *QueryError
	if ok := errors.As(err, &typedErr); !ok {
		t.Errorf("error not maptch QueryError: %v", err)
	}
}

func TestLogin_NoSuchLoginForEmptyDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	// Crash Early
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE managers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute query: %v", err)
	}

	result, err := Login("", "", db)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != false {
		t.Error("Login result not false for empty table")
	}
}

func TestLogin_LoginOk(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE managers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO managers(id, login, password) VALUES (1, 'vasya', 'secret')`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	result, err := Login("vasya", "secret", db)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != true {
		t.Error("Login result not true for existing account")
	}
}

func TestLogin_LoginNotOkForInvalidPassword(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	// shift 2 раза -> sql dialect
	_, err = db.Exec(`
	CREATE TABLE managers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO managers(id, login, password) VALUES (1, 'vasya', 'secret')`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = Login("vasya", "password", db)
	if !errors.Is(err, ErrInvalidPass) {
		t.Errorf("Not ErrInvalidPass error for invalid pass: %v", err)
	}
}
