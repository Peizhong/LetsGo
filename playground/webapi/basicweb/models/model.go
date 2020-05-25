package models

type MoneyAccount struct {
	Id          int    `gorm:"Column:id"`
	AccountName string `gorm:"Column:account_name;type:varchar(200);"`
}
