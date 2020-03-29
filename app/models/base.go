package models

import (
	"time"
	"math/rand"
	"github.com/jinzhu/gorm"
	"github.com/oklog/ulid"
)

// BaseObject basic object with a ULID as its id
type BaseObject struct {
	ID  string    `json: "id" gorm: primary_key`
}

func getULID() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

// BeforeCreate Run before creating an object. Set ID as a new ULID
func (base *BaseObject) BeforeCreate(scope *gorm.Scope) error {
	ulid := getULID()
	return scope.SetColumn("ID", ulid)
}


// DBMigrate automatically create / migrate all tables
func DBMigrate(db *gorm.DB) {
	db.AutoMigrate(&Image{})
}
