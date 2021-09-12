package main

import (
	"context"
	"fmt"
	"time"
	"log"

	"net/http"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func rockScissorPaper(c *gin.Context) {
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

	// JSON binding info
	var json struct{
		Username string `json:"username"`
		Choice string `json:"choice"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Fatal(err)
	}

	// name from param
	currentUsername := c.Param("username")
	currentUsernameChoice := json.Choice
	opponentUsername := json.Username
	db := client.Database("rock_scissors_paper")
	invitationColl := db.Collection("invitation")
	opts := options.FindOne().SetSort(bson.D{{Key: "date",Value: -1}})

	var currentUserInvitation invitation
	var opponentUserInvitation invitation
	var zeroValueResult invitation
	// fetch data for currentUser invitation
	err = invitationColl.FindOne(ctx, bson.D{{Key: "challenger",Value: currentUsername},{Key: "challenged",Value: opponentUsername},{Key: "matchStatus",Value: matchStart}},opts).Decode(&currentUserInvitation)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)		
		}
	}

	// fetch data for opponentUser invitation
	err = invitationColl.FindOne(ctx, bson.D{{Key: "challenger",Value: opponentUsername},{Key: "challenged",Value: currentUsername},{Key: "matchStatus",Value: matchStart}},opts).Decode(&opponentUserInvitation)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
		}
	}

	if currentUserInvitation != zeroValueResult {
		// currentUser already create invitation, maybe update a choice
		c.JSON(http.StatusOK, currentUserInvitation)
		return
	}else if opponentUserInvitation == zeroValueResult {		
		// opponentUser still not create invitation , create for him
		newInvitation := invitation{
			Challenger: currentUsername, 
			Challenged: opponentUsername,
			Choice: currentUsernameChoice,
			Date: primitive.NewDateTimeFromTime(time.Now()),
			MatchStatus: matchStart, 
		}
		insertResult, err := invitationColl.InsertOne(ctx, newInvitation)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusCreated, insertResult)
		return	
	}
	// Opponent already create invitation , no need to create it again just find a result
	// first change matchStatus to matchEnd in opponentUserInvitation
	filter := bson.D{{Key: "_id", Value: opponentUserInvitation.ID}}
	update := bson.D{{Key: "$set",Value : bson.D{{Key: "matchStatus",Value: matchEnd}}}}
	updateResult, err := invitationColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, updateResult)
	// proceed to find a winner and loser or tie
	challengerChoice := opponentUserInvitation.Choice
	challengedChoice :=	currentUsernameChoice
	newMatch := match{
		ChallengerUser: opponentUserInvitation.Challenger, 
		ChallengedUser: opponentUserInvitation.Challenged,
		ChallengerChoice: challengerChoice, 
		ChallengedChoice: challengedChoice,
		Date: primitive.NewDateTimeFromTime(time.Now()),
	}
	
	matchesColl := db.Collection("matches")
	usersColl := db.Collection("users")
	if challengerChoice == challengedChoice {
		// insert tie matched
		newMatch.Result = tie
		insertResult, err := matchesColl.InsertOne(ctx, newMatch)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusCreated, insertResult)
		// Update score in users collection for tie 
		filter := bson.D{{Key: "$or" ,Value: []interface{}	{
			 bson.M{ "username": currentUsername}, bson.M{"username":opponentUsername}}}}
		update := bson.D{{Key: "$inc",Value : bson.D{{Key: "tie",Value: 1}}}}
		updateResult, err := usersColl.UpdateMany(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
		// if updateResult.MatchedCount != 0 {
		// 	fmt.Println("matched and replaced an existing document")
		// }
		c.JSON(http.StatusOK, updateResult)
		return	
	} else if ((challengerChoice == "rock") && (challengedChoice == "scissors")) ||
	 		((challengerChoice == "scissors") && (challengedChoice == "paper")) ||
			((challengerChoice == "paper") && (challengedChoice == "rock")) {
		// insert challengerWin matched
		newMatch.Result = challengerWin
		insertResult, err := matchesColl.InsertOne(ctx, newMatch)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusCreated, insertResult)
		
		// Update score in users collection for ChallengerWin
		// Update win score to Challenger/Opponent
		filter := bson.D{{Key: "username", Value: opponentUsername}}
		update := bson.D{{Key: "$inc",Value : bson.D{{Key: "win",Value: 1}}}}
		updateResult, err := usersColl.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, updateResult)
		// Update lose score to Challenged/current		
		filter = bson.D{{Key: "username", Value: currentUsername}}
		update = bson.D{{Key: "$inc",Value : bson.D{{Key: "lose",Value: 1}}}}
		updateResult, err = usersColl.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, updateResult)
		return		
	}
	// insert challengerLose matched
	newMatch.Result = challengerLose
	insertResult, err := matchesColl.InsertOne(ctx, newMatch)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusCreated, insertResult)
	// Update score in users collection for ChallengerLose
	// Update lose score to Challenger/Opponent
	filter = bson.D{{Key: "username", Value: opponentUsername}}
	update = bson.D{{Key: "$inc",Value : bson.D{{Key: "lose",Value: 1}}}}
	updateResult, err = usersColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, updateResult)
	// Update win score to Challenged/current	
	filter = bson.D{{Key: "username", Value: currentUsername}}
	update = bson.D{{Key: "$inc",Value : bson.D{{Key: "win",Value: 1}}}}
	updateResult, err = usersColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, updateResult)	
}