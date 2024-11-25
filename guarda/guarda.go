package guarda

import (
	"context"
	"fmt"
	"strconv"
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

type Status17 struct {
	Ip            string                    `bson:"Ip"`
	SamplePlayers []minequery.PlayerEntry17 `bson:"PlayersOnline"`
	Dia           string                    `bson:"Dia"`
	Hora          string                    `bson:"Hora"`
}

type Database struct {
	Target      string `bson:"Ip"`
	Version     string `bson:"Version"`
	Players     string `bson:"Players"`
	Description string `bson:"Motd"`
	Dia         string `bson:"Dia"`
	Hora        string `bson:"Hora"`
}

func Guarda(ping time.Duration, pais string) {
	col := client.Database("Bigeye").Collection(pais)
	filtro := bson.D{{"Valido", true}}
	//filtro := bson.M{}
	cursor, err := col.Find(context.TODO(), filtro)
	if err != nil {
		fmt.Print(err)
	}

	pinger := minequery.NewPinger(
		minequery.WithTimeout(ping * time.Millisecond),
	)
	for cursor.Next(context.TODO()) {
		var result Ip
		cursor.Decode(&result)
		res, err := pinger.Ping17(result.Ip, 25565)
		if err != nil {
			//fmt.Println(err)
		} else {
			//fmt.Println("Server com ", res.OnlinePlayers, "\n")

			if res.OnlinePlayers > 0 {
				collection2 := client.Database("PlayersArchive").Collection("PT")
				dt := time.Now()

				final := Status17{
					result.Ip,
					append(res.SamplePlayers),
					dt.Format("01-02-2006"),
					dt.Format("15:04:05"),
				}

				fmt.Println(res.SamplePlayers)
				_, err = collection2.InsertOne(context.TODO(), final)
				if err != nil {
					//fmt.Println("Duplicado.")
				}
				fmt.Println("Servidor com pessoas online <-------------------")
				collection4 := client.Database("Servers").Collection(pais)
				final2 := Database{
					result.Ip,
					res.VersionName,
					strconv.Itoa(res.OnlinePlayers) + "/" + strconv.Itoa(res.MaxPlayers),
					res.Description.String(),
					dt.Format("01-02-2006"),
					dt.Format("15:04:05"),
				}

				opts := options.Update().SetUpsert(true)
				_, err = collection4.UpdateOne(context.TODO(), final2, opts)
				if err != nil {
					//fmt.Println("Duplicado.")
				}
			} else {
				collection2 := client.Database("Servers").Collection(pais)
				dt := time.Now()
				//opts := options.Update().SetUpsert(true)
				final := Database{
					result.Ip,
					res.VersionName,
					strconv.Itoa(res.OnlinePlayers) + "/" + strconv.Itoa(res.MaxPlayers),
					res.Description.String(),
					dt.Format("01-02-2006"),
					dt.Format("15:04:05"),
				}
				_, err = collection2.InsertOne(context.TODO(), final) //, opts)
				if err != nil {
					//fmt.Println("Duplicado.")
				}
			}
		}
	}
}
