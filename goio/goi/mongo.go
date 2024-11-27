package goi

import (
	gomongo "github.com/gif-gif/go.io/go-db/go-mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

func Mongo(names ...string) *mongo.Database {
	client := MongoClient(names...)
	if client == nil {
		return nil
	}
	return client.DB()
}

func MongoClient(names ...string) *gomongo.GoMongo {
	return gomongo.GetClient(names...)
}
