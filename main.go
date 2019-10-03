package main

import (
	"log"
	"net/http"
	_ "net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
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
func (form *Item) Validate() (ok bool, result map[string]string) {
	result = make(map[string]string)
	// 構造体のデータをタグで定義した検証方法でチェック
	// err := validator.New().Struct(*form)
	validate := validator.New()
	// validate.RegisterValidation("is_tarou", tarou) //第一引数をvalidateタグで設定した名前に合わせる
	err := validate.Struct(*form)
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

// DBの初期化処理
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}

	db.AutoMigrate(&Item{})
}

// create関数
func create(
	title string,
	description string,
	price int,
	maxParticipants int,
	numParticipants int,
	scheduledDateYear int,
	scheduledDateMonth int,
	scheduledDateDay int,
	scheduledDateHour int,
	scheduledDateMinute int,
	scheduledDateEndYear int,
	scheduledDateEndMonth int,
	scheduledDateEndDay int,
	scheduledDateEndHour int,
	scheduledDateEndMinute int,
	deadlineDateYear int,
	deadlineDateMonth int,
	deadlineDateDay int,
	deadlineDateHour int,
	deadlineDateMinute int,
	belongings string,
	target string,
	other string,
	createdTime string,
	updatedTime string,
) map[string]string {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}

	// 処理を追加
	form := Item{
		Title:                  title,
		Description:            description,
		Price:                  price,
		MaxParticipants:        maxParticipants,
		NumParticipants:        numParticipants,
		ScheduledDateYear:      scheduledDateYear,
		ScheduledDateMonth:     scheduledDateMonth,
		ScheduledDateDay:       scheduledDateDay,
		ScheduledDateHour:      scheduledDateHour,
		ScheduledDateMinute:    scheduledDateMinute,
		ScheduledDateEndYear:   scheduledDateEndYear,
		ScheduledDateEndMonth:  scheduledDateEndMonth,
		ScheduledDateEndDay:    scheduledDateEndDay,
		ScheduledDateEndHour:   scheduledDateEndHour,
		ScheduledDateEndMinute: scheduledDateEndMinute,
		DeadlineDateYear:       deadlineDateYear,
		DeadlineDateMonth:      deadlineDateMonth,
		DeadlineDateDay:        deadlineDateDay,
		DeadlineDateHour:       deadlineDateHour,
		DeadlineDateMinute:     deadlineDateMinute,
		Belongings:             belongings,
		Target:                 target,
		Other:                  other,
		CreatedTime:            createdTime,
		UpdatedTime:            updatedTime,
	}

	log.Print(form)

	// バリデーションの検証を行う
	ok, errorMessages := form.Validate()
	if !ok {
		log.Print("入力エラーあり")
		log.Print(errorMessages)
		return errorMessages
	}

	log.Print("入力エラーなし！！")
	db.Create(&Item{
		Title:                  title,
		Description:            description,
		Price:                  price,
		MaxParticipants:        maxParticipants,
		NumParticipants:        numParticipants,
		ScheduledDateYear:      scheduledDateYear,
		ScheduledDateMonth:     scheduledDateMonth,
		ScheduledDateDay:       scheduledDateDay,
		ScheduledDateHour:      scheduledDateHour,
		ScheduledDateMinute:    scheduledDateMinute,
		ScheduledDateEndYear:   scheduledDateEndYear,
		ScheduledDateEndMonth:  scheduledDateEndMonth,
		ScheduledDateEndDay:    scheduledDateEndDay,
		ScheduledDateEndHour:   scheduledDateEndHour,
		ScheduledDateEndMinute: scheduledDateEndMinute,
		DeadlineDateYear:       deadlineDateYear,
		DeadlineDateMonth:      deadlineDateMonth,
		DeadlineDateDay:        deadlineDateDay,
		DeadlineDateHour:       deadlineDateHour,
		DeadlineDateMinute:     deadlineDateMinute,
		Belongings:             belongings,
		Target:                 target,
		Other:                  other,
		CreatedTime:            createdTime,
		UpdatedTime:            updatedTime,
	})
	return errorMessages
}

