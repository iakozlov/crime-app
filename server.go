package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var Port = ":5555"
var memoryCache = cache.New(5*time.Minute, 10*time.Minute)
var client *mongo.Client
var tripsCollection *mongo.Collection
var ctx context.Context

func main() {

	//initDatabase()
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
func UserEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "POST":
		r.ParseMultipartForm(0)
		username := r.FormValue("tripsMessage")
		cursor, err := tripsCollection.Find(ctx, bson.M{"username": username})
		if err != nil {
			log.Fatal(err)
		}
		var trips []bson.M
		if err = cursor.All(ctx, &trips); err != nil {
			log.Fatal(err)
		}
		fmt.Println(fmt.Fprintf(w, "%v", trips))
	}
}
func ServeFiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		path := r.URL.Path
		if path == "/" {
			path = "./index.html"
		} else {
			path = "./" + path
		}
		http.ServeFile(w, r, path)
	case "POST":
		r.ParseMultipartForm(0)
		message := r.FormValue("message")
		fmt.Println("Message from client: ", message)
		lat := strings.Split(message, ";")[0]
		lng := strings.Split(message, ";")[1]
		date := strings.Split(message, ";")[2]
		hour := strings.Split(message, ";")[3]
		username := strings.Split(message, ";")[4]
		address := strings.Split(message, ";")[5]
		var answer string
		date += " " + hour
		_, flag := memoryCache.Get(message)
		if !flag {
			cmd := exec.Command("python",
				"parse_info.py",
				lng, lat, date, hour)
			//stderr, err := cmd.StderrPipe()
			err := cmd.Start()
			if err != nil {
				panic(err)
			}

			//go copyOutput(stderr)
			cmd.Wait()
			copyOutput(&answer)

			memoryCache.Set(message, answer, 5*time.Minute)
		} else {
			tmp, _ := memoryCache.Get(message)
			answer = fmt.Sprintf("%v", tmp)
		}
		go writeToCSv(answer)
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
		_, err := tripsCollection.UpdateMany(
			ctx,
			bson.M{"username": username},
			bson.D{
				{"$push", bson.D{{"trips", address + ":" + answer}}},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", answer)

	case "PUT":
		r.ParseMultipartForm(0)
		message := r.FormValue("tripsMessage")
		fmt.Println("Message from put client: ", message)
		fmt.Fprintf(w, "%s\n", message)
	default:
		fmt.Fprintf(w, "Request type other than get or post")
	}
}

func copyOutput(str *string) {
	dat, _ := ioutil.ReadFile("info.txt")
	*str = string(dat)
}
func writeToCSv(str string) {
	var data [][]string
	var slice []string
	slice = append(slice, "date", "value")
	data = append(data, slice)
	crimes := strings.Split(str, ";")
	for i := 0; i < 3; i++ {
		var slice []string
		arr := strings.Split(crimes[i], ":")
		slice = append(slice, arr[0], arr[1])
		data = append(data, slice)
	}
	file, err := os.Create("XYZ.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
}
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
func isUserInDatabase(username string) bool{
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
