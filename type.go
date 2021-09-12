package main

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	ID primitive.ObjectID `json:"UID" bson:"_id,omitempty"`
	Username string `json:"username" binding:"required" bson:"username,omitempty"`
	Win int	`json:"win" bson:"win"`
	Lost int `json:"lose" bson:"lose"`
	Tie int	`json:"tie" bson:"tie"`
}

type invitation struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Challenger string `json:"challenger" bson:"challenger,omitempty"`
	Challenged string `json:"challenged" bson:"challenged,omitempty"`
	Choice string `json:"choice" bson:"choice,omitempty"`
	Date primitive.DateTime `json:"date" bson:"date,omitempty"`
	MatchStatus int `json:"matchStatus" bson:"matchStatus,omitempty"`
}

type match struct {
	ID primitive.ObjectID	`json:"_id" bson:"_id,omitempty"`
	ChallengerUser string	`json:"challengerUser" bson:"challengerUser,omitempty"` 
	ChallengedUser string	`json:"challengedUser" bson:"challengedUser,omitempty"`
	ChallengerChoice	string	`json:"challengerChoice" bson:"challengerChoice,omitempty"` 
	ChallengedChoice	string	`json:"challengedChoice" bson:"challengedChoice,omitempty"`
	Date primitive.DateTime	`json:"date" bson:"date,omitempty"`
	Result	int	`json:"result" bson:"result,omitempty"`
} 