func update(
	id int,
	title string,
	description string,
	price int,
	maxParticipants int,
	numParticipants int,
	scheduledDateYear int,
	scheduledDateMonth int,
	scheduledDateDay int,
	scheduledDateHour int,
	scheduledDateMinute int,
	scheduledDateEndYear int,
	scheduledDateEndMonth int,
	scheduledDateEndDay int,
	scheduledDateEndHour int,
	scheduledDateEndMinute int,
	deadlineDateYear int,
	deadlineDateMonth int,
	deadlineDateDay int,
	deadlineDateHour int,
	deadlineDateMinute int,
	belongings string,
	target string,
	other string,
	createdTime string,
	updatedTime string,
) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	var item Item
	db.First(&item, id)
	item.Title = title
	item.Description = description
	item.Price = price
	item.MaxParticipants = maxParticipants
	item.NumParticipants = numParticipants
	item.ScheduledDateYear = scheduledDateYear
	item.ScheduledDateMonth = scheduledDateMonth
	item.ScheduledDateDay = scheduledDateDay
	item.ScheduledDateHour = scheduledDateHour
	item.ScheduledDateMinute = scheduledDateMinute
	item.ScheduledDateEndYear = scheduledDateEndYear
	item.ScheduledDateEndMonth = scheduledDateEndMonth
	item.ScheduledDateEndDay = scheduledDateEndDay
	item.ScheduledDateEndHour = scheduledDateEndHour
	item.ScheduledDateEndMinute = scheduledDateEndMinute
	item.DeadlineDateYear = deadlineDateYear
	item.DeadlineDateMonth = deadlineDateMonth
	item.DeadlineDateDay = deadlineDateDay
	item.DeadlineDateHour = deadlineDateHour
	item.DeadlineDateMinute = deadlineDateMinute
	item.Belongings = belongings
	item.Target = target
	item.Other = other
	item.CreatedTime = createdTime
	item.UpdatedTime = updatedTime
	db.Save(&item)
	db.Close()
}

// 全てのItemを取得する
func getAll() []Item {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	var items []Item
	db.Order("created_at desc").Find(&items)
	return items
}

