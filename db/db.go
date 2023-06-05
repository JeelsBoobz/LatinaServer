package db

import (
	"os"

	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/nedpals/supabase-go"
)

var PendingUpdateQuotas = map[string]int64{}

type PremiumList struct {
	Id       int64  `json:"id"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Domain   string `json:"domain"`
	Quota    int64  `json:"quota"`
}

func connect() *supabase.Client {
	return supabase.CreateClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))
}

func GetPremiumList() map[string][]PremiumList {
	var (
		premiumList = map[string][]PremiumList{}
		rows        = []PremiumList{}
	)

	if err := connect().DB.From("premium").Select("*").Eq("domain", os.Getenv("DOMAIN")).Execute(&rows); err != nil {
		panic(err)
	}

	for _, premium := range rows {
		premiumList[premium.Type] = append(premiumList[premium.Type], premium)
	}

	return premiumList
}

func UpdatePremiumQuota(name string) {
	var (
		curUsage = PendingUpdateQuotas[name] + helper.GetUserStats(name)/1000000 // In MB
	)

	row := PremiumList{}

	if err := connect().DB.From("premium").Update(PremiumList{
		Quota: curUsage,
	}).Eq("id", name).Execute(&row); err != nil {
		PendingUpdateQuotas[name] = curUsage
	} else {
		PendingUpdateQuotas[name] = 0
	}
}
