package model

import (
	"log"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

// Item は出品アイテム情報
type Item struct {
	gorm.Model
	Title                  string `validate:"required"` // イベントタイトル
	Description            string `validate:"required"` // イベントの説明
	Price                  int    // イベントの参加価格
	MaxParticipants        int    // イベントの最大参加可能人数
	NumParticipants        int    // イベントの現在参加予定人数
	ScheduledDateYear      int    // イベントの開催予定日時 (年)
	ScheduledDateMonth     int    // イベントの開催予定日時 (月)
	ScheduledDateDay       int    // イベントの開催予定日時 (日)
	ScheduledDateHour      int    // イベントの開催予定日時 (時)
	ScheduledDateMinute    int    // イベントの開催予定日時 (分)
	ScheduledDateEndYear   int    // イベントの開催終了日時 (年)
	ScheduledDateEndMonth  int    // イベントの開催終了日時 (月)
	ScheduledDateEndDay    int    // イベントの開催終了日時 (日)
	ScheduledDateEndHour   int    // イベントの開催終了日時 (時)
	ScheduledDateEndMinute int    // イベントの開催終了日時 (分)
	DeadlineDateYear       int    // 参加申し込みの締切日時 (年)
	DeadlineDateMonth      int    // 参加申し込みの締切日時 (月)
	DeadlineDateDay        int    // 参加申し込みの締切日時 (日)
	DeadlineDateHour       int    // 参加申し込みの締切日時 (時)
	DeadlineDateMinute     int    // 参加申し込みの締切日時 (分))
	Belongings             string // 持ち物リスト
	Target                 string // 参加対象者
	Other                  string // その他
	CreatedTime            string `validate:"required"` // 作成日時
	UpdatedTime            string `validate:"required"` // 更新日時
}

// Validate about Item structure.
// ==============================
// 構造体Itemのバリデーション
// ==============================
func (item *Item) Validate() (ok bool, result map[string]string) {
	result = make(map[string]string)
	// 構造体のデータをタグで定義した検証方法でチェック
	// err := validator.New().Struct(*item)
	validate := validator.New()
	// validate.RegisterValidation("is_tarou", tarou) //第一引数をvalidateタグで設定した名前に合わせる
	err := validate.Struct(*item)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		if len(errors) != 0 {
			for i := range errors {
				// フィールドごとに、検証
				switch errors[i].StructField() {
				case "Title":
					result["Title"] = "タイトルの入力は必須です．"
				case "Description":
					result["Description"] = "本文の入力は必須です．"
				case "CreatedTime":
					result["CreatedTime"] = "CreatedTimeの入力は必須です．"
				case "UpdatedTime":
					result["UpdatedTime"] = "UpdatedTimeの入力は必須です．"
				}
			}
		}
		return false, result
	}
	return true, result
}

// Create is a function
// ===================
// create関数
// ===================
func (item *Item) Create() map[string]string {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}

	// バリデーションの検証を行う
	ok, errorMessages := item.Validate()
	if !ok {
		log.Print("入力エラーあり")
		log.Print(errorMessages)
		return errorMessages
	}

	log.Print("入力エラーなし！！")
	db.Create(&item)
	return errorMessages
}

// Update is a function
// ====================
// Update関数
// ====================
func (item *Item) Update() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	db.Save(&item)
	db.Close()
}

// GetAllItems is a function
// ====================
// GetAllItems関数
// 全てのItemを取得する
// ====================
func GetAllItems() []Item {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	var items []Item
	db.Order("created_at desc").Find(&items)
	return items
}

// SearchItems is a function
// =============================
// search関数
// 検索条件を満たすItemを取得する
// =============================
func SearchItems(searchWords map[string]string) []Item {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	var items []Item
	query := db.Order("created_at desc")
	log.Print(searchWords)

	if len(searchWords["words"]) != 0 {
		query = query.Where("title LIKE (?) OR description LIKE (?)", "%"+searchWords["words"]+"%", "%"+searchWords["words"]+"%")
	}

	if searchWords["scheduledDateFrom"] != "" {
		scheduledDateFrom := searchWords["scheduledDateFrom"]
		scheduledDateFromSlice := strings.Split(scheduledDateFrom, "-")
		scheduledDateFromYear, _ := strconv.ParseUint(scheduledDateFromSlice[0], 10, 32)
		scheduledDateFromMonth, _ := strconv.ParseUint(scheduledDateFromSlice[1], 10, 32)
		scheduledDateFromDay, _ := strconv.ParseUint(scheduledDateFromSlice[2], 10, 32)

		query = query.Where("scheduled_date_end_year >= (?) ", scheduledDateFromYear)
		query = query.Where("scheduled_date_end_month >= (?) ", scheduledDateFromMonth)
		query = query.Where("scheduled_date_end_day >= (?) ", scheduledDateFromDay)
	}

	if searchWords["scheduledDateTo"] != "" {
		scheduledDateTo := searchWords["scheduledDateTo"]
		scheduledDateToSlice := strings.Split(scheduledDateTo, "-")
		scheduledDateToYear, _ := strconv.ParseUint(scheduledDateToSlice[0], 10, 32)
		scheduledDateToMonth, _ := strconv.ParseUint(scheduledDateToSlice[1], 10, 32)
		scheduledDateToDay, _ := strconv.ParseUint(scheduledDateToSlice[2], 10, 32)

		query = query.Where("scheduled_date_year <= (?) ", scheduledDateToYear)
		query = query.Where("scheduled_date_month <= (?) ", scheduledDateToMonth)
		query = query.Where("scheduled_date_day <= (?) ", scheduledDateToDay)
	}

	if searchWords["priceFrom"] != "" {
		priceFrom, _ := strconv.Atoi(searchWords["priceFrom"])
		log.Print("priceFrom = ", priceFrom)
		query = query.Where("price >= (?) ", priceFrom)
	}

	if searchWords["priceTo"] != "" {
		priceTo, _ := strconv.Atoi(searchWords["priceTo"])
		log.Print("priceTo = ", priceTo)
		query = query.Where("price <= (?) ", priceTo)
	}

	query.Find(&items)
	return items
}
