package main

import (
	"csm/database"
	csmhttp "csm/http"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dbDriver := &database.DBDriver{}

	err := dbDriver.Connect()
	if err != nil {
		log.Printf("Failed to connect to MySQL database. Error %v", err)
		os.Exit(1)
	}
	defer dbDriver.Close()

	err = dbDriver.MaintainDatabasesAndTables()
	if err != nil {
		log.Printf("Database and table not initalized")
		os.Exit(1)
	}

	db := database.NewDatabase(dbDriver)
	go csmhttp.InitHttp(db)
	handleOsSignal()
}

func handleOsSignal() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	for {
		s := <-sigc
		log.Printf("OS signal received %s", s)
		log.Println("STOPPING CSM APP")
		return
	}
}
