package main

import (
	"encoding/json"
	"fmt"
	"goapp1/controllers"
	u "goapp1/util"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, Newman")
	s := "hey now"
	s = strings.ToUpper(s)
	fmt.Println(s)

	env := godotenv.Load("..\\src\\goapp1\\.env")
	//C:\Users\Robert\go\src\goapp1
	if env != nil {
		fmt.Println(env)
	}

	dbName := os.Getenv("db_name")
	fmt.Println(dbName)
	router := mux.NewRouter()
	port := "8000"

	router.HandleFunc("/api/greating", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("welcome to the api")
	})

	router.HandleFunc("/api/thing1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("thing1 endpoint hit")
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(2)

	})

	router.HandleFunc("/api/thing2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("thing2 endpoint hit")
		w.Header().Add("Content-Type", "application/json")
		u.Message("success", "You hit an endpoint")
	})

	router.HandleFunc("/api/thing3", controllers.GetMessage()).Methods("GET")

	err := http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, router))

	if err != nil {
		fmt.Print(err)
	}

}
