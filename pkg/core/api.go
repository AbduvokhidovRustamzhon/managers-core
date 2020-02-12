package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	Address string
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
//func Init(db *sql.DB) (err error) {
//	ddls := []string{managersDDL, productsDDL, salesDDL, clients, atm, managers, services, cards}
//	for _, ddl := range ddls {
//		_, err = db.Exec(ddl)
//		if err != nil {
//			return err
//		}
//	}
//
//	initialData := []string{managersInitialData, productsInitialData}
//	for _, datum := range initialData {
//		_, err = db.Exec(datum)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
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



//////////////////////////////////////

func Init(db *sql.DB) (err error) {
	ddls := []string{managers, clients, atm, services, cards}
	for _, ddl := range ddls {
		_, err = db.Exec(ddl)
		if err != nil {
			return err
		}
	}

	initialData := []string{managersInitialData}
	for _, datum := range initialData {
		_, err = db.Exec(datum)
		if err != nil {
			return err
		}
	}

	return nil
}




func LoginManager(login, password string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string

	err := db.QueryRow(
		loginManagersSQL,
		login).Scan(&dbLogin, &dbPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, queryError(loginManagersSQL, err)
	}

	if dbPassword != password {
		return false, ErrInvalidPass
	}

	return true, nil
}



func LoginUsers(login, password string, db *sql.DB) (bool, error) {
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

func AddAtm( atmName string, atmAddress string, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertAtmSQL,

		sql.Named("name", atmName),
		sql.Named("adress", atmAddress),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllAtms(db *sql.DB) (atms []ATM, err error) {
	rows, err := db.Query(listAtmsSQL)
	if err != nil {
		return nil, queryError(listAtmsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			atms, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		atm := ATM{}
		err = rows.Scan(&atm.Id, &atm.Name, &atm.Address)
		if err != nil {
			return nil, dbError(err)
		}
		atms = append(atms, atm)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return atms, nil
}

func AddService( serviceName string, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertServiceSQL,

		sql.Named("name", serviceName),
		sql.Named("balance", 0),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllServices(db *sql.DB) (services []Service, err error) {
	rows, err := db.Query(listServicesSQL)
	if err != nil {
		return nil, queryError(listServicesSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			services, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		service := Service{}
		err = rows.Scan(&service.Id, &service.Name, &service.Price)
		if err != nil {
			return nil, dbError(err)
		}
		services = append(services, service)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return services, nil
}

func AddCard( cardName string, cardBalance int64, cardUser_id, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertCardsSQL,

		sql.Named("name", cardName),
		sql.Named("balance", cardBalance),
		sql.Named("user_id", cardUser_id),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllCards(db *sql.DB) (cards []Card, err error) {
	rows, err := db.Query(listCards)
	if err != nil {
		return nil, queryError(listCards, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			cards, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		card := Card{}
		err = rows.Scan(&card.Id, &card.Name, &card.Balance,&card.UserId)
		if err != nil {
			return nil, dbError(err)
		}
		cards = append(cards, card)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return cards, nil
}

func AddUser( userName string, userLogin string, userPassword string, userPassportSeries string, userPhoneNumber int, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		insertUserSQL,

		sql.Named("name", userName),
		sql.Named("login", userLogin),
		sql.Named("password", userPassword),
		sql.Named("passportSeries", userPassportSeries),
		sql.Named("phone", userPhoneNumber),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllUsers(db *sql.DB) (users []User, err error) {
	rows, err := db.Query(lisUsers)
	if err != nil {
		return nil, queryError(lisUsers, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			users, err = nil, dbError(innerErr)
		}
	}()


	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Name,&user.Login, &user.Password,&user.Phone)
		if err != nil {
			return nil, dbError(err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return users, nil
}

type Products struct {
	Id int64
	Name string
	Price int64
	Qty int64
}

func Json(db *sql.DB){
	rows, err := db.Query(getAllProductsSQL)
	if err != nil{
		log.Println(err)
	}
	defer rows.Close()

	stats := make([]*Products, 0)

	for rows.Next(){
		b:= new(Products)
		err := rows.Scan(&b.Price,&b.Qty,&b.Name,&b.Id)
		if err != nil {
			log.Fatal(err)
		}
		stats = append(stats, b)
	}
	jsonData, err := json.Marshal(&stats)
}

