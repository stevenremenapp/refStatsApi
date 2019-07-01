package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// dbname, exists := os.LookupEnv("DB_NAME")
// 	if exists {
// 		fmt.Println(dbname)
// 	}

// DB_HOST=127.0.0.1
// DB_PORT=5432
// DB_USER=postgres
// DB_PASS=Octopus2017!
// DB_NAME=refstats

var dbHost = os.Getenv("HOST")
var dbPort = os.Getenv("DBPORT")
var dbUser = os.Getenv("USER")
var dbPassword = os.Getenv("PASSWORD")
var dbName = os.Getenv("NAME")
var serverPort = os.Getenv("SERVERPORT")

var (
	host       = dbHost
	dbport     = dbPort
	user       = dbUser
	password   = dbPassword
	dbname     = dbName
	serverport = serverPort
)

type Interaction struct {
	// gorm.Model

	// ID        int
	// Type      string
	// Timestamp string

	ID        int    `json:"id"`
	Type      string `json:"type"`
	Timestamp string `json:"time"`
}

type Interactions []Interaction

func allInteractions(w http.ResponseWriter, r *http.Request) {
	// interactions := Interactions{
	// 	Interaction{
	// 		ID:        1,
	// 		Type:      "tech",
	// 		Timestamp: time.Now().Format(time.Kitchen),
	// 	},
	// 	Interaction{
	// 		ID:        2,
	// 		Type:      "reference",
	// 		Timestamp: time.Now().Add(time.Hour * 5).Format(time.Kitchen),
	// 	},
	// 	Interaction{
	// 		ID:        3,
	// 		Type:      "tech",
	// 		Timestamp: time.Now().Add(time.Hour * 2).Format(time.Kitchen),
	// 	},
	// }

	var interactions []Interaction
	//w.Header().Set("Content-Type", "application/json")
	db.Find(&interactions)

	// fmt.Println("Hit all interactions endpoint")
	json.NewEncoder(w).Encode(&interactions)

	// dbHost, exists := os.LookupEnv("dbHost")
	// dbPort, exists := os.LookupEnv("dbPort")
	// if exists {
	// 	fmt.Println(dbHost)
	// 	fmt.Println(dbPort)
	// }

	// dbHost := os.Getenv("dbHost")
	// fmt.Println(dbHost)
}

func postInteraction(w http.ResponseWriter, r *http.Request) {
	var interaction Interaction
	json.NewDecoder(r.Body).Decode(&interaction)
	db.Create(&interaction)
	json.NewEncoder(w).Encode(&interaction)
}

func deleteInteraction(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	// w.Header().Set("Content-Type", "application/json")
	// fmt.Println("Delete endpoint hit")
	params := mux.Vars(r)
	var interaction Interaction
	db.First(&interaction, params["id"])
	db.Delete(&interaction)

	var interactions []Interaction
	db.Find(&interactions)
	json.NewEncoder(w).Encode(&interactions)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "HomePage!")
}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/interactions", allInteractions).Methods("GET")
	myRouter.HandleFunc("/interactions", postInteraction).Methods("POST")
	myRouter.HandleFunc("/interactions/{id}", deleteInteraction).Methods("DELETE", "OPTIONS")
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	})
	handler := c.Handler(myRouter)
	http.ListenAndServe(fmt.Sprintf(":%s", serverport), handler)
}

var db *gorm.DB
var err error

func main() {

	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)

	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }

	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	db, err = gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			host, dbport, user, dbname, password),
	)

	// db, err = gorm.Open(
	// 	"postgres",
	// 	fmt.Sprintf(
	// 		"host="+os.Getenv("HOST")+" port="+os.Getenv("PORT")+" user="+os.Getenv("USER")+
	// 			" dbname="+os.Getenv("NAME")+" sslmode=disable password="+
	// 			os.Getenv("PASSWORD")),
	// )

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	handleRequests()
}
