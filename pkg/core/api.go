package core

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"

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
	Address string
}

type Services struct {
	Id int64
	Name string
	Price int64
}

type Client struct {
	Id int64
	Name string
	Login string
	Password string
	Balance uint64
	BalanceNumber uint64
	PhoneNumber int64
	PassportSeries int64
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


func LoginUser(login, password string, db *sql.DB) (int64 ,bool, error) {
	var dbLogin, dbPassword string
	var dbId int64
	err := db.QueryRow(
		LoginForClient,
		login).Scan(&dbId,&dbLogin, &dbPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1, false, nil
		}

		return -1,false, queryError(LoginForClient, err)
	}

	if dbPassword != password {
		return -1,false, ErrInvalidPass
	}

	return dbId ,true, nil
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
		sql.Named("address", atmAddress),
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

func AddService( serviceName string, servicePrice int64, db *sql.DB) (err error) {
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
		sql.Named("price", servicePrice),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetAllServices(db *sql.DB) (services []Services, err error) {
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
		service := Services{}
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

func AddCard( cardName string, cardBalance int64, cardUserId int64, db *sql.DB) (err error) {
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
		sql.Named("user_id", cardUserId),
	)
	if err != nil {
		return err
	}

	return nil
}

func AddUser( userName string, userLogin string, userPassword string, userPassportSeries string, userPhoneNumber int, balance uint64, balanceNumber int64, db *sql.DB) (err error) {
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
		sql.Named("passport_series", userPassportSeries),
		sql.Named("phone", userPhoneNumber),
		sql.Named("balance", balance),
		sql.Named("balance_number", balanceNumber),
	)
	if err != nil {
		return err
	}

	return nil
}



func UpdateBalanceClient(id int64, balance int64,  db *sql.DB) (err error) {
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
		updateCardBalanceSQL,
		sql.Named("id", id),
		sql.Named("balance", balance),
	)
	if err != nil {
		return err
	}

	return nil
}


func UpdateBalanceClientForService(login string, balance int64,  db *sql.DB) (err error) {
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
		updateClientBalanceMinusSQL,
		sql.Named("login", login),
		sql.Named("balance", balance),
	)
	if err != nil {
		return err
	}

	return nil
}




func PayForService(id int64, balance int64,  db *sql.DB) (err error) {
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
		payForServices,
		sql.Named("id", id),
		sql.Named("price", balance),
	)
	if err != nil {
		return err
	}

	return nil
}

func GetBalanceList(db *sql.DB, userId int64) (listBalance []Client, err error) {
	rows, err := db.Query(lisUsers, userId)
	if err != nil {
		return nil, queryError(lisUsers, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			listBalance, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		listAccount := Client{}
		err = rows.Scan(&listAccount.Id, &listAccount.Name, &listAccount.BalanceNumber, &listAccount.Balance )
		if err != nil {
			return nil, dbError(err)
		}
		listBalance = append(listBalance,listAccount)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return listBalance, nil
}

func TransactionBalanceNumberMinus(tranzaction Client, db *sql.DB) (err error) {
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
		updateTransactionWithBalanceNumberMinus,
		sql.Named("balance_number", tranzaction.BalanceNumber),
		sql.Named("balance", tranzaction.Balance),
	)
	if err != nil {
		return err
	}

	return nil
}

func CheckByBalanceNumber(balanceNumber uint64, db *sql.DB)(err error)  {
	var id int
	err = db.QueryRow("select id from client where balance_number=?", balanceNumber).Scan(&id)
	return err
}

func CheckByPhoneNumber(phoneNumber int64,db *sql.DB) (err error) {
	var id int
	err = db.QueryRow("select id from client where phone_number=?", phoneNumber).Scan(&id)
	return err
}

func TransactionPlus(phoneNumber int64,balance uint64, db *sql.DB) (err error) {
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
		updateTransactionWithPhoneNumberPlus,
		sql.Named("phone_number", phoneNumber),
		sql.Named("balance", balance),
	)
	if err != nil {
		return err
	}

	return nil
}

func TransactionMinus(tranzaction Client, db *sql.DB) (err error) {
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
		updateTransactionWithPhoneNumberMinus,
		sql.Named("phone_number", tranzaction.PhoneNumber),
		sql.Named("balance", tranzaction.Balance),
	)
	if err != nil {
		return err
	}

	return nil
}

func TransactionBalanceNumberPlus(balanceNumber uint64,balance uint64, db *sql.DB) (err error) {
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
		updateTransactionWithBalanceNumberPlus,
		sql.Named("balance_number", balanceNumber),
		sql.Named("balance", balance),
	)
	if err != nil {
		return err
	}

	return nil
}

// export

func ExportClientsToJSON(db *sql.DB) error {
	return ExportToFile(db, getAllClientsDataSQL, "clients.json",
		mapRowToClient, json.Marshal, mapInterfaceSliceToClients)
}
func ExportAtmsToJSON(db *sql.DB) error {
	return ExportToFile(db, getAllAtmDataSQL, "atms.json",
		mapRowToAtm, json.Marshal,
		mapInterfaceSliceToAtms)
}

//XML

func ExportClientsToXML(db *sql.DB) error {
	return ExportToFile(db, getAllClientsDataSQL, "clients.xml",
		mapRowToClient, xml.Marshal, mapInterfaceSliceToClients)
}
func ExportAtmsToXML(db *sql.DB) error {
	return ExportToFile(db, getAllAtmDataSQL, "atms.xml",
		mapRowToAtm, xml.Marshal,
		mapInterfaceSliceToAtms)
}

