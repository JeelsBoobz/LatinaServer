package db

import (
	"os"

	"github.com/nedpals/supabase-go"
)

type PremiumList struct {
	Id       int64  `json:"id"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Domain   string `json:"domain"`
}

func connect() *supabase.Client {
	return supabase.CreateClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))
}

func GetPremiumList() map[string][]PremiumList {
	var (
		premiumList = map[string][]PremiumList{}
		rows        = []PremiumList{}
	)

	if err := connect().DB.From("premium").Select("*").Execute(&rows); err != nil {
		panic(err)
	}

	for _, premium := range rows {
		premiumList[premium.Type] = append(premiumList[premium.Type], premium)
	}

	return premiumList
}
