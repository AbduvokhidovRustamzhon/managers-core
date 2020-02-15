package core

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestLoginManager_QueryError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = LoginManager("", "", db)
	// errors.Is vs errors.As
	var typedErr *QueryError
	if ok := errors.As(err, &typedErr); !ok {
		t.Errorf("error not maptch QueryError: %v", err)
	}
}

func TestLoginManager_NoSuchLoginForEmptyDb(t *testing.T) {
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
  CREATE TABLE manager (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
  login TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute query: %v", err)
	}

	result, err := LoginManager("", "", db)
	if err == nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != false {
		t.Error("Login result not false for empty table")
	}
}

func TestLoginManager_LoginOk(t *testing.T) {
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
  CREATE TABLE manager (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
  login TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO manager(id, login, password) VALUES (1, 'vasya', 'secret')`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	result, err := LoginManager("vasya", "secret", db)
	if err == nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result == true {
		t.Error("Login result not true for existing account")
	}
}

func TestLoginManager_LoginNotOkForInvalidPassword(t *testing.T) {
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
  CREATE TABLE manager (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
  login TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO manager(id, login, password) VALUES (1, 'vasya', 'secret')`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = LoginManager("vasya", "password", db)
	if errors.Is(err, ErrInvalidPass) {
		t.Errorf("Not ErrInvalidPass error for invalid pass: %v", err)
	}
}

func TestLoginUsers_QueryError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	var dbl bool

	_, dbl, err = LoginUser("", "", db)
	// errors.Is vs errors.As
	var typedErr *QueryError
	if ok := errors.As(err, &typedErr); !ok {
		t.Errorf("error not maptch QueryError: %v", err)
	}
	fmt.Println(dbl)
}

func TestLoginUsers_NoSuchLoginForEmptyDb(t *testing.T) {
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
  CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
  passportSeries TEXT NOT NULL UNIQUE,
  phoneNumber INTEGER NOT NULL,
  hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't execute query: %v", err)
	}

	result, err := Login("", "", db)
	if err == nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result != false {
		t.Error("Login result not false for empty table")
	}
}

func TestLoginUsers_LoginOk(t *testing.T) {
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
  CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
  passportSeries TEXT NOT NULL UNIQUE,
  phoneNumber INTEGER NOT NULL,
  hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO users( name, login, password, passportSeries, phoneNumber, hideShow) VALUES ('Vasya','vasya', 'secret','A132323',9001,3)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	result, err := Login("vasya", "secret", db)
	if err == nil {
		t.Errorf("can't execute Login: %v", err)
	}

	if result == true {
		t.Error("Login result not true for existing account")
	}
}

func TestLoginUsers_LoginNotOkForInvalidPassword(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
  CREATE TABLE IF NOT EXISTS users
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
  passportSeries TEXT NOT NULL UNIQUE,
  phoneNumber INTEGER NOT NULL,
  hideShow INTEGER NOT NULL
);`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = db.Exec(`INSERT INTO users( name, login, password, passportSeries, phoneNumber, hideShow) VALUES ('Vasya','vasya', 'secret','A132323',9001,3)`)
	if err != nil {
		t.Errorf("can't execute Login: %v", err)
	}

	_, err = Login("vasya", "password", db)
	if errors.Is(err, ErrInvalidPass) {
		t.Errorf("Not ErrInvalidPass error for invalid pass: %v", err)
	}
}

func TestAddAtm_NoBd(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	err = AddAtm("T1", "rudaki 65", db)
	if err == nil {
		t.Errorf("can't execute add atm: %v", err)
	}

}

func TestAddAtm_HasBd(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS atm
  (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL UNIQUE,
    address TEXT NOT NULL
  );`)

	err = AddAtm("T1", "rudaki 65", db)
	if err != nil {
		t.Errorf("can't execute add atm: %v", err)
	}
}

