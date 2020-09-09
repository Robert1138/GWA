package misc

import (
	"encoding/json"
	"fmt"
	"goapp1/pkg/item"
	"goapp1/util"
	lg "goapp1/util/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var log *lg.StandardLogger

// HTTPRoutes registers endpoints and their appropriate HandlerFuncs to the provided router as well as any subrouters.
func HTTPRoutes(router *mux.Router, newLogger *lg.StandardLogger) {
	log = newLogger
	router.HandleFunc("/hello", Hello()).Methods("GET")
	router.HandleFunc("/time", Time()).Methods("GET")
	router.HandleFunc("/pic", Pic()).Methods("GET")
}

// Hello lets you say hello to Obi-Wan Kenobi
func Hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		//json.NewEncoder(w).Encode(util.Message("200", "Hello there"))
		util.Response(w, util.Message("200", "Hello there"))
	}
}

// C:\Users\Robert\go\src\goapp1\file1.jpg
func Pic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("..\\src\\goapp1\\file2.jpg")
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		}

		FileHeader := make([]byte, 512)
		file.Read(FileHeader)

		FileContentType := http.DetectContentType(FileHeader)
		fmt.Println(FileContentType)
		FileStat, err := file.Stat()
		if err != nil {
			fmt.Println(err)
		}
		FileSize := FileStat.Size()
		w.Header().Set("Content-Disposition", "inline") // set headers after reading the files !!!!!!
		w.Header().Set("Content-Type", FileContentType)
		w.Header().Set("Conent-Length", strconv.FormatInt(FileSize, 10))
		file.Seek(0, 0)
		io.Copy(w, file)
		return
	}
}

// Time displays time
func Time() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// w.WriteHeader(http.StatusFound)
		item.GetItem(1)
		util.Response(w, util.Message("200", strconv.FormatInt(time.Now().Unix(), 10)))
	}
}

// CsrfPostTest tests post with csrf protection
func CsrfPostTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("POST WENT THROUGH WITH VALID CSRF")

	}
}

// GetMessage returns a simple response -- placeholder to get csrf token
func GetMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.Response(w, util.Message("success", "it worked"))
	}
}
