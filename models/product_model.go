package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	SKU      string             `json:"sku,omitempty" validate:"required"`
	Quantity int                `json:"quantity,omitempty" validate:"required"`
}