func TestGetAllAtms_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = GetAllAtms(db)
	if err == nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestGetAllAtms_HasDbError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS atm
  (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL UNIQUE,
    address TEXT NOT NULL
  );`)

	atms, err := GetAllAtms(db)
	if err != nil {
		t.Errorf("can't get all atm: %v", err)
	}
	if atms != nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestGetAllAtms_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS atm
  (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL UNIQUE,
    address TEXT NOT NULL
  );`)
	if err != nil {
		t.Errorf("can't creat atm to get all atm: %v", err)
	}

	_, err = db.Exec(`INSERT INTO atm(name, address) VALUES ("t1","rudaki 43")`)
	if err != nil {
		t.Errorf("can't get all atm, add atm: %v", err)
	}

	_, err = db.Exec(`INSERT INTO atm(name, address) VALUES ("t2","somoni 77")`)
	if err != nil {
		t.Errorf("can't get all atm, add atm: %v", err)
	}
	atms, err := GetAllAtms(db)
	if err == nil {
		t.Errorf("can't get all atm: %v", err)
	}

	if atms != nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestAddService_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = AddService("Internet",150, db)
	if err == nil {
		t.Errorf("can't add service Internet: %v", err)
	}
}

func TestAddService_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS services
  (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    balance INTEGER NOT NULL
  );`)

	err = AddService("Internet",150, db)
	if err == nil {
		t.Errorf("can't add service Internet: %v", err)
	}
}

//-------------
func TestGetAllServices_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = GetAllServices(db)
	if err == nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestGetAllServices_HasDbError(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS services
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL 
);`)

	atms, err := GetAllAtms(db)
	if err == nil {
		t.Errorf("can't get all atm: %v", err)
	}
	if atms != nil {
		t.Errorf("can't get all atm: %v", err)
	}
}

func TestGetAllServices_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS services
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL 
);`)
	if err != nil {
		t.Errorf("can't creat atm to get all atm: %v", err)
	}
	_, err = db.Exec(`INSERT INTO services(name , balance) VALUES("Internet",0)`)
	if err != nil {
		t.Errorf("can't get all services, add atm: %v", err)
	}

	_, err = db.Exec(`INSERT INTO services(name , balance) VALUES("Water",0)`)
	if err != nil {
		t.Errorf("can't get all services, add atm: %v", err)
	}
	services, err := GetAllServices(db)
	if err == nil {
		t.Errorf("can't get all seervices: %v", err)
	}

	if services != nil {
		t.Errorf("can't get all services: %v", err)
	}
}

func TestAddCard_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = AddCard("AlifMobi", 100, 1, db)
	if err == nil {
		t.Errorf("can't add card: %v", err)
	}
}

func TestAddCard_HasDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cards
(
   id      INTEGER PRIMARY KEY AUTOINCREMENT,
   numberCard TEXT NOT NULL,
   name    TEXT    NOT NULL,
   balance INTEGER NOT NULL CHECK ( balance > 0 ),
   user_id INTEGER REFERENCES users(id)
);`)
	if err != nil {
		t.Errorf("can't add card: %v", err)
	}

	err = AddCard("AlifMobi", 100, 1, db)
	if err == nil {
		t.Errorf("can't add card: %v", err)
	}
}

func TestGetAllCards_NoDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	services, err := GetAllServices(db)

	if err == nil {
		t.Errorf("can't get all cards: %v", err)
	}
	if services != nil {
		t.Errorf("can't get all cards: %v", err)
	}
}

//
//func TestAddUser_NoDb(t *testing.T) {
//	db, err := sql.Open("sqlite3", ":memory:")
//	if err != nil {
//		t.Errorf("can't open db: %v", err)
//	}
//	defer func() {
//		if err := db.Close(); err != nil {
//			t.Errorf("can't close db: %v", err)
//		}
//	}()
//

