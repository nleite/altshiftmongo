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

func main() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal("Make sure your ``mongod`` is up and running on port **27017**")
		panic(err)
	}
	defer session.Close()

	fmt.Println("Getting a database object")
	db := session.DB("altshiftmongo")
	collection := db.C("arrays")
	n, err := collection.Count()
	if err != nil {
		log.Fatal("Could not count collection")
	}
	fmt.Printf("Connected to arrays. Current count: %d\n", n)

	doc1 := SomeArray{[]int{1, 2, 3, 4}}
	insertErr := collection.Insert(&doc1)
	if insertErr != nil {
		log.Fatal(err)
	}

}
