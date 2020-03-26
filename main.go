package main

import (
	"fmt"
	"goapp1/controllers"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/csrf"
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

	f, err1 := os.OpenFile("..\\src\\goapp1\\info.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err1 != nil {
		fmt.Println("log failed")
	}
	defer f.Close()
	logrus.SetOutput(f)
	// log can be set to json
	// logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.Info("it works")
	router := mux.NewRouter()
	port := "8080"

	fmt.Println("Router setup")
	csrfMiddleware := csrf.Protect([]byte(os.Getenv("csrfkey")), csrf.Secure(false))

	router.HandleFunc("/hello", controllers.Hello()).Methods("GET")
	router.HandleFunc("/Login", controllers.Login()).Methods("POST")       // csrf middleware not intended to be attached
	router.HandleFunc("/Register", controllers.Register()).Methods("POST") // csrf middleware not intended  to be attached

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/stuff", controllers.GetMessage()).Methods("GET")   // GET placeholder to get csrf token
	api.HandleFunc("/CSRF", controllers.CsrfPostTest()).Methods("POST") // POST placeholder that requires a csrf token in the header

	api.Use(csrfMiddleware)
	api.Use(controllers.CsrfTokenMiddleware) // sets csrf token in header for all get request
	router.Use(controllers.JwtMiddleware)

	fmt.Println("Starting server")
	err := http.ListenAndServe("localhost:"+port, handlers.LoggingHandler(os.Stdout, router))

	if err != nil {
		fmt.Print(err)
	}

}
