package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"hw2/controllers/stdhttp" 
	"hw2/gate/psg"
)

func main() {
	
	dbHost := flag.String("dbhost", "localhost:5432", "Database host address")
	dbUser := flag.String("dbuser", "postgres", "Database username")
	dbPassword := flag.String("dbpassword", "12345", "Database password")
	serverAddr := flag.String("serveraddr", "localhost:8080", "Server address")

	flag.Parse()

	db, err := psg.NewPsg(*dbHost, *dbUser, *dbPassword)
	if err != nil {
		log.Fatal("Error creating database connection:", err)
	}

	controller := stdhttp.NewController(*serverAddr, db)

	http.HandleFunc("/record/add", controller.RecordAdd)
	http.HandleFunc("/records/get", controller.RecordsGet)
	http.HandleFunc("/record/update", controller.RecordUpdate)
	http.HandleFunc("/record/delete", controller.RecordDeleteByPhone)

	serverAddress := fmt.Sprintf(":%s", *serverAddr)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}
