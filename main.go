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
	Title               string `validate:"required"` // イベントタイトル
	Description         string `validate:"required"` // イベントの説明
	Price               string `validate:"required"` // イベントの参加価格
	MaxParticipants     uint64 `validate:"required"` // イベントの最大参加可能人数
	NumParticipants     uint64 `validate:"required"` // イベントの現在参加予定人数
	ScheduledDateYear   uint64 `validate:"required"` // イベントの開催予定日時 (年)
	ScheduledDateMonth  uint64 `validate:"required"` // イベントの開催予定日時 (月)
	ScheduledDateDay    uint64 `validate:"required"` // イベントの開催予定日時 (日)
	ScheduledDateHour   uint64 `validate:"required"` // イベントの開催予定日時 (時)
	ScheduledDateMinute uint64 `validate:"required"` // イベントの開催予定日時 (分)
	DeadlineDateYear    uint64 `validate:"required"` // 参加申し込みの締切日時 (年)
	DeadlineDateMonth   uint64 `validate:"required"` // 参加申し込みの締切日時 (月)
	DeadlineDateDay     uint64 `validate:"required"` // 参加申し込みの締切日時 (日)
	DeadlineDateHour    uint64 `validate:"required"` // 参加申し込みの締切日時 (時)
	DeadlineDateMinute  uint64 `validate:"required"` // 参加申し込みの締切日時 (分))
	Belongings          string `validate:"required"` // 持ち物リスト
	Target              string `validate:"required"` // 参加対象者
	Other               string // その他
	CreatedTime         string `validate:"required"` // 作成日時
	UpdatedTime         string `validate:"required"` // 更新日時
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
				case "Price":
					result["Price"] = "点数の入力は必須です．"
				case "MaxParticipants":
					result["MaxParticipants"] = "最大参加人数の入力は必須です．"
				case "NumParticipants":
					result["NumParticipants"] = "現在参加人数の入力は必須です．"
				case "ScheduledDateYear":
					result["ScheduledDateYear"] = "開催予定日時(年)の入力は必須です．"
				case "ScheduledDateMonth":
					result["ScheduledDateMonth"] = "開催予定日時(月)の入力は必須です．"
				case "ScheduledDateDay":
					result["ScheduledDateDay"] = "開催予定日時(日)の入力は必須です．"
				case "ScheduledDateHour":
					result["ScheduledDateHour"] = "開催予定日時(時)の入力は必須です．"
				case "ScheduledDateMinute":
					result["ScheduledDateMinute"] = "開催予定日時(分)の入力は必須です．"
				case "DeadlineDateYear":
					result["DeadlineDateYear"] = "申し込み締切日時(年)の入力は必須です．"
				case "DeadlineDateMonth":
					result["DeadlineDateMonth"] = "申し込み締切日時(月)の入力は必須です．"
				case "DeadlineDateDay":
					result["DeadlineDateDay"] = "申し込み締切日時(日)の入力は必須です．"
				case "DeadlineDateHour":
					result["DeadlineDateHour"] = "申し込み締切日時(時))の入力は必須です．"
				case "DeadlineDateMinute":
					result["DeadlineDateMinute"] = "申し込み締切日時(分)の入力は必須です．"
				case "Belongings":
					result["Belongings"] = "持ち物の入力は必須です．"
				case "Target":
					result["Target"] = "参加対象の入力は必須です．"
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
	price string,
	maxParticipants uint64,
	numParticipants uint64,
	scheduledDateYear uint64,
	scheduledDateMonth uint64,
	scheduledDateDay uint64,
	scheduledDateHour uint64,
	scheduledDateMinute uint64,
	deadlineDateYear uint64,
	deadlineDateMonth uint64,
	deadlineDateDay uint64,
	deadlineDateHour uint64,
	deadlineDateMinute uint64,
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
		Title:               title,
		Description:         description,
		Price:               price,
		MaxParticipants:     maxParticipants,
		NumParticipants:     numParticipants,
		ScheduledDateYear:   scheduledDateYear,
		ScheduledDateMonth:  scheduledDateMonth,
		ScheduledDateDay:    scheduledDateDay,
		ScheduledDateHour:   scheduledDateHour,
		ScheduledDateMinute: scheduledDateMinute,
		DeadlineDateYear:    deadlineDateYear,
		DeadlineDateMonth:   deadlineDateMonth,
		DeadlineDateDay:     deadlineDateDay,
		DeadlineDateHour:    deadlineDateHour,
		DeadlineDateMinute:  deadlineDateMinute,
		Belongings:          belongings,
		Target:              target,
		Other:               other,
		CreatedTime:         createdTime,
		UpdatedTime:         updatedTime,
	}

	log.Print("MaxParticipants: ", maxParticipants)
	log.Print("NumParticipants: ", numParticipants)

	// バリデーションの検証を行う
	ok, errorMessages := form.Validate()
	if !ok {
		log.Print("入力エラーあり")
		log.Print(errorMessages)
		return errorMessages
	}

	log.Print("入力エラーなし！！")
	db.Create(&Item{
		Title:               title,
		Description:         description,
		Price:               price,
		MaxParticipants:     maxParticipants,
		NumParticipants:     numParticipants,
		ScheduledDateYear:   scheduledDateYear,
		ScheduledDateMonth:  scheduledDateMonth,
		ScheduledDateDay:    scheduledDateDay,
		ScheduledDateHour:   scheduledDateHour,
		ScheduledDateMinute: scheduledDateMinute,
		DeadlineDateYear:    deadlineDateYear,
		DeadlineDateMonth:   deadlineDateMonth,
		DeadlineDateDay:     deadlineDateDay,
		DeadlineDateHour:    deadlineDateHour,
		DeadlineDateMinute:  deadlineDateMinute,
		Belongings:          belongings,
		Target:              target,
		Other:               other,
		CreatedTime:         createdTime,
		UpdatedTime:         updatedTime,
	})
	return errorMessages
}

