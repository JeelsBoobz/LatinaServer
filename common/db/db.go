package db

import (
	"database/sql"
	"strconv"

	"github.com/LalatinaHub/LatinaSub-go/db"
)

type PremiumList struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

func GetPremiumList() map[string][]PremiumList {
	var (
		premiumList       = map[string][]PremiumList{}
		query             = "SELECT * FROM premium;"
		name              sql.NullInt64
		password, vpnType sql.NullString
	)

	rows, err := db.New().Conn().Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &password, &vpnType)
		if err != nil {
			panic(err)
		}

		premium := PremiumList{
			Name:     strconv.Itoa(int(name.Int64)),
			Password: password.String,
			Type:     vpnType.String,
		}

		premiumList[premium.Type] = append(premiumList[premium.Type], premium)
	}

	return premiumList
}
