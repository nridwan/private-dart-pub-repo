package utils

import (
	"context"

	"gorm.io/gorm"
)

func GetDBWithContext(db *gorm.DB, contexts []context.Context) *gorm.DB {
	if len(contexts) == 0 {
		return db
	}
	return db.WithContext(contexts[0])
}
