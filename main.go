package main

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"net/http"
	_ "net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

// User はユーザ情報
type User struct {
	// 大文字だと Public 扱い
	ID       int    `json:"id"`
	Name     string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SessionInfo は セッション情報を保持する構造体
type SessionInfo struct {
	UserID         interface{} //ログインしているユーザのID
	Name           interface{} //ログインしているユーザの名前
	IsSessionAlive bool        //セッションが生きているかどうか
}

// Validate about Item structure.
// ==============================
// 構造体Itemのバリデーション
// ==============================
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

// ======================
// DBの初期化処理
// Initialization of DB
// ======================
func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}

	db.AutoMigrate(&Item{})
	db.AutoMigrate(&User{})
}

// ===================
// create関数
// ===================
func create(item Item) map[string]string {
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

// ===================
// createUser関数
// ===================
func createUser(user User) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}

	db.Create(&user)
}

// ====================
// update関数
// ====================
func update(item Item) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	db.Save(&item)
	db.Close()
}

// ====================
// getAll関数
// 全てのItemを取得する
// ====================
func getAll() []Item {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("failed to connect database\n")
	}
	var items []Item
	db.Order("created_at desc").Find(&items)
	return items
}

// =============================
// search関数
// 検索条件を満たすItemを取得する
// =============================
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

// Login is a function
// =====================
// Login 関数
// =====================
func Login(g *gin.Context, user User) {
	session := sessions.Default(g)
	session.Set("alive", true)
	session.Set("userID", user.ID)
	session.Set("name", user.Name)
	session.Save()
}

// GetSessionInfo is a function
// =====================
// GetSessionInfo 関数
// =====================
func GetSessionInfo(c *gin.Context) SessionInfo {
	var info SessionInfo
	session := sessions.Default(c)
	userID := session.Get("userID")
	name := session.Get("name")
	alive := session.Get("alive")
	// if isNil(userID) && isNil(name) && isNil(alive) {
	if userID == nil && name == nil && alive == nil {
		info = SessionInfo{
			UserID: -1, Name: "", IsSessionAlive: false,
		}
	} else {
		info = SessionInfo{
			UserID:         userID.(int),
			Name:           name.(string),
			IsSessionAlive: alive.(bool),
		}
	}
	log.Println(info)
	return info
}

// isNil is a function
// =====================
// isNil 関数
// =====================
// func isNil(a interface{}) bool {
// 	return a == nil || reflect.ValueOf(a).IsNil()
// }

