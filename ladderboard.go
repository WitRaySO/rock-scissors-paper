package main

import(
	"context"
	"fmt"
	"time"
	"log"
	"sort"

	"net/http"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"	
)

type formatUserRanking struct {
	Username string `json:"username" binding:"required" bson:"username,omitempty"`
	Winrate float64 `json:"winrate" bson:"winrate"`
	Win int	`json:"win" bson:"win"`
	Lost int `json:"lose" bson:"lose"`
	Tie int	`json:"tie" bson:"tie"`
}

func ladderboard(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	// Disconnect client by defer
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Ping the primary for checking the connection to database
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected and pinged.")

	// fetch data from all users
	coll := client.Database("rock_scissors_paper").Collection("users")
	cursor, err :=  coll.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var users []user
	if err = cursor.All(ctx,&users); err != nil {
		log.Fatal(err)
	}
	// Calculate winrate then put it into slice
	var	fmtUsers []formatUserRanking
	for _,u := range users {
		newfmtUsers := formatUserRanking{
			Username: u.Username,
			Winrate: calculateWinrate(u.Win, u.Lost),
			Win: u.Win,
			Lost: u.Lost,
			Tie: u.Tie,
		}
		fmtUsers = append(fmtUsers,newfmtUsers)		
	}
	// Sort them by descending
	sort.Slice(fmtUsers, func(i, j int) bool {
		return fmtUsers[i].Winrate > fmtUsers[j].Winrate
	})
	c.JSON(http.StatusOK, fmtUsers)
}

func calculateWinrate(win int, lost int) (winrate float64) {
	if lost == 0 {
		return float64(0)
	}
	return float64(win)/float64(lost)
} 