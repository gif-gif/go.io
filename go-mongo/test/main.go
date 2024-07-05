package main

import (
	"context"
	golog "github.com/gif-gif/go.io/go-log"
	gomongo "github.com/gif-gif/go.io/go-mongo"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

func main() {
	config := gomongo.Config{
		Name:     "",
		Addr:     "127.0.0.1:27017",
		User:     "",
		Password: "",
		Database: "",
		AutoPing: true,
	}

	client, err := gomongo.New(config)
	if err != nil {
		golog.ErrorF("error:%v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.DB().Collection("users")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 500)
	golog.InfoF("end of gomongo")
}
