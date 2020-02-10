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




const loginSQL = `SELECT login, password FROM managers WHERE login = ?`
const getAllProductsSQL = `SELECT id, name, price, qty FROM products;`
const getProductPriceAndQtyByIdSQL = `SELECT price, qty FROM products WHERE id = ?;`
const insertSaleSQL = `INSERT INTO sales(manager_id, product_id, price, qty) VALUES (:manager_id, :product_id, :price, :qty);`

const clients = `CREATE TABLE IF NOT EXISTS clients(
id INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT NOT NULL,
surname TEXT NOT NULL,
login TEXT NOT NULL UNIQUE,
password TEXT NOT NULL,
phone TEXT NOT NULL,
ban BOOLEAN NOT NULL
);`

