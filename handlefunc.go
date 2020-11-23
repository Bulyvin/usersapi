package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, Im is company person servis!\n")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	databody := ""
	for true {
		bs := make([]byte, 1024)
		n, err := r.Body.Read(bs)
		databody = (databody + string(bs[:n]))
		if n == 0 || err != nil {
			break
		}
	}
	defer r.Body.Close()

	NewReq := new(User)
	if err := json.Unmarshal([]byte(databody), &NewReq); err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte("400- Bad Request not valid json "))
		return
	}
	C := GetClient()
	err := C.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Println("Couldn't connect to the database", err)
		w.WriteHeader(503)
		w.Write([]byte("503 Service Unavailable"))
		return
	}
	log.Println("Connected!")

	switch r.Method {

	case "GET":

		allusers := ReturnAllUsers(C, bson.M{"idcompany": NewReq.Idcompany})
		for _, user := range allusers {
			//log.Println(user.ID, user.Name, user.Surname, user.Phone, user.Passport, user.Idcompany)

			otvet, err := json.Marshal(user)
			if err != nil {
				log.Println("Ошибка маршалинга")
			}
			w.WriteHeader(200)
			io.WriteString(w, string(otvet))

		}

	case "POST":

		insertedID, err := InsertNewUser(C, *NewReq)
		if err != nil {
			w.WriteHeader(503)
			w.Write([]byte("503 Service Unavailable dont can create ID"))
			return
		}
		log.Println(insertedID.ID)
		otvet, err := json.Marshal(insertedID.ID)
		if err != nil {
			log.Println("Ошибка маршалинга")
		}
		w.WriteHeader(200)
		io.WriteString(w, "createdID="+string(otvet))

	case "PUT":

		userUpdated := UpdateUser(C, *NewReq, bson.M{"id": NewReq.ID})
		log.Println("Users updated count:", userUpdated)
		w.WriteHeader(200)
		w.Write([]byte("200 - update successful"))

	case "DELETE":

		userRemoved := RemoveOneUser(C, bson.M{"id": NewReq.ID})
		log.Println("Users removed count:", userRemoved)
		w.WriteHeader(200)
		w.Write([]byte("200 - deleted"))

	default:
		log.Println("Err func processor")
		w.WriteHeader(405)
		w.Write([]byte("405 Method Not Allowed"))
	}
}