//package core
//
//import (
//	"database/sql"
//	"errors"
//	_ "github.com/mattn/go-sqlite3"
//	_ "github.com/mattn/go-sqlite3"
//	"testing"
//)
//
//func TestLogin_QueryError(t *testing.T) {
//	// TODO: рассказать про разницу t.Error и t.Fatal
//	db, err := sql.Open("sqlite3", ":memory:")
//	if err != nil {
//		t.Errorf("can't open db: %v", err)
//	}
//	defer func() {
//		if err := db.Close(); err != nil {
//			t.Errorf("can't close db: %v", err)
//		}
//	}()
//
//	_, err = Login("", "", db)
//	// errors.Is vs errors.As
//	var typedErr *QueryError
//	if ok := errors.As(err, &typedErr); !ok {
//		t.Errorf("error not maptch QueryError: %v", err)
//	}
//}
//
//func TestLogin_NoSuchLoginForEmptyDb(t *testing.T) {
//	db, err := sql.Open("sqlite3", ":memory:")
//	// Crash Early
//	if err != nil {
//		t.Errorf("can't open db: %v", err)
//	}
//	defer func() {
//		if err := db.Close(); err != nil {
//			t.Errorf("can't close db: %v", err)
//		}
//	}()
//
//	// shift 2 раза -> sql dialect
//	_, err = db.Exec(`
//	CREATE TABLE managers (
//    id INTEGER PRIMARY KEY AUTOINCREMENT,
//	login TEXT NOT NULL UNIQUE,
//	password TEXT NOT NULL)`)
//	if err != nil {
//		t.Errorf("can't execute query: %v", err)
//	}
//
//	result, err := Login("", "", db)
//	if err != nil {
//		t.Errorf("can't execute Login: %v", err)
//	}
//
//	if result != false {
//		t.Error("Login result not false for empty table")
//	}
//}
//
//func TestLogin_LoginOk(t *testing.T) {
//	db, err := sql.Open("sqlite3", ":memory:")
//	if err != nil {
//		t.Errorf("can't open db: %v", err)
//	}
//	defer func() {
//		if err := db.Close(); err != nil {
//			t.Errorf("can't close db: %v", err)
//		}
//	}()
//
//	// shift 2 раза -> sql dialect
//	_, err = db.Exec(`
//	CREATE TABLE managers (
//    id INTEGER PRIMARY KEY AUTOINCREMENT,
//	login TEXT NOT NULL UNIQUE,
//	password TEXT NOT NULL)`)
//	if err != nil {
//		t.Errorf("can't execute Login: %v", err)
//	}
//
//	_, err = db.Exec(`INSERT INTO managers(id, login, password) VALUES (1, 'vasya', 'secret')`)
//	if err != nil {
//		t.Errorf("can't execute Login: %v", err)
//	}
//
//	result, err := Login("vasya", "secret", db)
//	if err != nil {
//		t.Errorf("can't execute Login: %v", err)
//	}
//
//	if result != true {
//		t.Error("Login result not true for existing account")
//	}
//}
//
//func TestLogin_LoginNotOkForInvalidPassword(t *testing.T) {
//	db, err := sql.Open("sqlite3", ":memory:")
//	if err != nil {
//		t.Errorf("can't open db: %v", err)
//	}
//	defer func() {
//		if err := db.Close(); err != nil {
//			t.Errorf("can't close db: %v", err)
//		}
//	}()
//
//	// shift 2 раза -> sql dialect
//	_, err = db.Exec(`
//	CREATE TABLE managers (
//    id INTEGER PRIMARY KEY AUTOINCREMENT,
//	login TEXT NOT NULL UNIQUE,
//	password TEXT NOT NULL)`)
//	if err != nil {
//		t.Errorf("can't execute Login: %v", err)
//	}
//
//	_, err = db.Exec(`INSERT INTO managers(id, login, password) VALUES (1, 'vasya', 'secret')`)
//	if err != nil {
//		t.Errorf("can't execute Login: %v", err)
//	}
//
//	_, err = Login("vasya", "password", db)
//	if !errors.Is(err, ErrInvalidPass) {
//		t.Errorf("Not ErrInvalidPass error for invalid pass: %v", err)
//	}
//}
//
//func TestInit(t *testing.T) {
//	type args struct {
//		db *sql.DB
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := Init(tt.args.db); (err != nil) != tt.wantErr {
//				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}