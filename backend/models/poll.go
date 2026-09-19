package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PollOption struct {
	ID   string `bson:"id" json:"id"`
	Text string `bson:"text" json:"text"`
}

type Poll struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Question  string        `bson:"question" json:"question"`
	Options   []PollOption  `bson:"options" json:"options"`
	CreatedBy bson.ObjectID `bson:"createdBy" json:"createdBy"`
	ShareCode string        `bson:"shareCode" json:"shareCode"`
	IsClosed  bool          `bson:"isClosed" json:"isClosed"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
}