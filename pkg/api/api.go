package api

import (
	"fmt"
	"os"

	"github.com/Robert1138/GWA/pkg/auth"
	"github.com/Robert1138/GWA/pkg/csrftest"
	"github.com/Robert1138/GWA/pkg/item"
	"github.com/Robert1138/GWA/pkg/middleware"
	"github.com/Robert1138/GWA/pkg/misc"
	lg "github.com/Robert1138/GWA/util/log"
	"github.com/Robert1138/GWA/util/server"

	"github.com/Robert1138/GWA/pkg/account"

	_ "github.com/Robert1138/GWA/util/log"
	"github.com/rs/cors"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func testFileCloser(newFile *os.File) {
	newFile.Close()
}

// StartAPI is where all the logger, route/handlerfunc functions, middleware and serverStart gets called
func StartAPI() {

	// logger is passed to HTTPRoutes and then used in file
	// Instead of defer, the logging output file is closed on server shutdown
	logger := lg.NewLogger()
	logger.Info("logger in api is working")
	router := mux.NewRouter()
	fmt.Println("Router setup")

	csrfMiddleware := csrf.Protect([]byte(os.Getenv("csrfkey")), csrf.Secure(false)) //
	// At this point endpoints are address:port/example    If it doesnt need auth it goes here
	misc.HTTPRoutes(router, logger)
	item.HTTPRoutes(router, logger)
	auth.HTTPRoutes(router, logger) // csrf middleware not intended to be attached
	// At this point endpoints are address:port/v1/whatever If it needs auth it goes here
	api := router.PathPrefix("/v1").Subrouter()
	account.HTTPRoutes(api, logger)
	csrftest.HTTPRoutes(api, logger)

	api.Use(middleware.JwtMiddleware)
	api.Use(csrfMiddleware)
	api.Use(middleware.CsrfTokenMiddleware) // sets csrf token in header for all get request

	logHandler := handlers.LoggingHandler(os.Stdout, router)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "https://auction-e5422.web.app"},
		AllowCredentials: true,
		Debug:            true, // Debug true for testing, disable for production

	})
	corsHandler := c.Handler(logHandler)

	fmt.Println("Starting server")
	server.StartServer(&corsHandler)
	//fmt.Println("got here after server closed")

}
