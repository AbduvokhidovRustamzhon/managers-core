package core

const managersDDL = `
CREATE TABLE IF NOT EXISTS managers
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL,
    salary  INTEGER NOT NULL CHECK ( salary > 0 ),
    plan    INTEGER NOT NULL CHECK ( plan >= 0 ),
    unit 	TEXT,
    boss_id INTEGER REFERENCES managers
);`

const productsDDL = `
CREATE TABLE IF NOT EXISTS products
(
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    name  TEXT    NOT NULL UNIQUE,
    price INTEGER NOT NULL CHECK ( price > 0 ),
    qty INTEGER NOT NULL CHECK ( qty > 0 )
);`

const salesDDL = `
CREATE TABLE IF NOT EXISTS sales (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    manager_id INTEGER NOT NULL REFERENCES managers,
    product_id INTEGER NOT NULL REFERENCES products,
    qty INTEGER NOT NULL CHECK ( qty > 0 ),
    price INTEGER NOT NULL CHECK ( price > 0 )
);`

const managersInitialData = `INSERT INTO managers
VALUES (1, 'Vasya', 'vasya', 'secret', 100000, 0, NULL, NULL),
       (2, 'Petya', 'petya', 'secret', 90000, 90000, 'boys', 1),
       (3, 'Vanya', 'vanya', 'secret', 80000, 80000, 'boys', 2),
       (4, 'Masha', 'masha', 'secret', 80000, 80000, 'girls', 1),
       (5, 'Dasha', 'dasha', 'secret', 60000, 60000, 'girls', 4),
       (6, 'Sasha', 'sasha', 'secret', 40000, 40000, 'girls', 5)
       ON CONFLICT DO NOTHING;`

const productsInitialData = `INSERT INTO products(id, name, price, qty)
VALUES (1, 'Big Mac', 200, 10),       -- 1
       (2, 'Chicken Mac', 150, 15),   -- 2
       (3, 'Cheese Burger', 100, 20), -- 3
       (4, 'Tea', 50, 10),            -- 4
       (5, 'Coffee', 80, 10),         -- 5
       (6, 'Cola', 100, 20)           -- 6
       ON CONFLICT DO NOTHING;`




const loginUserSQL = `SELECT login, password FROM managers WHERE login = ?;`
const getAllProductsSQL = `SELECT id, name, price, qty FROM products;`
const getProductPriceAndQtyByIdSQL = `SELECT price, qty FROM products WHERE id = ?;`
const insertSaleSQL = `INSERT INTO sales(manager_id, product_id, price, qty) VALUES (:manager_id, :product_id, :price, :qty);`

const loginManagersSQL  = `SELECT login, password FROM managers WHERE login = ?;`
const listAtmsSQL = `SELECT name, address FROM atm;`
const listServicesSQL  = `SELECT id, name, price FROM service;`
const listCards = ` SELECT id, name, balance, user_id FROM card;`
const lisUsers = `SELECT id, name, surname, middle_name, sex, email, login, password, phone, address, ban FROM client;`



const clients = `CREATE TABLE IF NOT EXISTS client(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	-- surname TEXT NOT NULL,
	-- middle_name TEXT NOT NULL,
	-- sex TEXT NOT NULL,
	-- email TEXT NOT NULL UNIQUE,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	passport_series TEXT NOT NULL UNIQUE,
	phone TEXT NOT NULL UNIQUE,
	balance INTEGER NOT NULL,
	balance_number INTEGER NOT NULL UNIQUE
	-- address TEXT NOT NULL,
	-- ban BOOLEAN NOT NULL,
);`

const managers = `CREATE TABLE IF NOT EXISTS manager(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	login TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	passport_series TEXT NOT NULL UNIQUE,
	phone TEXT NOT NULL
);`

const atm  = `CREATE TABLE IF NOT EXISTS atm(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	address TEXT NOT NULL
);`

const services  = `CREATE TABLE IF NOT EXISTS service(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	price INTEGER NOT NULL CHECK(price > 0)
);`


const cards = `CREATE TABLE IF NOT EXISTS card(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	balance INTEGER NOT NULL,
	user_id TEXT NOT NULL REFERENCES client
);`
// -- Insertes
const insertAtmSQL = `INSERT INTO atm( name, address)VALUES( :name,:address);`
const insertServiceSQL = `INSERT INTO service( name, price)VALUES( :name, :price);`
const insertCardsSQL = `INSERT INTO card(name, balance, user_id)VALUES( :name, :balance);`
const insertUserSQL = `INSERT INTO client(name, login, password, passport_series, phone, balance, balance_number)VALUES( :name, :login, :password, :passport_series, :phone, :balance, :balance_number )`
const getAllAtmDataSQL = `SELECT * FROM atm;`
const getAllClientsDataSQL = `SELECT * FROM client`
// -- Updates
const updateCardBalanceSQL = `UPDATE client SET balance=balance + :balance WHERE id = :id;`
const updateClientBalanceMinusSQL = `UPDATE client SET balance = balance - :balance WHERE login=:login;`
const updateClientBalancePlusSQL =	`UPDATE client SET balance = balance + :balance WHERE id = :id;`

const priceOfService = `SELECT price FROM service WHERE id= :id;`
const payForServices = `UPDATE service SET price = price + :price WHERE id =:id;`
const updateTransactionWithPhoneNumberMinus = `UPDATE client SET balance = balance - :balance WHERE phone_number = :phone_number;`
const updateTransactionWithPhoneNumberPlus = `UPDATE client SET balance = balance + :balance where phone_number = :phone_number;`
const updateTransactionWithBalanceNumberMinus = `UPDATE client SET balance = balance - :balance WHERE balance_number = :balance_number;`
const updateTransactionWithBalanceNumberPlus = `UPDATE client SET balance = balance + :balance where balance_number = :balance_number;`
const LoginForClient = `select id, login,password from client where login = ?;`
const getAllAtmSql = `select id,name,street from atm;`

const insertClientSQL  = `insert into client(id,name,login,password,passport_series,phone) values(:id, :name, :login, :password,:passport_series,:phone) ON CONFLICT DO NOTHING;`