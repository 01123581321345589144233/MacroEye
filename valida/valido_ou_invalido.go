package valida

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dreamscached/minequery/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientOptions = options.Client().ApplyURI("mongodb://127.0.0.1:27017/")
	client, err   = mongo.Connect(context.TODO(), clientOptions)
)

type Ip struct {
	Ip string `bson:"Ip"`
}

func Valida(pais string, ping time.Duration) {
	// Access a MongoDB collection through a database
	col := client.Database("Bigeye").Collection(pais)
	filtro := bson.D{{"Valido", false}}
	//filtro := bson.M{}
	cursor, err := col.Find(context.TODO(), filtro)
	if err != nil {
		os.Exit(1)
	}

	pinger := minequery.NewPinger(
		minequery.WithTimeout(ping * time.Millisecond),
	)

	for cursor.Next(context.TODO()) {
		var result Ip
		cursor.Decode(&result)

		_, err := pinger.Ping17(result.Ip, 25565)
		switch err {
		case nil:
			_, err := col.UpdateOne(
				context.Background(),
				bson.M{"Ip": result.Ip},
				bson.M{"$set": bson.M{"Valido": true}},
			)
			if err != nil {
				fmt.Print(err)
			}
			fmt.Println("Online status update em -->", result.Ip)
		default:
			_, err := col.UpdateOne(
				context.Background(),
				bson.M{"Ip": result.Ip},
				bson.M{"$set": bson.M{"Valido": false}},
			)
			if err != nil {
				fmt.Print(err)
			}
			// fmt.Println(result.Ip, "nunca foi visto online até hoje")
		}
	}
}

func Online(pais string, ping time.Duration) {
	// Access a MongoDB collection through a database
	col := client.Database("Bigeye").Collection(pais)
	filtro := bson.M{}
	cursor, err := col.Find(context.TODO(), filtro)
	if err != nil {
		os.Exit(1)
	}

	pinger := minequery.NewPinger(
		minequery.WithTimeout(ping * time.Millisecond),
	)

	for cursor.Next(context.TODO()) {
		var result Ip
		cursor.Decode(&result)

		_, err := pinger.Ping17(result.Ip, 25565)
		switch err {
		case nil:
			_, err := col.UpdateOne(
				context.Background(),
				bson.M{"Ip": result.Ip},
				bson.M{"$set": bson.M{"Online": true}},
			)
			if err != nil {
				fmt.Print(err)
			}
			fmt.Println("Online status update em -->", result.Ip)
		default:
			_, err := col.UpdateOne(
				context.Background(),
				bson.M{"Ip": result.Ip},
				bson.M{"$set": bson.M{"Online": false}},
			)
			if err != nil {
				fmt.Print(err)
			}
			// fmt.Println(result.Ip, "nunca foi visto online até hoje")
		}
	}
}