func mapRowToClient(rows *sql.Rows) (interface{}, error) {
	client := Client{}
	err := rows.Scan(&client.Id, &client.Login, &client.Password,
		&client.Name, &client.PhoneNumber, &client.Balance, &client.BalanceNumber, &client.PassportSeries)
	if err != nil {
		return nil, err
	}
	return client, nil
}
func mapRowToAtm(rows *sql.Rows) (interface{}, error) {
	atm := ATM{}
	err := rows.Scan(&atm.Id,&atm.Name, &atm.Address)
	if err != nil {
		return nil, err
	}
	return atm, nil
}
type ClientsExport struct {
	Clients []Client
}
func mapInterfaceSliceToClients(ifaces []interface{}) interface{} {
	clients := make([]Client, len(ifaces))
	for i := range ifaces {
		clients[i] = ifaces[i].(Client)
	}
	clientsExport := ClientsExport{Clients: clients}
	return clientsExport
}
func mapInterfaceSliceToAtms(ifaces []interface{}) interface{} {
	atms := make([]ATM, len(ifaces))
	for i := range ifaces {
		atms[i] = ifaces[i].(ATM)
	}
	atmsExport := AtmsExport{Atms: atms}
	return atmsExport
}
func ImportClientsFromJSON(db *sql.DB) error {
	return ImportFromFile(
		db,
		"clients.json",
		func(data []byte) ([]interface{}, error) {
			return mapBytesToClients(data, json.Unmarshal)
		},
		insertClientToDB,
	)
}
func ImportAtmsFromJSON(db *sql.DB) error {
	return ImportFromFile(
		db,
		"atms.json",
		func(data []byte) ([]interface{}, error) {
			return mapBytesToAtms(data, json.Unmarshal)
		},
		insertAtmToDB,
	)
}
func ImportClientsFromXML(db *sql.DB) error {
	return ImportFromFile(
		db,
		"clients.xml",
		func(data []byte) ([]interface{}, error) {
			return mapBytesToClients(data, xml.Unmarshal)
		},
		insertClientToDB,
	)
}
func ImportAtmsFromXML(db *sql.DB) error {
	return ImportFromFile(
		db,
		"atms.xml",
		func(data []byte) ([]interface{}, error) {
			return mapBytesToAtms(data, xml.Unmarshal)
		},
		insertAtmToDB,
	)
}
func mapBytesToClients(data []byte,
	unmarshal func([]byte, interface{}) error,
) ([]interface{}, error) {
	clientsExport := ClientsExport{}
	err := unmarshal(data, &clientsExport)
	if err != nil {
		return nil, err
	}
	ifaces := make([]interface{}, len(clientsExport.Clients))
	for index := range ifaces {
		ifaces[index] = clientsExport.Clients[index]
	}
	return ifaces, nil
}
func insertClientToDB(iface interface{}, db *sql.DB) error {
	client := iface.(Client)
	_, err := db.Exec(
		insertClientSQL,
		sql.Named("id", client.Id),
		sql.Named("name", client.Name),
		sql.Named("login", client.Login),
		sql.Named("password", client.Password),
		sql.Named("phone", client.PhoneNumber),
		sql.Named("balance_number", client.BalanceNumber),
		sql.Named("balance", client.Balance),
	)
	if err != nil {
		return err
	}
	return nil
}

type AtmsExport struct {
	Atms []ATM
}

func mapBytesToAtms(data []byte,
	unmarshal func([]byte, interface{}) error,
) ([]interface{}, error) {
	atmsExport := AtmsExport{}
	err := unmarshal(data, &atmsExport)
	if err != nil {
		return nil, err
	}
	ifaces := make([]interface{}, len(atmsExport.Atms))
	for index := range ifaces {
		ifaces[index] = atmsExport.Atms[index]
	}
	return ifaces, nil
}
func insertAtmToDB(iface interface{}, db *sql.DB) error {
	atm := iface.(ATM)
	_, err := db.Exec(
		insertAtmSQL,
		sql.Named("id", atm.Id),
		sql.Named("name", atm.Name),
		sql.Named("address", atm.Address),
	)
	if err != nil {
		return err
	}
	return nil
}


type MapperRowTo func(rows *sql.Rows) (interface{}, error)
type MapperInterfaceSliceTo func([]interface{}) interface{}
type Marshaller func(interface{}) ([]byte, error)

func ExportToFile(
	db *sql.DB,
	getDataFromDbSQL string,
	filename string,
	mapRow MapperRowTo,
	marshal Marshaller,
	mapDataSlice MapperInterfaceSliceTo) error {

	rows, err := db.Query(getDataFromDbSQL)
	if err != nil {
		return err
	}
	defer func() {
		err = rows.Close()
	}()
	var dataSlice []interface{}
	for rows.Next() {
		dataElement, err := mapRow(rows)
		if err != nil {
			return err
		}
		dataSlice = append(dataSlice, dataElement)
	}
	exportData := mapDataSlice(dataSlice)
	data, err := marshal(exportData)
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}
	return nil
}


type MapperBytesTo func([]byte) ([]interface{}, error)

func ImportFromFile(
	db *sql.DB,
	filename string,
	mapBytes MapperBytesTo,
	insertToDB func(interface{}, *sql.DB) error,
) error {
	itemsData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	sliceData, err := mapBytes(itemsData)

	for _, datum := range sliceData {
		err = insertToDB(datum, db)
		if err != nil {
			return err
		}
	}

	return nil
}
