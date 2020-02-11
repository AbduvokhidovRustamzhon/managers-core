package core

import (
	"database/sql"
	"errors"
	"fmt"
)

// ошибки - это тоже часть API
var ErrInvalidPass = errors.New("invalid password")

type QueryError struct { // alt + enter
	Query string
	Err   error
}

type DbError struct {
	Err error
}

type DbTxError struct {
	Err         error
	RollbackErr error
}

type Product struct {
	Id    int64
	Name  string
	Price int64
	Qty   int64
}

type ATM struct {
	Id int64
	Name string
	AddressU string
}

type Service struct {
	Id int64
	Name string
	Price int64
}

type Card struct {
	Id int64
	Name string
	Balance int64
	UserId int64
}

type User struct {
	Id int64
	Name string
	Surname string
	MiddleName string
	Login string
	Email string
	Password string
	Phone int64
	Ban bool
}

func (receiver *QueryError) Unwrap() error {
	return receiver.Err
}

func (receiver *QueryError) Error() string {
	return fmt.Sprintf("can't execute query %s: %s", loginUserSQL, receiver.Err.Error())
}

func queryError(query string, err error) *QueryError {
	return &QueryError{Query: query, Err: err}
}

func (receiver *DbError) Error() string {
	return fmt.Sprintf("can't handle db operation: %v", receiver.Err.Error())
}

func (receiver *DbError) Unwrap() error {
	return receiver.Err
}

func dbError(err error) *DbError {
	return &DbError{Err: err}
}

// TODO: INIT
func Init(db *sql.DB) (err error) {
	ddls := []string{managersDDL, productsDDL, salesDDL, clients, atm, managers, services, cards}
	for _, ddl := range ddls {
		_, err = db.Exec(ddl)
		if err != nil {
			return err
		}
	}

	initialData := []string{managersInitialData, productsInitialData}
	for _, datum := range initialData {
		_, err = db.Exec(datum)
		if err != nil {
			return err
		}
	}

	return nil
}
// TODO: INIT


func Login(login, password string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string

	err := db.QueryRow(
		loginUserSQL,
		login).Scan(&dbLogin, &dbPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, queryError(loginUserSQL, err)
	}

	if dbPassword != password {
		return false, ErrInvalidPass
	}

	return true, nil
}

//func queryData(db *sql.DB, query string, mapRow func(rows *sql.Rows)) {
//	// mapping -> отображение одних данных в другие
//	// map
//}


func GetAllProducts(db *sql.DB) (products []Product, err error) {
	rows, err := db.Query(getAllProductsSQL)
	if err != nil {
		return nil, queryError(getAllProductsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			products, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		product := Product{}
		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Qty)
		if err != nil {
			return nil, dbError(err)
		}
		products = append(products, product)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return products, nil
}

// TODO: add manager_id
func Sale(productId int64, productQty int64, db *sql.DB) (err error) {
	// begin + commit|rollback
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// если произошла ошибка, пробуем делать rollback
	// если нет, делаем commit (вернёт ошибку или nil)
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var currentPrice int64
	var currentQty int64

	err = tx.QueryRow(
		getProductPriceAndQtyByIdSQL,
		productId,
	).Scan(&currentPrice, &currentQty)
	if err != nil {
		return err
	}

	// TODO: желательно ещё проверить, что продаём меньше или равно
	_, err = tx.Exec(
		insertSaleSQL,
		sql.Named("manager_id", 1),
		sql.Named("product_id", productId),
		sql.Named("price", currentPrice),
		sql.Named("qty", productQty),
	)
	if err != nil {
		return err
	}

	return nil
}
