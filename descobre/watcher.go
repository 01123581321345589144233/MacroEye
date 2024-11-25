package watcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ip2location/ip2location-go"
	"github.com/zan8in/masscan"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	Ip        string  `bson:"Ip"`
	Dia       string  `bson:"Dia"`
	Hora      string  `bson:"Hora"`
	Pais      string  `bson:"Pais"`
	Valido    bool    `bson:"Valido"`
	Longitude float64 `bson:"Longitude"`
	Latitude  float64 `bson:"Latitude"`
}

func Watcher(rate int, ranges string, pais string) {
	var (
		scannerResult []masscan.ScannerResult
	)

	scanner, err := masscan.NewScanner(
		masscan.SetParamRate(rate),
		masscan.SetParamTargets(ranges),
		masscan.SetParamPorts("25565"),
		masscan.EnableDebug(),
	)

	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}

	if err := scanner.RunAsync(); err != nil {
		log.Printf("error async")
	}

	stdout := scanner.GetStdout()

	for stdout.Scan() {

		srs := masscan.ParseResult(stdout.Bytes())
		scannerResult = append(scannerResult, srs)
		dt := time.Now()

		fmt.Print(dt.Format("01-02-2006 15:04:05 Monday "), srs.IP, "\n")

		// MONGO DB MERDAS
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// get collection as ref
		collection := client.Database("Bigeye").Collection(pais)

		path := "path_of_ips"
		db, err := ip2location.OpenDB(path)
		results, err := db.Get_all(srs.IP)
		if err != nil {
			fmt.Print(err)
			return
		}
		fmt.Printf("longitude: %f\n", results.Longitude)
		fmt.Printf("latitude: %f\n", results.Latitude)

		lon := float64(results.Longitude)
		lat := float64(results.Latitude)

		dia := dt.Format("01-02-2006")
		hora := dt.Format("15:04:05")
		serverdb := Server{
			srs.IP,
			dia,
			hora,
			pais,
			false,
			lon,
			lat,
		}

		_, err = collection.InsertOne(context.TODO(), serverdb)
		if err != nil {
			fmt.Println("^----> Duplicado")
		}
		db.Close()
	}
}
