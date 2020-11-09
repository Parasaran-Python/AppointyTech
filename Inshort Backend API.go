package main

import (
    "net/http"
    "fmt"
    "log"
    "context"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	Content string
}

var client *mongo.Client

func writeToDB(data string){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.TODO(), nil)

    if err != nil {
        log.Fatal(err)
    }
    
    collection:=client.Database("tango").Collection("bingo")
    
    data1:=Article{data}
    
    collection.InsertOne(context.TODO(), data1)
    if err != nil {
        log.Fatal(err)
    }
    
}

func getDocById(id string){
	docID, err := primitive.ObjectIDFromHex(id)

	
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    
    if err != nil {
        log.Fatal(err)
    }
    
    result := Article{}
    
    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    col := client.Database("tango").Collection("bingo")
    
    err = col.FindOne(ctx, bson.M{"_id": docID}).Decode(&result)
	if err != nil {
		fmt.Println("FindOne() ObjectIDFromHex ERROR:", err)
	}
	
	return result.Content
	
}

func getAllArticles(){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    col := client.Database("tango").Collection("bingo")
	cursor, err := col.Find(context.TODO(), bson.D{})

    if err != nil {
        fmt.Println("Finding all documents ERROR:", err)
        defer cursor.Close(ctx)
	
	}
	
	for cursor.Next(ctx) {

            var result bson.M
            err := cursor.Decode(&result)

			ret:=""
            if err != nil {
                fmt.Println("cursor.Next() error:", err)
               
            } else {
                ret+=result["content"]+" "
            }
        }
        return ret
        //fmt.Println(cursor)
}

func main() {
	
	//getAllArticles()
	
	//getDocById("5fa97e8c535d0951acbcedf7")
	
	//writeToDB("Doc")
	
    http.HandleFunc("/articles", docWriter) // write to article to mongodb using post request or return all articles


    log.Fatal(http.ListenAndServe(":8080",nil))
}

func docWriter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "POST":
			if err := r.ParseForm(); err != nil {
            fmt.Fprintf(w, "ParseForm() err: %v", err)
            return
			}
			name := r.FormValue("doc")
			writeToDB(name)
		case "GET":
			str:=getAllArticles()
			fmt.Fprintf(w,str)
	}
}
