package db

import (
	"fmt"
	"os"

	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/nedpals/supabase-go"
)

var domain = os.Getenv("DOMAIN")

type PremiumList struct {
	Id       int64  `json:"id"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Domain   string `json:"domain"`
	Quota    int64  `json:"quota"`
	CC       string `json:"cc"`
}

func Connect() *supabase.Client {
	return supabase.CreateClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))
}

func GetPremiumList() map[string][]PremiumList {
	var (
		premiumList = map[string][]PremiumList{}
		rows        = []PremiumList{}
	)

	if err := Connect().DB.From("premium").Select("*").Execute(&rows); err != nil {
		panic(err)
	}

	for _, premium := range rows {
		if premium.Quota > 0 {
			premiumList[premium.Type] = append(premiumList[premium.Type], premium)
		}
	}

	return premiumList
}

func UpdatePremiumQuota(name string) bool {
	rows := []PremiumList{}
	if err := Connect().DB.From("premium").Select("*").Eq("id", name).Execute(&rows); err != nil {
		fmt.Println(err)
		return true
	}

	row := rows[0]
	row.Quota = row.Quota - (helper.GetUserStats(name) / 1000000)
	if err := Connect().DB.From("premium").Update(row).Eq("id", name).Execute(&rows); err != nil {
		fmt.Println(err)
	}

	if row.Quota > 0 {
		return true
	}
	return false
}
