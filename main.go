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

	. "./model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

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
		items := GetAllItems()
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

		items := SearchItems(searchWords)
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
			user.Create()
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

	// *********************
	// url: GET "/logout"
	// ログインページ
	// *********************
	r.GET("/logout", func(c *gin.Context) {
		r.LoadHTMLGlob("templates/main/*")

		ClearSession(c)
		c.Redirect(302, "/")
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

		item.Create()
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

		item.Update()
		c.Redirect(302, "/ughfkhszdlvjkdjsbfkjsdabfl/sadmin")
	})

	// assets フォルダの読み取り
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	r.Run()
}
