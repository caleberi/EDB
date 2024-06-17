package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sender struct {
	Name     string `json:"name"`
	Country  string `json:"country"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	DOB      string `json:"dob"`
	Email    string `json:"email"`
	IDNumber string `json:"idNumber"`
	IDType   string `json:"idType"`
}

type Destination struct {
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	AccountType   string `json:"accountType"`
	NetworkID     string `json:"networkId"`
}

type Payment struct {
	ID              string      `json:"id"`
	ChannelID       string      `json:"channelId"`
	SequenceID      string      `json:"sequenceId"`
	Currency        string      `json:"currency"`
	Country         string      `json:"country"`
	Amount          float64     `json:"amount"`
	Reason          string      `json:"reason"`
	ConvertedAmount float64     `json:"convertedAmount"`
	Status          string      `json:"status"`
	Rate            float64     `json:"rate"`
	Sender          Sender      `json:"sender"`
	Destination     Destination `json:"destination"`
	CreatedAt       string      `json:"createdAt"`
	UpdatedAt       string      `json:"updatedAt"`
	ExpiresAt       string      `json:"expiresAt"`
}

type Disbursement struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" validate:"required"`
	ReceiverID   primitive.ObjectID `bson:"receiver_id,omitempty" json:"receiver_id,omitempty" validate:"required"`
	CreatedAt    *time.Time         `bson:"createdAt,omitempty" json:"-" validate:"required"`
	UpdatedAt    *time.Time         `bson:"updatedAt,omitempty" json:"-" validate:"required"`
	SalaryAmount float64            `bson:"salary_amount,omitempty" json:"salary_amount,omitempty" validate:"required"`
	SenderID     primitive.ObjectID `bson:"sender_id,omitempty" json:"sender_id,omitempty" validate:"required"`
	Status       string             `bson:"status,omitempty" json:"status,omitempty" validate:"required"`
	Payment      Payment            `bson:"payment,omitempty" json:"payment,omitempty" validate:"required"`
}
