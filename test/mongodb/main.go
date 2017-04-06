package main

import (
	"fmt"
	"log"
	"../../src/dependents/gopkg.in/mgo.v2"
	"../../src/dependents/gopkg.in/mgo.v2/bson"
	"server/utils/logger"
)

type Person struct {
	Name string
	Phone string
}

func main1() {
	session, err := mgo.Dial("192.168.18.88:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}

func main() {
	hello()
	fmt.Println("444444444")
	hello()
}

func hello(){
	defer func() {
		if err:=recover();err!=nil{
			logger.Error(nil, "StartAmqpConsumer error", err)
		}
	}()
	fmt.Println("hi")
	go func(){
		panic("false")
	}()
	fmt.Println("hi 2")
}