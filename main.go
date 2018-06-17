package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

func main() {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server, user, password, port, database)
	db, err = gorm.Open("mssql", connectionString)

	if err != nil {
		log.Fatal("Failed to create connection pool. Error: " + err.Error())
	}
	gorm.DefaultCallback.Create().Remove("mssql:set_identity_insert")
	defer db.Close()

	db.AutoMigrate(&Contact{}, &ContactType{}, &ContactDetail{})
	// db.Create(&ContactType{TypeName: "Email", TypeDescription: "Email Address"})
	// db.Create(&ContactType{TypeName: "Phone", TypeDescription: "Phone Number"})

	router := mux.NewRouter()
	router.HandleFunc("/", RootEndpoint).Methods("GET")
	router.HandleFunc("/contact", GetContactsEndpoint).Methods("GET")
	router.HandleFunc("/contact/{id}", GetContactEndpoint).Methods("GET")
	router.HandleFunc("/contact/create/", CreateContactEndpoint).Methods("POST")
	router.HandleFunc("/contact/update/{id}", UpdateContactEndpoint).Methods("PUT")
	router.HandleFunc("/contact/delete/{id}", DeleteContactEndpoint).Methods("DELETE")
	router.HandleFunc("/contact/upload/", UploadContactsEndpoint).Methods("POST")
	router.HandleFunc("/contact/download/", DownloadContactsEndpoint).Methods("POST")
	fmt.Println("Listening at :12344")
	log.Fatal(http.ListenAndServe(":12344", router))

}
