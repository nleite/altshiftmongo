package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
)

// SomeArray first document array example
type SomeArray struct {
	SomeArray []int `json:"some_array" bson:"some_array"`
}

// StringArray - different document type to be inserted on arrays collection
type StringArray struct {
	StringArray []string `json:"string_array" bson:"string_array"`
}

func Count(collection *mgo.Collection){
	n, err := collection.Count()
	if err != nil {
		log.Fatal("Could not count collection")
	}
	fmt.Printf("Connected to arrays. Current count: %d\n", n)
}


func Drop(collection *mgo.Collection){
	collection.DropCollection()
}


func InsertStringArray(collection *mgo.Collection){
	doc2 := StringArray{[]string{"bernie", "ernie", "dottie"}}
	insertErr := collection.Insert(&doc2)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
}


func FindMathingString(collection *mgo.Collection, matching string) {
	filter := bson.M{"some_array": matching}
	results := []bson.D{}

}

func main() {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		log.Fatal("Make sure your ``mongod`` is up and running on port **27017**")
		panic(err)
	}
	defer session.Close()

	fmt.Println("Getting a database object")
	db := session.DB("altshiftmongo")
	collection := db.C("arrays")
	Count(collection)

	doc1 := SomeArray{[]int{1, 2, 3, 4}}
	insertErr := collection.Insert(&doc1)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	Count(collection)

	InsertStringArray(collection)
	Count(collection)
	//Drop(collection)
}
