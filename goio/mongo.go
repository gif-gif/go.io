package goio

import (
	gomongo "github.com/gif-gif/go.io/go-db/go-mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

func Mongo(names ...string) *mongo.Database {
	return gomongo.GetClient(names...).DB()
}

func MongoClient(names ...string) *gomongo.GoMongo {
	return gomongo.GetClient(names...)
}
