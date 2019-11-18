package main

import (
	"encoding/json"
	"fmt"
	"goapp1/controllers"
	u "goapp1/util"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

func main() {

	env := godotenv.Load("..\\src\\goapp1\\.env")
	if env != nil {
		fmt.Println(env)
	}

	router := mux.NewRouter()
	port := "8000"

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Api functioning")
	})

	router.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Hello there")
	})

	router.HandleFunc("/api/thing1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Thing1 rep")

	})

	router.HandleFunc("/api/thing2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		u.Message("success", "You hit an endpoint")
		u.Response(w, u.Message("success", "You hit an endpoint"))

	})

	router.HandleFunc("/api/thing31", controllers.GetMessage()).Methods("GET")
	router.HandleFunc("/api/Login", controllers.Login()).Methods("GET")

	router.Use(controllers.JwtMiddleware)

	err := http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, router))

	if err != nil {
		fmt.Print(err)
	}

}