// 検索条件を満たすItemを取得する
func search(searchWords map[string]string) []Item {
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

//
// main関数
//
func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")
	dbInit()

	// 一覧取得
	r.GET("/", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/main/*")
		items := getAll()
		c.HTML(200, "index.tmpl", gin.H{
			"items": items,
		})
	})

	// 検索結果取得
	r.GET("/search", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/main/*")
		// var searchWords map[string]string
		searchWords := map[string]string{}
		searchWords["words"] = c.Query("words")
		searchWords["scheduledDateFrom"] = c.Query("scheduledDateFrom")
		searchWords["scheduledDateTo"] = c.Query("scheduledDateTo")
		searchWords["priceFrom"] = c.Query("priceFrom")
		searchWords["priceTo"] = c.Query("priceTo")

		items := search(searchWords)
		c.HTML(200, "index.tmpl", gin.H{
			"items":       items,
			"searchWords": searchWords,
		})
	})

	// sadmin TOPページ
	r.GET("/ughfkhszdlvjkdjsbfkjsdabfl/sadmin", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/sadmin/*")
		c.HTML(200, "sadmin_index.tmpl", gin.H{})
	})

	// sadmin イベント情報作成ページ
	r.POST("/ughfkhszdlvjkdjsbfkjsdabfl/sadmin/create", func(c *gin.Context) {
		title := c.PostForm("title")
		description := c.PostForm("description")
		price, _ := strconv.Atoi(c.PostForm("price"))
		maxParticipants, _ := strconv.Atoi(c.PostForm("maxParticipants"))
		numParticipants := 0

		scheduledDate := c.PostForm("scheduledDate")
		scheduledDateSlice := strings.Split(scheduledDate, "-")
		scheduledDateYear, _ := strconv.Atoi(scheduledDateSlice[0])
		scheduledDateMonth, _ := strconv.Atoi(scheduledDateSlice[1])
		scheduledDateDay, _ := strconv.Atoi(scheduledDateSlice[2])

		scheduledTime := c.PostForm("scheduledTime")
		scheduledTimeSlice := strings.Split(scheduledTime, ":")
		scheduledDateHour, _ := strconv.Atoi(scheduledTimeSlice[0])
		scheduledDateMinute, _ := strconv.Atoi(scheduledTimeSlice[1])

		scheduledDateEnd := c.PostForm("scheduledDateEnd")
		scheduledDateEndSlice := strings.Split(scheduledDateEnd, "-")
		scheduledDateEndYear, _ := strconv.Atoi(scheduledDateEndSlice[0])
		scheduledDateEndMonth, _ := strconv.Atoi(scheduledDateEndSlice[1])
		scheduledDateEndDay, _ := strconv.Atoi(scheduledDateEndSlice[2])

		scheduledEndTime := c.PostForm("scheduledEndTime")
		scheduledEndTimeSlice := strings.Split(scheduledEndTime, ":")
		scheduledDateEndHour, _ := strconv.Atoi(scheduledEndTimeSlice[0])
		scheduledDateEndMinute, _ := strconv.Atoi(scheduledEndTimeSlice[1])

		deadlineDate := c.PostForm("deadlineDate")
		deadlineDateSlice := strings.Split(deadlineDate, "-")
		deadlineDateYear, _ := strconv.Atoi(deadlineDateSlice[0])
		deadlineDateMonth, _ := strconv.Atoi(deadlineDateSlice[1])
		deadlineDateDay, _ := strconv.Atoi(deadlineDateSlice[2])

		deadlineTime := c.PostForm("deadlineTime")
		deadlineTimeSlice := strings.Split(deadlineTime, ":")
		deadlineDateHour, _ := strconv.Atoi(deadlineTimeSlice[0])
		deadlineDateMinute, _ := strconv.Atoi(deadlineTimeSlice[1])

		belongings := c.PostForm("belongings")
		target := c.PostForm("target")
		other := c.PostForm("other")
		createdTime := time.Now().Format("2006/1/2 15:04:05")
		updatedTime := time.Now().Format("2006/1/2 15:04:05")

		create(
			title,
			description,
			price,
			maxParticipants,
			numParticipants,
			scheduledDateYear,
			scheduledDateMonth,
			scheduledDateDay,
			scheduledDateHour,
			scheduledDateMinute,
			scheduledDateEndYear,
			scheduledDateEndMonth,
			scheduledDateEndDay,
			scheduledDateEndHour,
			scheduledDateEndMinute,
			deadlineDateYear,
			deadlineDateMonth,
			deadlineDateDay,
			deadlineDateHour,
			deadlineDateMinute,
			belongings,
			target,
			other,
			createdTime,
			updatedTime,
		)
		c.Redirect(302, "/ughfkhszdlvjkdjsbfkjsdabfl/sadmin")
	})

	// sadmin Editページ
	r.GET("/ughfkhszdlvjkdjsbfkjsdabfl/edit", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/sadmin/*")
		db, err := gorm.Open("sqlite3", "test.sqlite3")
		if err != nil {
			panic("failed to connect database\n")
		}

		var item Item
		id := c.Query("id")
		db.First(&item, id)

		if item.ID == 0 {
			log.Print("item.ID = ", item.ID)
			c.Redirect(302, "/ughfkhszdlvjkdjsbfkjsdabfl/sadmin")
		} else {

			scheduledDate := strconv.Itoa(item.ScheduledDateYear) + "-" + strconv.Itoa(item.ScheduledDateMonth) + "-" + strconv.Itoa(item.ScheduledDateDay)
			scheduledTime := strconv.Itoa(item.ScheduledDateHour) + ":" + strconv.Itoa(item.ScheduledDateMinute)
			scheduledDateEnd := strconv.Itoa(item.ScheduledDateEndYear) + "-" + strconv.Itoa(item.ScheduledDateEndMonth) + "-" + strconv.Itoa(item.ScheduledDateEndDay)
			scheduledEndTime := strconv.Itoa(item.ScheduledDateEndHour) + ":" + strconv.Itoa(item.ScheduledDateEndMinute)
			deadlineDate := strconv.Itoa(item.DeadlineDateYear) + "-" + strconv.Itoa(item.DeadlineDateMonth) + "-" + strconv.Itoa(item.DeadlineDateDay)
			deadlineTime := strconv.Itoa(item.DeadlineDateHour) + ":" + strconv.Itoa(item.DeadlineDateMinute)
			log.Print("scheduledDate = ", scheduledDate)

			c.HTML(200, "sadmin_index.tmpl", gin.H{
				"item":             item,
				"scheduledDate":    scheduledDate,
				"scheduledTime":    scheduledTime,
				"scheduledDateEnd": scheduledDateEnd,
				"scheduledEndTime": scheduledEndTime,
				"deadlineDate":     deadlineDate,
				"deadlineTime":     deadlineTime,
			})
		}
	})

	// assets フォルダの読み取り
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	r.Run()
}
