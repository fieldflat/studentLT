package model

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionInfo は セッション情報を保持する構造体
type SessionInfo struct {
	UserID         interface{} //ログインしているユーザのID
	Name           interface{} //ログインしているユーザの名前
	IsSessionAlive bool        //セッションが生きているかどうか
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

// ClearSession is a function
// =====================
// ClearSession 関数
// =====================
func ClearSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}
