package api

import (
	db "simple_bank/db/sqlc"
	"simple_bank/utils"
)

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    utils.RandomString(10),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}
