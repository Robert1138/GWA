package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	lg "github.com/Robert1138/GWA/util/log"
	"github.com/joho/godotenv"
)

// Planning for serverConfig and start server to handle everything from the .env file
func loadConfig() (map[string]string, error) {
	configPath, err := filepath.Abs("../src/github.com/Robert1138/GWA/.env")
	if err != nil {
		return map[string]string{}, err
	}

	err = godotenv.Load(configPath)
	if err != nil {
		return map[string]string{}, err
	}

	serverConfig, err := godotenv.Read(configPath)
	if err != nil {
		fmt.Println("in server LoadConfig Read()")
		fmt.Println(err)
		return map[string]string{}, err
	}
	return serverConfig, err
}

// StartServer will start the server will details specified by env
func StartServer(handler *http.Handler) {
	// if loadConfig fails set a default config.  TODO either add setting default config or make it fail entirely
	serverConfig, err := loadConfig()
	if err != nil {
		fmt.Println("in server startServer")
		fmt.Println(err)
	}

	cfgReadTimeout, errRT := time.ParseDuration(serverConfig["readtimeout"])
	if errRT != nil {
		cfgReadTimeout = 10
	}

	cfgWriteTimeout, errWT := time.ParseDuration(serverConfig["writetimeout"])
	if errWT != nil {
		cfgReadTimeout = 10
	}

	newServer := &http.Server{
		Addr:         serverConfig["port"],
		ReadTimeout:  cfgReadTimeout,
		WriteTimeout: cfgWriteTimeout,
		Handler:      *handler,
	}

	/*
		The server gets started in a go function while a channel is set up to listen to any signals that
		might shutdown this process.  When the channel gets a signal we create the initial context using
		the top-level context (context.Background()) and cancels when the timeout period has ended.
		At this point we do some clean up and shutdown the server

	*/
	go func() { // this could be switched with listening for signals portion in the following after this block
		err = newServer.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	signalChannel := make(chan os.Signal, 4) //buffer for each of the signals.  Might only need one but not sure
	// for windows it appears killing the process via any means other than "ctrl-c" doesnt trigger SIGTERM, SIGKILL, SIGINT
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	<-signalChannel

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel() // this cancels the parent when we finish with the following
	lg.CloseOutputFile()

	if err = newServer.Shutdown(ctx); err != nil {
		fmt.Println("error during shutdown")
	} else {
		//lg.CloseOutputFile()
		fmt.Println("no shutdown errors")
	}

}
