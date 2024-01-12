package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
	Mongo *mongo.Client
}