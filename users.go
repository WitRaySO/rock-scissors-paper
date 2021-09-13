package main

import (
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

func getAllUsers(c *gin.Context) {
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

	coll := client.Database("rock_scissors_paper").Collection("users")
	cursor, err :=  coll.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var users []user
	if err = cursor.All(ctx,&users); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, users)
}

func signUp(c *gin.Context) {
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

	// JSON binding	
	var signupUser user
	if err := c.ShouldBindJSON(&signupUser); err != nil {
		log.Fatal(err)
	}
	coll := client.Database("rock_scissors_paper").Collection("users")
	insertResult, err := coll.InsertOne(ctx, signupUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "DuplicateKeyError"})
			return
		}
		log.Fatal(err)
	}
	c.JSON(http.StatusCreated,insertResult)
}