// ================
// main関数
// ================
func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.Static("/assets", "./assets")
	dbInit()

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~
	//
	// 以下，ユーザー操作の処理
	//
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~

	// *********************
	// url: "/"
	// 一覧取得
	// *********************
	r.GET("/", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/main/*")
		items := getAll()
		info := GetSessionInfo(c)
		c.HTML(200, "index.tmpl", gin.H{
			"items":       items,
			"SessionInfo": info,
		})
	})

	// *********************
	// url: "/search"
	// 検索結果取得
	// *********************
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
		info := GetSessionInfo(c)
		c.HTML(200, "index.tmpl", gin.H{
			"items":       items,
			"searchWords": searchWords,
			"SessionInfo": info,
		})
	})

	// *********************
	// url: "/detail"
	// 詳細ページ
	// *********************
	r.GET("/detail", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/main/*")
		db, err := gorm.Open("sqlite3", "test.sqlite3")
		if err != nil {
			panic("failed to connect database\n")
		}

		var item Item
		id := c.Query("id")
		db.First(&item, id)

		info := GetSessionInfo(c)
		c.HTML(200, "detail.tmpl", gin.H{
			"item":        item,
			"SessionInfo": info,
		})
	})

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~
	//
	// 以下，ログインの処理
	//
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~

	// *********************
	// url: GET "/signin"
	// サインインページ
	// *********************
	r.GET("/signin", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/main/*")

		c.HTML(200, "signin.tmpl", gin.H{})
	})

	// *********************
	// url: POST "/signin"
	// ユーザ登録
	// *********************
	r.POST("/signin", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/main/*")

		user := User{}
		user.Name = c.PostForm("name")
		user.Email = c.PostForm("email")
		bytes := []byte(c.PostForm("password"))
		hashPassword := sha512.Sum512(bytes)
		user.Password = hex.EncodeToString(hashPassword[:])
		log.Printf("%x\n", sha512.Sum512(bytes))

		if user.Name == "" || user.Email == "" {
			c.HTML(200, "signin.tmpl", gin.H{})
		} else {
			createUser(user)
			c.HTML(200, "index.tmpl", gin.H{})
		}
	})

	// *********************
	// url: GET "/login"
	// ログインページ
	// *********************
	r.GET("/login", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/main/*")

		c.HTML(200, "login.tmpl", gin.H{})
	})

	// *********************
	// url: POST "/login"
	// ユーザ登録
	// *********************
	r.POST("/login", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/main/*")
		db, err := gorm.Open("sqlite3", "test.sqlite3")
		if err != nil {
			panic("failed to connect database\n")
		}

		user := User{}
		email := c.PostForm("email")
		bytes := []byte(c.PostForm("password"))
		hashPassword := sha512.Sum512(bytes)
		password := hex.EncodeToString(hashPassword[:])

		db.First(&user, "email = (?) AND password = (?)", email, password)
		log.Println(user)

		if user.ID == 0 {
			c.HTML(200, "login.tmpl", gin.H{})
		} else {
			Login(c, user)
			info := GetSessionInfo(c)
			c.HTML(200, "index.tmpl", gin.H{
				"SessionInfo": info,
			})
		}
	})

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~
	//
	// 以下，sadminの処理
	//
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~

	// *******************************************
	// url: "/ughfkhszdlvjkdjsbfkjsdabfl/sadmin"
	// sadmin TOPページ
	// *******************************************
	r.GET("/ughfkhszdlvjkdjsbfkjsdabfl/sadmin", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/sadmin/*")
		c.HTML(200, "sadmin_index.tmpl", gin.H{})
	})

	// *************************************************
	// url: "/ughfkhszdlvjkdjsbfkjsdabfl/sadmin/create"
	// sadmin イベント情報作成ページ
	// *************************************************
	r.POST("/ughfkhszdlvjkdjsbfkjsdabfl/sadmin/create", func(c *gin.Context) {
		item := Item{}

		item.Title = c.PostForm("title")
		item.Description = c.PostForm("description")
		item.Price, _ = strconv.Atoi(c.PostForm("price"))
		item.MaxParticipants, _ = strconv.Atoi(c.PostForm("maxParticipants"))
		item.NumParticipants = 0

		scheduledDate := c.PostForm("scheduledDate")
		scheduledDateSlice := strings.Split(scheduledDate, "-")
		item.ScheduledDateYear, _ = strconv.Atoi(scheduledDateSlice[0])
		item.ScheduledDateMonth, _ = strconv.Atoi(scheduledDateSlice[1])
		item.ScheduledDateDay, _ = strconv.Atoi(scheduledDateSlice[2])

		scheduledTime := c.PostForm("scheduledTime")
		scheduledTimeSlice := strings.Split(scheduledTime, ":")
		item.ScheduledDateHour, _ = strconv.Atoi(scheduledTimeSlice[0])
		item.ScheduledDateMinute, _ = strconv.Atoi(scheduledTimeSlice[1])

		scheduledDateEnd := c.PostForm("scheduledDateEnd")
		scheduledDateEndSlice := strings.Split(scheduledDateEnd, "-")
		item.ScheduledDateEndYear, _ = strconv.Atoi(scheduledDateEndSlice[0])
		item.ScheduledDateEndMonth, _ = strconv.Atoi(scheduledDateEndSlice[1])
		item.ScheduledDateEndDay, _ = strconv.Atoi(scheduledDateEndSlice[2])

		scheduledEndTime := c.PostForm("scheduledEndTime")
		scheduledEndTimeSlice := strings.Split(scheduledEndTime, ":")
		item.ScheduledDateEndHour, _ = strconv.Atoi(scheduledEndTimeSlice[0])
		item.ScheduledDateEndMinute, _ = strconv.Atoi(scheduledEndTimeSlice[1])

		deadlineDate := c.PostForm("deadlineDate")
		deadlineDateSlice := strings.Split(deadlineDate, "-")
		item.DeadlineDateYear, _ = strconv.Atoi(deadlineDateSlice[0])
		item.DeadlineDateMonth, _ = strconv.Atoi(deadlineDateSlice[1])
		item.DeadlineDateDay, _ = strconv.Atoi(deadlineDateSlice[2])

		deadlineTime := c.PostForm("deadlineTime")
		deadlineTimeSlice := strings.Split(deadlineTime, ":")
		item.DeadlineDateHour, _ = strconv.Atoi(deadlineTimeSlice[0])
		item.DeadlineDateMinute, _ = strconv.Atoi(deadlineTimeSlice[1])

		item.Belongings = c.PostForm("belongings")
		item.Target = c.PostForm("target")
		item.Other = c.PostForm("other")
		item.CreatedTime = time.Now().Format("2006/1/2 15:04:05")
		item.UpdatedTime = time.Now().Format("2006/1/2 15:04:05")

		create(item)
		c.Redirect(302, "/ughfkhszdlvjkdjsbfkjsdabfl/sadmin")
	})

	// *************************************************
	// url: "/ughfkhszdlvjkdjsbfkjsdabfl/sadmin/edit"
	// sadmin Editページ
	// *************************************************
	r.GET("/ughfkhszdlvjkdjsbfkjsdabfl/sadmin/edit", func(c *gin.Context) {
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
				"isEdit":           true,
			})
		}
	})

	// *************************************************
	// url: "/ughfkhszdlvjkdjsbfkjsdabfl/sadmin/update"
	// sadmin update
	// *************************************************
	r.POST("/ughfkhszdlvjkdjsbfkjsdabfl/sadmin/update", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/sadmin/*")

		item := Item{}

		id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
		item.ID = uint(id)
		item.Title = c.PostForm("title")
		item.Description = c.PostForm("description")
		item.Price, _ = strconv.Atoi(c.PostForm("price"))
		item.MaxParticipants, _ = strconv.Atoi(c.PostForm("maxParticipants"))
		item.NumParticipants = 0

		scheduledDate := c.PostForm("scheduledDate")
		scheduledDateSlice := strings.Split(scheduledDate, "-")
		item.ScheduledDateYear, _ = strconv.Atoi(scheduledDateSlice[0])
		item.ScheduledDateMonth, _ = strconv.Atoi(scheduledDateSlice[1])
		item.ScheduledDateDay, _ = strconv.Atoi(scheduledDateSlice[2])

		scheduledTime := c.PostForm("scheduledTime")
		scheduledTimeSlice := strings.Split(scheduledTime, ":")
		item.ScheduledDateHour, _ = strconv.Atoi(scheduledTimeSlice[0])
		item.ScheduledDateMinute, _ = strconv.Atoi(scheduledTimeSlice[1])

		scheduledDateEnd := c.PostForm("scheduledDateEnd")
		scheduledDateEndSlice := strings.Split(scheduledDateEnd, "-")
		item.ScheduledDateEndYear, _ = strconv.Atoi(scheduledDateEndSlice[0])
		item.ScheduledDateEndMonth, _ = strconv.Atoi(scheduledDateEndSlice[1])
		item.ScheduledDateEndDay, _ = strconv.Atoi(scheduledDateEndSlice[2])

		scheduledEndTime := c.PostForm("scheduledEndTime")
		scheduledEndTimeSlice := strings.Split(scheduledEndTime, ":")
		item.ScheduledDateEndHour, _ = strconv.Atoi(scheduledEndTimeSlice[0])
		item.ScheduledDateEndMinute, _ = strconv.Atoi(scheduledEndTimeSlice[1])

		deadlineDate := c.PostForm("deadlineDate")
		deadlineDateSlice := strings.Split(deadlineDate, "-")
		item.DeadlineDateYear, _ = strconv.Atoi(deadlineDateSlice[0])
		item.DeadlineDateMonth, _ = strconv.Atoi(deadlineDateSlice[1])
		item.DeadlineDateDay, _ = strconv.Atoi(deadlineDateSlice[2])

		deadlineTime := c.PostForm("deadlineTime")
		deadlineTimeSlice := strings.Split(deadlineTime, ":")
		item.DeadlineDateHour, _ = strconv.Atoi(deadlineTimeSlice[0])
		item.DeadlineDateMinute, _ = strconv.Atoi(deadlineTimeSlice[1])

		item.Belongings = c.PostForm("belongings")
		item.Target = c.PostForm("target")
		item.Other = c.PostForm("other")
		item.CreatedTime = time.Now().Format("2006/1/2 15:04:05")
		item.UpdatedTime = time.Now().Format("2006/1/2 15:04:05")

		update(item)
		c.Redirect(302, "/ughfkhszdlvjkdjsbfkjsdabfl/sadmin")
	})

	// assets フォルダの読み取り
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	r.Run()
}
