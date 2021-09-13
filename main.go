package main

import (
	"github.com/gin-gonic/gin"
)

// matchResult
const (
	_ = iota
	tie
	challengerWin
	challengerLose
)
// matchStatus
const (
	_ = iota
	matchEnd
	matchStart
)

const uri string = "mongodb+srv://witty:1234@project0.vzh5p.mongodb.net/user?retryWrites=true&w=majority"

func main() {
	router := gin.Default()
	router.GET("/getAllUsers", getAllUsers)
	router.PUT("/signup", signUp)
	router.POST("/user/:username/invitation", rockScissorPaper)
	router.GET("/user/:username/comparison", comparing)
	router.GET("/leaderboard", leaderboard)	
	router.Run("localhost:8080")
}

