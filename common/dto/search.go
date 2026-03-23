// Package dto provides shared data transfer objects and GORM scope builders
// for search, pagination, and common CRUD operations across app modules.
package dto

import (
	"go-admin-api/common/global"

	"github.com/go-admin-team/go-admin-core/tools/search"
	"gorm.io/gorm"
)

// GeneralDelDto carries single and batch IDs for delete operations.
type GeneralDelDto struct {
	Id  int   `uri:"id" json:"id" validate:"required"`
	Ids []int `json:"ids"`
}

// GetIds returns the normalized list of valid IDs for delete operations.
// It includes the single Id field and any positive values in Ids.
// When no valid IDs are provided, it returns []int{0}.
func (g GeneralDelDto) GetIds() []int {
	ids := make([]int, 0, len(g.Ids)+1)

	if g.Id > 0 {
		ids = append(ids, g.Id)
	}

	for _, id := range g.Ids {
		if id > 0 {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return []int{0}
	}

	return ids
}

// GeneralGetDto carries a single resource ID for get-by-ID operations.
type GeneralGetDto struct {
	Id int `uri:"id" json:"id" validate:"required"`
}

// MakeCondition returns a GORM scope that applies WHERE, OR, JOIN, and ORDER
// clauses derived from the search struct tags on q.
func MakeCondition(q interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		condition := &search.GormCondition{
			GormPublic: search.GormPublic{},
			Join:       make([]*search.GormJoin, 0),
		}
		search.ResolveSearchQuery(global.Driver, q, condition)
		for _, join := range condition.Join {
			if join == nil {
				continue
			}
			db = db.Joins(join.JoinOn)
			for k, v := range join.Where {
				db = db.Where(k, v...)
			}
			for k, v := range join.Or {
				db = db.Or(k, v...)
			}
			for _, o := range join.Order {
				db = db.Order(o)
			}
		}
		for k, v := range condition.Where {
			db = db.Where(k, v...)
		}
		for k, v := range condition.Or {
			db = db.Or(k, v...)
		}
		for _, o := range condition.Order {
			db = db.Order(o)
		}
		return db
	}
}

// Paginate returns a GORM scope that applies OFFSET and LIMIT for the given
// page. A negative offset is clamped to zero.
func Paginate(pageSize, pageIndex int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageIndex - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(pageSize)
	}
}

