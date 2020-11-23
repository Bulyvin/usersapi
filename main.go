package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User struct represent a User
type User struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Surname   string      `json:"surname"`
	Idcompany int         `json:"idcompany"`
	Phone     string      `json:"phone"`
	Passport  PassportObg `json:"passport"`
}

//PassportObg struct
type PassportObg struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

//GetClient returns a MongoDB Client
func GetClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func main() {
	http.HandleFunc("/users/", userHandler)
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":4444", nil)
}

// ReturnAllUsers return all documents from the collection Users
func ReturnAllUsers(client *mongo.Client, filter bson.M) []*User {
	var allusers []*User
	collection := client.Database("userstest").Collection("users")
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}
	for cur.Next(context.TODO()) {
		var user User
		err = cur.Decode(&user)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		allusers = append(allusers, &user)
	}
	return allusers
}

// ReturnOneUser just one document from the collection Users
func ReturnOneUser(client *mongo.Client, filter bson.M) User {
	var user User
	collection := client.Database("userstest").Collection("users")
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&user)
	return user
}

// InsertNewUser insert a new User in the users
func InsertNewUser(client *mongo.Client, user User) (User, error) {
	collection := client.Database("userstest").Collection("users")
	count := 0
LOOP:
	count = count + 1
	user.ID = randomID(10, 1000)
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("Error on inserting new User", err)
		if count == 5 {
			log.Println("Error Index create")
			return user, err
		}
		goto LOOP
	}
	log.Println("CreatNewUserID", insertResult)
	return user, nil
}

// RemoveOneUser remove one existing User
func RemoveOneUser(client *mongo.Client, filter bson.M) int64 {
	collection := client.Database("userstest").Collection("users")
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on deleting one User", err)
	}
	return deleteResult.DeletedCount
}

// UpdateUser update the info of a informed User
func UpdateUser(client *mongo.Client, userforupdate User, filter bson.M) int64 {
	collection := client.Database("userstest").Collection("users")
	var user User
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&user)
	log.Println("Клиент для исправления=", &user)
	if userforupdate.Idcompany != 0 {
		user.Idcompany = userforupdate.Idcompany
	}
	if userforupdate.Name != "" {
		user.Name = userforupdate.Name
	}
	if userforupdate.Surname != "" {
		user.Name = userforupdate.Surname
	}
	if userforupdate.Phone != "" {
		user.Phone = userforupdate.Phone
	}
	if userforupdate.Passport.Type != "" {
		user.Passport.Type = userforupdate.Passport.Type
	}
	if userforupdate.Passport.Number != "" {
		user.Passport.Number = userforupdate.Passport.Number
	}
	log.Println("Исправленный клиент=", user)
	atualizacao := bson.D{{Key: "$set", Value: user}}
	updatedResult, err := collection.UpdateOne(context.TODO(), filter, atualizacao)
	if err != nil {
		log.Fatal("Error on updating one User", err)
	}
	return updatedResult.ModifiedCount
}

func randomID(min, max int) int {
	rand.Seed(time.Now().Unix())
	if min > max {
		return min
	}
	return rand.Intn(max-min) + min
}
