package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var Port = ":5555"
var memoryCache = cache.New(5*time.Minute, 10*time.Minute)
var client *mongo.Client
var tripsCollection *mongo.Collection
var ctx context.Context

//struct describing request to get crime analysis
type Request struct {
	Lat      string `json:"lat"`
	Lng      string `json:"lng"`
	Date     string `json:"date"`
	Time     string `json:"time"`
	Username string `json:"username"`
	Address  string `json:"address"`
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI(
		"mongodb+srv://ivankozlov:crime12345@cluster0.4o5vr.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	quickstartDatabase := client.Database("quickstartNew")
	tripsCollection = quickstartDatabase.Collection("trips-san-francisco")

	http.HandleFunc("/", ServeFiles)
	http.HandleFunc("/reg", UserEndpoint)
	fmt.Println("Serving @ : ", "127.0.0.1"+Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}

//func that works on /reg endpoint
//serves get request for history of requests
func UserEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		username := r.URL.Query().Get("username")
		cursor, err := tripsCollection.Find(ctx, bson.M{"username": username})
		if err != nil {
			log.Fatal(err)
		}
		var trips []bson.M
		if err = cursor.All(ctx, &trips); err != nil {
			log.Fatal(err)
		}
		jsonTrips, err := json.Marshal(trips)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%v", string(jsonTrips))
	}
}

//func serving main endpoint
//server post method for getting crime analysis
func ServeFiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := r.URL.Path
		if path == "/" {
			path = "../Frontend/index.html"
		} else {
			path = "../Frontend/" + path
		}
		http.ServeFile(w, r, path)
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("error in reading body")
		}
		fmt.Printf("%s", body)
		var request Request
		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Fatal("error here")
		}
		fmt.Println("Message from client: ", request.Lat)
		lat := request.Lat
		lng := request.Lng
		date := request.Date
		hour := request.Time
		username := request.Username
		address := request.Address
		var answer string
		date += " " + hour
		_, flag := memoryCache.Get(request.Lat + ";" + request.Lng + ";" + request.Date + ";" + request.Time)
		if !flag {
			cmd := exec.Command("python",
				"parse_info.py",
				lng, lat, date, hour)
			err = cmd.Start()
			if err != nil {
				panic(err)
			}
			cmd.Wait()
			copyOutput(&answer)
			memoryCache.Set(request.Lat+";"+request.Lng+";"+request.Date+";"+request.Time, answer, 5*time.Minute)
		} else {
			tmp, _ := memoryCache.Get(request.Lat + ";" + request.Lng + ";" + request.Date + ";" + request.Time)
			answer = fmt.Sprintf("%v", tmp)
		}
		if username != "" {
			go updateUserRequests(username, address, answer)
		}
		answer = strings.ReplaceAll(answer, "\n", ";")
		jsonAnswer, err := json.Marshal(answer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", string(jsonAnswer))

	default:
		fmt.Fprintf(w, "Request type other than get or post")
	}
}

//func to read data from model output
func copyOutput(str *string) {
	dat, _ := ioutil.ReadFile("info.txt")
	*str = string(dat)
}

//func to check if user is in database
func isUserInDatabase(username string) bool {
	cursor, err := tripsCollection.Find(ctx, bson.M{"username": username})
	if err != nil {
		log.Fatal(err)
	}
	var trips []bson.M
	if err = cursor.All(ctx, &trips); err != nil {
		log.Fatal(err)
	}
	return len(trips) > 0
}
func updateUserRequests(username string, address string, answer string) {
	userFlag := isUserInDatabase(username)
	if !userFlag {
		_, err := tripsCollection.InsertOne(ctx, bson.D{
			{Key: "username", Value: username},
			{Key: "trips", Value: bson.A{}},
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	ans := make(chan string)
	go modifyAnswer(answer, ans)
	_, err := tripsCollection.UpdateMany(
		ctx,
		bson.M{"username": username},
		bson.D{
			{"$push", bson.D{{"trips", address + ":" + <-ans}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

//func that changes answer from float to percent
func modifyAnswer(answer string, res chan string) {
	result := ""
	tmp := answer[:len(answer)-1]
	crimes := strings.Split(tmp, "\n")
	for i := 0; i < 3; i++ {
		arr := strings.Split(crimes[i], ":")
		result += arr[0]
		result += " - "
		probability, _ := strconv.ParseFloat(strings.TrimSpace(arr[1]), 32)
		percentage := int(math.Round(probability) * 100)
		result += strconv.Itoa(percentage)
		result += "%;"
	}
	res <- result
}
