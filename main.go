package main

import (
	"github.com/gin-gonic/gin"
)

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
	// router.GET("/allUser", getUsers)
	router.GET("/getAllUsers", getAllUsers)
	router.PUT("/signup", signUp)
	router.POST("/user/:username/invitation", rockScissorPaper)
	router.GET("/user/:username/comparison", comparing)
	router.GET("/ladderboard", leadderboard)	
	router.Run("localhost:8080")
}
// func comparing(c *gin.Context) {

// }

// func getLadderboard(c *gin.Context) {

// }

