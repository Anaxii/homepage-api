package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

var HomeStats interface{}

func serveHomeStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	data, _ := json.Marshal(HomeStats)
	w.Write(data)
}

func main() {
	go fetchTimer()
	router := mux.NewRouter()
	router.HandleFunc("/homestats", serveHomeStats)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func fetchTimer() {
	getHomeStats()
	s := gocron.NewScheduler()
	s.Every(60).Seconds().Do(getHomeStats)
	<- s.Start()
}

func getHomeStats() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://anaxii:B4ngFestina123@cluster0.1lvz50h.mongodb.net/test"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("Exposure").Collection("exposurelandingstats")
	var result map[string]interface{}
	_options := options.FindOne()
	_options.SetSort(bson.M{"$natural":-1})
	err = coll.FindOne(context.TODO(), bson.M{}, _options).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return
	}
	delete(result, "_id")
	HomeStats = result
}
