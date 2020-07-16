package database

const initQuery = `
CREATE TABLE accounts
(
    id             INTEGER,
    email          TEXT,
    btc_address    TEXT,
    password       TEXT,
    balance        REAL,
    reward_points  INTEGER,
    fp_count       INTEGER,
    last_fp_time   DATE,
    stats          TEXT,
    settings       TEXT,
    browser        TEXT,
    cookies        TEXT,
    active_boosts  TEXT,
    referrer_id    INTEGER,
    login_time     DATE,
    signup_time    DATE,
    serial         TEXT,
    proxy          TEXT
);

CREATE UNIQUE INDEX id_uindex ON accounts (id);`

const insertAccountQuery = `
INSERT INTO accounts(id,email,btc_address,password,balance,reward_points,fp_count,last_fp_time,stats,settings,browser,
cookies,active_boosts,referrer_id,login_time,signup_time,serial,proxy) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

const getAllAccountsQuery = `
SELECT id,email,btc_address,password,balance,reward_points,fp_count,last_fp_time,stats,settings,browser,
cookies,active_boosts,referrer_id,login_time,signup_time,serial,proxy FROM accounts`

const updateAccountAfterRollQuery = `
UPDATE accounts SET balance=?,reward_points=?,fp_count=?,last_fp_time=?,stats=?,settings=?,cookies=?,active_boosts=?,proxy=? WHERE id=?`
