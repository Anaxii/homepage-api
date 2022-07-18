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

type Token struct {
	Name         string `json:"name"`
	Token        string `json:"token"`
	Quote        string `json:"quote"`
	PairAddress  string `json:"pairAddress"`
	TokenAddress string `json:"tokenAddress"`
	QuoteAddress string `json:"quoteAddress"`
}

type Basket struct {
	TokenAddress string `json:"tokenAddress"`
	JoeAddress   string `json:"joeAddress"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Tokens       []Token `json:"tokens"`
}

var HomeStats interface{}
var Baskets []Basket
var Tokens []Token

func serve(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

func serveHomeStats(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(HomeStats)
	serve(w, data)
}

func serveBaskets(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(Baskets)
	serve(w, data)
}

func serveTokens(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(Tokens)
	serve(w, data)
}

func main() {
	go fetchTimer()
	router := mux.NewRouter()
	router.HandleFunc("/homestats", serveHomeStats)
	router.HandleFunc("/baskets", serveBaskets)
	router.HandleFunc("/tokens", serveTokens)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func fetchTimer() {
	getHomeStats()
	getBaskets()
	getTokens()
	s := gocron.NewScheduler()
	s.Every(60).Seconds().Do(getHomeStats)
	s.Every(60).Seconds().Do(getBaskets)
	s.Every(60).Seconds().Do(getTokens)
	<- s.Start()
}

func getHomeStats() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://anaxii:B4ngFestina123@cluster0.1lvz50h.mongodb.net/test"))
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
		}
	}()
	coll := client.Database("Exposure").Collection("exposurelandingstats")
	var result map[string]interface{}
	_options := options.FindOne()
	_options.SetSort(bson.M{"$natural":-1})
	err = coll.FindOne(context.TODO(), bson.M{}, _options).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return err
	}
	delete(result, "_id")
	HomeStats = result
	return nil
}

func getBaskets() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://anaxii:B4ngFestina123@cluster0.1lvz50h.mongodb.net/test"))
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("Exposure").Collection("baskets")
	var result []Basket
	cursor, err := coll.Find(context.TODO(), bson.D{})
	ctx := context.Background()
	if err = cursor.All(ctx, &result); err != nil {
		return err
	}
	Baskets = result
	return nil
}

func getTokens() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://anaxii:B4ngFestina123@cluster0.1lvz50h.mongodb.net/test"))
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("Exposure").Collection("tokens")
	var result []Token
	cursor, err := coll.Find(context.TODO(), bson.D{})
	ctx := context.Background()
	if err = cursor.All(ctx, &result); err != nil {
		return err
	}
	Tokens = result
	return nil
}