func update(
	id int,
	title string,
	description string,
	price string,
	maxParticipants uint64,
	numParticipants uint64,
	scheduledDateYear uint64,
	scheduledDateMonth uint64,
	scheduledDateDay uint64,
	scheduledDateHour uint64,
	scheduledDateMinute uint64,
	deadlineDateYear uint64,
	deadlineDateMonth uint64,
	deadlineDateDay uint64,
	deadlineDateHour uint64,
	deadlineDateMinute uint64,
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

// main関数
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

	// 新規作成
	// r.POST("/new", func(c *gin.Context) {
	// 	r.LoadHTMLGlob("templates/main/*")
	// 	title := c.PostForm("title")
	// 	description := c.PostForm("description")
	// 	price := c.PostForm("price")
	// 	maxParticipants, _ := strconv.ParseUint(c.PostForm("maxParticipants"), 10, 32)
	// 	numParticipants, _ := strconv.ParseUint(c.PostForm("numParticipants"), 10, 32)
	// 	scheduledDateYear, _ := strconv.ParseUint(c.PostForm("scheduledDateYear"), 10, 32)
	// 	scheduledDateMonth, _ := strconv.ParseUint(c.PostForm("scheduledDateMonth"), 10, 32)
	// 	scheduledDateDay, _ := strconv.ParseUint(c.PostForm("scheduledDateDay"), 10, 32)
	// 	deadlineDateYear, _ := strconv.ParseUint(c.PostForm("deadlineDateYear"), 10, 32)
	// 	deadlineDateMonth, _ := strconv.ParseUint(c.PostForm("deadlineDateMonth"), 10, 32)
	// 	deadlineDateDay, _ := strconv.ParseUint(c.PostForm("deadlineDateDay"), 10, 32)
	// 	belongings := c.PostForm("belongings")
	// 	target := c.PostForm("target")
	// 	other := c.PostForm("other")
	// 	createdTime := c.PostForm("createdTime")
	// 	updatedTime := c.PostForm("updatedTime")

	// 	errors := create(
	// 		title,
	// 		description,
	// 		price,
	// 		maxParticipants,
	// 		numParticipants,
	// 		scheduledDateYear,
	// 		scheduledDateMonth,
	// 		scheduledDateDay,
	// 		deadlineDateYear,
	// 		deadlineDateMonth,
	// 		deadlineDateDay,
	// 		belongings,
	// 		target,
	// 		other,
	// 		createdTime,
	// 		updatedTime,
	// 	)
	// 	log.Print("errors : ", errors)
	// 	if len(errors) != 0 {
	// 		items := getAll()
	// 		c.HTML(200, "index.tmpl", gin.H{
	// 			"items":  items,
	// 			"errors": errors,
	// 		})
	// 	} else {
	// 		log.Print("こっちあよ22！！！！！！")
	// 		c.Redirect(302, "/")
	// 	}
	// })

	// admin TOPページ
	r.GET("/ughfkhszdlvjkdjsbfkjsdabfl/sadmin", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/sadmin/*")
		c.HTML(200, "sadmin_index.tmpl", gin.H{})
	})

	r.POST("/ughfkhszdlvjkdjsbfkjsdabfl/sadmin/create", func(c *gin.Context) {
		title := c.PostForm("title")
		description := c.PostForm("description")
		price := c.PostForm("price")
		maxParticipants, _ := strconv.ParseUint(c.PostForm("maxParticipants"), 10, 32)
		numParticipants, _ := strconv.ParseUint("1", 10, 32)

		scheduledDate := c.PostForm("scheduledDate")
		scheduledDateSlice := strings.Split(scheduledDate, "-")
		scheduledDateYear, _ := strconv.ParseUint(scheduledDateSlice[0], 10, 32)
		scheduledDateMonth, _ := strconv.ParseUint(scheduledDateSlice[1], 10, 32)
		scheduledDateDay, _ := strconv.ParseUint(scheduledDateSlice[2], 10, 32)

		scheduledTime := c.PostForm("scheduledTime")
		scheduledTimeSlice := strings.Split(scheduledTime, ":")
		scheduledDateHour, _ := strconv.ParseUint(scheduledTimeSlice[0], 10, 32)
		scheduledDateMinute, _ := strconv.ParseUint(scheduledTimeSlice[1], 10, 32)

		deadlineDate := c.PostForm("deadlineDate")
		deadlineDateSlice := strings.Split(deadlineDate, "-")
		deadlineDateYear, _ := strconv.ParseUint(deadlineDateSlice[0], 10, 32)
		deadlineDateMonth, _ := strconv.ParseUint(deadlineDateSlice[1], 10, 32)
		deadlineDateDay, _ := strconv.ParseUint(deadlineDateSlice[2], 10, 32)

		deadlineTime := c.PostForm("deadlineTime")
		deadlineTimeSlice := strings.Split(deadlineTime, ":")
		deadlineDateHour, _ := strconv.ParseUint(deadlineTimeSlice[0], 10, 32)
		deadlineDateMinute, _ := strconv.ParseUint(deadlineTimeSlice[1], 10, 32)

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
		c.Redirect(302, "/")
	})

	// r.POST("/update/:id", func(c *gin.Context) {
	// 	title := c.PostForm("title")
	// 	description := c.PostForm("description")
	// 	point := c.PostForm("point")
	// 	n := c.Param("id")
	// 	id, err := strconv.Atoi(n)
	// 	if err != nil {
	// 		panic("failed to get id\n")
	// 	}

	// 	update(id, title, description, point, time.Now().Format("2006/1/2 15:04:05"))
	// 	c.Redirect(302, "/")
	// })

	// assets フォルダの読み取り
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	r.Run()
}
