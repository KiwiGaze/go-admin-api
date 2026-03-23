package dto

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrderDest returns a GORM scope that orders results by the given column.
// When desc is true the order is descending, otherwise ascending.
func OrderDest(sort string, desc bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: sort}, Desc: desc})
	}
}
