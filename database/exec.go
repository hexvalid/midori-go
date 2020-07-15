package database

import (
	"database/sql"
	"encoding/json"
	"github.com/fatih/color"
	"github.com/hexvalid/midori-go/bot"
	"strconv"
)

func InsertAccount(db *sql.DB, a *bot.Account) error {

	statsJson, err := json.Marshal(&a.Stats)
	if err != nil {
		return err
	}
	browserJson, err := json.Marshal(&a.Browser)
	if err != nil {
		return err
	}
	settingsJson, err := json.Marshal(&a.Settings)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(insertAccountQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		a.ID,
		a.Email,
		a.BTCAddress,
		a.Password,
		a.Balance,
		a.RewardPoints,
		a.FPCount,
		a.LastFPDate,
		statsJson,
		settingsJson,
		browserJson,
		a.JarToString(),
		"",
		a.ReferrerID,
		a.LoginTime,
		a.SignUpTime,
		a.Serial,
		a.Proxy)
	if err == nil {
		log.Info("Account successfully inserted to database: %s", color.BlueString(strconv.Itoa(a.ID)))
	}
	return err
}

func GetAllAccounts(db *sql.DB) (accs []*bot.Account, err error) {

	stmt, err := db.Prepare(getAllAccountsQuery)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	for rows.Next() {

		var statsString string
		var settingsString string
		var browserString string
		var cookiesString string
		var activeBoostsString string

		a := bot.Account{}
		if err = rows.Scan(
			&a.ID,
			&a.Email,
			&a.BTCAddress,
			&a.Password,
			&a.Balance,
			&a.RewardPoints,
			&a.FPCount,
			&a.LastFPDate,
			&statsString,
			&settingsString,
			&browserString,
			&cookiesString,
			&activeBoostsString,
			&a.ReferrerID,
			&a.LoginTime,
			&a.SignUpTime,
			&a.Serial,
			&a.Proxy,
		); err != nil {
			return
		}

		if err = json.Unmarshal([]byte(statsString), &a.Stats); err != nil {
			return
		}

		if err = json.Unmarshal([]byte(settingsString), &a.Settings); err != nil {
			return
		}

		if err = json.Unmarshal([]byte(browserString), &a.Browser); err != nil {
			return
		}

		a.StringToJar(cookiesString)

		accs = append(accs, &a)
	}
	log.Info("%s account loaded from database", color.YellowString(strconv.Itoa(len(accs))))
	return accs, err
}
