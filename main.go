package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
)

// Article - Our struct for all articles
type Article struct {
	Id      int    `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Articles []Article

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	articles := Articles{
		Article{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	fmt.Println("Endpoint Hit: returnAll")

	json.NewEncoder(w).Encode(articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Fprintf(w, "Key: "+key)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/create", createTable)
	myRouter.HandleFunc("/all", returnAllArticles)
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func createTable(w http.ResponseWriter, r *http.Request) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// createTable()

	params := &dynamodb.CreateTableInput{
		TableName: aws.String("Articles"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("Id"), KeyType: aws.String("HASH")},
			{AttributeName: aws.String("Title"), KeyType: aws.String("RANGE")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			// {AttributeName: aws.String("Desc"), AttributeType: aws.String("S")},
			// {AttributeName: aws.String("Content"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("Title"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("Id"), AttributeType: aws.String("S")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(100),
		},
	}

	resp, err := svc.CreateTable(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// print the response data
	fmt.Println(resp)
}

func main() {
	handleRequests()
	// createTable()
}
