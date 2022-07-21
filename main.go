package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type TokenStruct struct {
	Name         string `json:"name"`
	Token        string `json:"token"`
	Quote        string `json:"quote"`
	PairAddress  string `json:"pairAddress"`
	TokenAddress string `json:"tokenAddress"`
	QuoteAddress string `json:"quoteAddress"`
}

type BasketStruct struct {
	TokenAddress string        `json:"tokenAddress"`
	JoeAddress   string        `json:"joeAddress"`
	Name         string        `json:"name"`
	Symbol       string        `json:"symbol"`
	Tokens       []TokenStruct `json:"tokens"`
}

type HoldersStruct struct {
	NumberOfHolders int                `json:"numberOfHolders"`
	HolderBalances  map[string]float64 `json:"holderBalances"`
	CurrentSupply   float64            `json:"currentSupply"`
	Block           int                `json:"block"`
	Time            int64              `json:"time"`
}

type PricesStruct struct {
	Index               float64            `json:"index"`
	PortionIndex        float64            `json:"portionIndex"`
	TrackedPortionIndex float64            `json:"trackedPortionIndex"`
	ActualIndex         float64            `json:"actualIndex"`
	TokenPrices         map[string]float64 `json:"tokenPrices"`
	BasketName          string             `json:"basketName"`
	Time                int64              `json:"time"`
}

var HomeStats interface{}
var Baskets []BasketStruct
var Tokens []TokenStruct
var Holders []HoldersStruct
var Prices []PricesStruct

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

func serveHolders(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(Holders)
	serve(w, data)
}

func servePrices(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(Prices)
	serve(w, data)
}

func main() {
	fetchData()
	go fetchTimer()
	router := mux.NewRouter()
	router.HandleFunc("/homestats", serveHomeStats)
	router.HandleFunc("/baskets", serveBaskets)
	router.HandleFunc("/tokens", serveTokens)
	router.HandleFunc("/holders", serveHolders)
	router.HandleFunc("/prices", servePrices)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func fetchTimer() {
	s := gocron.NewScheduler()
	s.Every(60).Seconds().Do(fetchData)
	<-s.Start()
}

func fetchData() {
	go getHomeStats()
	go getBaskets()
	go getTokens()
	go getHolders()
	go getPrices()
}

func getHomeStats() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://anaxii:B4ngFestina123@cluster0.1lvz50h.mongodb.net/test"))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())
	coll := client.Database("Exposure").Collection("exposurelandingstats")
	var result map[string]interface{}
	_options := options.FindOne()
	_options.SetSort(bson.M{"$natural": -1})
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
	defer client.Disconnect(context.TODO())

	coll := client.Database("Exposure").Collection("baskets")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	ctx := context.Background()
	if err = cursor.All(ctx, &Baskets); err != nil {
		return err
	}
	return nil
}

func getTokens() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://anaxii:B4ngFestina123@cluster0.1lvz50h.mongodb.net/test"))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	coll := client.Database("Exposure").Collection("tokens")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	ctx := context.Background()
	if err = cursor.All(ctx, &Tokens); err != nil {
		return err
	}
	return nil
}

func getHolders() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://anaxii:B4ngFestina123@cluster0.1lvz50h.mongodb.net/test"))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	coll := client.Database("Exposure").Collection("holders")
	var result []HoldersStruct
	cursor, err := coll.Find(context.TODO(), bson.D{})
	ctx := context.Background()
	if err = cursor.All(ctx, &result); err != nil {
		return err
	}

	var parsedResult []HoldersStruct
	lastSupply := 0.0
	for _, v := range result {
		if v.CurrentSupply == lastSupply || v.CurrentSupply == 0 {
			continue
		}
		lastSupply = v.CurrentSupply
		parsedResult = append(parsedResult, v)
	}
	Holders = parsedResult
	return nil
}

func getPrices() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://anaxii:B4ngFestina123@cluster0.1lvz50h.mongodb.net/test"))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	coll := client.Database("Exposure").Collection("prices")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	ctx := context.Background()
	if err = cursor.All(ctx, &Prices); err != nil {
		return err
	}
	return nil
}