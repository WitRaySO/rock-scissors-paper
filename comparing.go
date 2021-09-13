package main

import(
	"context"
	"fmt"
	"time"
	"log"

	"net/http"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"	
)

type stat struct {
	CurrentUser string	`json:"currentUser"`
	OpponentUser string	`json:"opponentUser"`
	RecentMatches []match	`json:"recentMatches"`
	WinScore int	`json:"winScore"`
	LoseScore int	`json:"loseScore"` 
	Status string	`json:"status"`
}

func comparing(c *gin.Context) {
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

	// JSON binding currentlyUser info
	var opponentUser struct{
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&opponentUser); err != nil {
		log.Fatal(err)
	}

	currentUsername := c.Param("username")
	opponentUsername := opponentUser.Username	
	db := client.Database("rock_scissors_paper")

	// check if someone compare him to himself or not
	if currentUsername == opponentUsername {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You can't comparing yourself with yourself, compare to the others"})	
		return	
	}	
	
	// check if opponentUser exist or not
	usersColl := db.Collection("users")
	filter := bson.D{{Key: "username",Value: opponentUsername}}
	var u user
	err = usersColl.FindOne(ctx,filter).Decode(&u)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "There is not that username in system"})	
		return
	}

	// fetch data from matches collection
	matchesColl := db.Collection("matches")
	// fetch recent matches
	findOpts := options.Find().SetSort(bson.D{{"date", -1}}).SetLimit(8)
	filter = bson.D{{"$or", []interface{} { bson.D{ {"$and" ,[]interface{} { bson.D{ {"challengerUser",currentUsername}}, bson.D{{"challengedUser",opponentUsername}}}}}, bson.D{{"$and" ,[]interface{} { bson.D{ {"challengerUser",opponentUsername}}, bson.D{ {"challengedUser",currentUsername}}}}}}}}
	cursor, err :=  matchesColl.Find(ctx, filter,findOpts)
	if err != nil {
		log.Fatal(err)
	}
	var recentMatches []match
	if err = cursor.All(ctx,&recentMatches); err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("recentMatches")
	// win/lose history
	// fetch only win matches
	filter = bson.D{ { "$or", []interface{} { bson.D{ {"$and" ,[]interface{} { bson.D{ {"challengerUser",currentUsername} }, bson.D{ {"challengedUser",opponentUsername} }, bson.D{ {"result",challengerWin} } }} }, bson.D{ {"$and" ,[]interface{} { bson.D{ {"challengerUser",opponentUsername} }, bson.D{ {"challengedUser",currentUsername} }, bson.D{ {"result",challengerLose} } }} } } } } 
	cursor, err = matchesColl.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	var onlyWinMatches []match
	if err = cursor.All(ctx,&onlyWinMatches); err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("win history")
	// fetch only lose matches
	filter = bson.D{ { "$or", []interface{} { bson.D{ {"$and" ,[]interface{} { bson.D{ {"challengerUser",currentUsername} }, bson.D{ {"challengedUser",opponentUsername} }, bson.D{ {"result",challengerLose} } }} }, bson.D{ {"$and" ,[]interface{} { bson.D{ {"challengerUser",opponentUsername} }, bson.D{ {"challengedUser",currentUsername} }, bson.D{ {"result",challengerWin} } }} } } } } 
	cursor, err = matchesColl.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	var onlyLoseMatches []match
	if err = cursor.All(ctx,&onlyLoseMatches); err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("lose history")
	// status [challenger,challenged,neutral]
	// fetch data from invitation collection
	invitationColl := db.Collection("invitation")
	// check if there is a invitation that matchStatus is set to 2(matchStart)
	// findOneOpts := options.FindOne().SetSort(bson.D{{"date", -1}})
	filter = bson.D{ { "$or", []interface{} { bson.D{ {"$and" ,[]interface{} { bson.D{ {"challenger",currentUsername} }, bson.D{ {"challenged",opponentUsername} }, bson.D{ {"matchStatus",matchStart} } }} }, bson.D{ {"$and" ,[]interface{} { bson.D{ {"challenger",opponentUsername} }, bson.D{ {"challenged",currentUsername} }, bson.D{ {"matchStatus",matchStart} } }} } } } }
	var i invitation
	var userStatus string
	err = invitationColl.FindOne(ctx,filter).Decode(&i)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
		}
		userStatus = "neutral"
	}
	if i.Challenger == currentUsername {
		userStatus = "challenger"
	} else if i.Challenged == currentUsername {
		userStatus = "challenged"
	}
		
	// fmt.Printf("lose history")
	var currentUserStat stat
	currentUserStat.CurrentUser = currentUsername  
	currentUserStat.OpponentUser = opponentUsername
	currentUserStat.RecentMatches = recentMatches
	currentUserStat.WinScore = len(onlyWinMatches)
	currentUserStat.LoseScore = len(onlyLoseMatches)
	currentUserStat.Status = userStatus
	c.JSON(http.StatusOK, currentUserStat)
}