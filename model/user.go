package model

import "github.com/jinzhu/gorm"

// User はユーザ情報
type User struct {
	// 大文字だと Public 扱い
	ID       int    `json:"id"`
	Name     string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create is a function
// ===================
// Create関数
// ===================
func (user *User) Create() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}

	db.Create(&user)
}
