package models

import (
	"regexp"
	"unicode"
)

//User struct for database entries
type Payment struct {
	Id        int64
	UserID        int64
	Created   string
	InvoiceId string
	Amount int
	PackageID int64
	Method_Payment string
	Status int
	InforReturn string
}
