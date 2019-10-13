package model

import (
	"log"

	"github.com/jinzhu/gorm"
)

// Participant is a struct
type Participant struct {
	gorm.Model
	ItemID int // itemID
	UserID int // userID
}

// Create is a function
// ===================
// Create関数
// ===================
func (part *Participant) Create() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	var parts []Participant
	query := db.Order("created_at desc").Where("item_id = (?) AND user_id = (?)", part.ItemID, part.UserID)
	query.Find(&parts)

	if len(parts) == 0 && part.UserID != -1 {
		db.Create(&part)
	} else {
		log.Printf("申し込みできませんでした．すでに申し込みが完了している可能性があります．")
	}
}
