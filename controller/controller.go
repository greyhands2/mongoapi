package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/greyhands2/mongoapi/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//note i replaced "@" in the password with "%40"
const connectionString = "mongodb://localhost:27017"

const dbName = "netflix"
const colName = "watchlist"

//most important

var collection *mongo.Collection

//connect with mongodb
//an function runs only one time;at the beginning
func init() {
	//client option

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	collection = client.Database(dbName).Collection(colName)

	//colleciton instance
	fmt.Println("collection instance is ready")
}

//mongodb helpers-file

//insert 1 record
func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted one movie in db with id:", inserted.InsertedID)
}

//update onerecord
func updateOneMovie(movieId string) {

	id, _ := primitive.ObjectIDFromHex(movieId)

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count:", result.ModifiedCount)
}

//delete 1 record
func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)

	filter := bson.M{"_id": id}

	deleteCount, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("movie got deleted with delete count: ", deleteCount)
}

//delete all records

func deleteAllMovies() int64 {

	deleteResult, err := collection.DeleteMany(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of movies deleted", deleteResult.DeletedCount)

	return deleteResult.DeletedCount

}

//get all movies from db
func getAllMovies() []primitive.M {

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M

	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)

		if err != nil {
			log.Fatal(err)
		}

		movies = append(movies, movie)
	}
	defer cursor.Close(context.Background())
	return movies
}

//actual controllers-file

func GetAllMovies(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/x-www-form-urlencode")

	allMovies := getAllMovies()
	json.NewEncoder(res).Encode(allMovies)

}

func CreateMovie(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/x-www-form-urlencode")

	res.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(req.Body).Decode(&movie)

	insertOneMovie(movie)
	json.NewEncoder(res).Encode(movie)
}

func MarkAsWatched(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/x-www-form-urlencode")

	res.Header().Set("Allow-Control-Allow-Methods", "POST")

	params := mux.Vars(req)
	updateOneMovie(params["id"])

	json.NewEncoder(res).Encode(params["id"])
}

func DeleteAMovie(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	res.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(req)

	deleteOneMovie(params["id"])

	json.NewEncoder(res).Encode(params["id"])

}

func DeleteAllMovies(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	res.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deleteAllMovies()

	json.NewEncoder(res).Encode(count)

}
