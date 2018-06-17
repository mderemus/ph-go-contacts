package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var (
	server   = "DESKTOP-PO84PUU"
	port     = 1433
	user     = "PHContacts"
	password = "PHContacts"
	database = "PHContacts"
)

var db *gorm.DB
var err error

// Contact struct: Main contact holder
type Contact struct {
	gorm.Model
	FirstName      string
	LastName       string
	ContactDetails []ContactDetail `gorm:"foreignkey:ContactID"`
}

// ContactRequest struct: Like a view model
type ContactRequest struct {
	gorm.Model
	FirstName     string
	LastName      string
	ContactDetail []ContactDetail
}

// ContactType struct: The Type of Contact Info a Contact can have
type ContactType struct {
	gorm.Model
	TypeName        string // e.g. Cell, Home, Email
	TypeDescription string // e.g. Non-land line, Personal Email
}

// ContactDetail struct: To enable storing multiple types of contact info for a contact
type ContactDetail struct {
	gorm.Model
	ContactID     uint
	ContactTypeID uint   // select ContactInfo from ContactInfo where ContactTypeID = 1 and ContactID = 1
	ContactInfo   string // e.g. 555-555-5555 or bob@gmail.com. Tried to make this as generic as possible, it's just strings. I can validate on save
	Active        bool   // Numbers or email address may no longer be valid
}

// Contacts type
type Contacts []Contact

// ContactDetails type
type ContactDetails []ContactDetail

// RootEndpoint gets a root
func RootEndpoint(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte("Hello Root"))
}

// GetContactsEndpoint get Contacts
func GetContactsEndpoint(response http.ResponseWriter, request *http.Request) {
	var contacts []Contact
	fmt.Println("Getting all contacts...")
	// Get All Active Contacts with Corresponding Contact Details
	db.Where("deleted_at is NULL").Preload("ContactDetails").Find(&contacts)

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(&contacts)
}

// GetContactEndpoint get a Contact
func GetContactEndpoint(response http.ResponseWriter, request *http.Request) {
	var contact Contact
	params := mux.Vars(request)

	// Get Contact with Corresponding Contact Details
	db.Preload("ContactDetails", "contact_id = ?", params["id"]).Find(&contact)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)
	json.NewEncoder(response).Encode(&contact)
}

// CreateContactEndpoint create a Contact
func CreateContactEndpoint(response http.ResponseWriter, request *http.Request) {
	/*
		Todo:
		1.) Get id of inserted - Done
		2.) Pass id of inserted to contact detail - Done
		3.) Check for duplicates of contacts - Done
	*/
	var contactReq ContactRequest
	var contactDet []ContactDetail
	var contact Contact
	var conCheck []Contact

	if request.Body == nil {
		http.Error(response, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&contactReq)
	if err != nil {
		http.Error(response, err.Error(), 400)
		return
	}

	contact.FirstName = contactReq.FirstName
	contact.LastName = contactReq.LastName
	contactDet = contactReq.ContactDetail

	// Dupe check
	db.Where(Contact{FirstName: contact.FirstName, LastName: contact.LastName}).First(&conCheck)

	if len(conCheck) == 0 {
		contact.FirstName = contactReq.FirstName
		contact.LastName = contactReq.LastName

		db.Create(&Contact{FirstName: contact.FirstName, LastName: contact.LastName})
		db.Last(&contact)

		if len(contactDet) != 0 {
			for k := range contactDet {
				fmt.Println(contactDet[k].ContactInfo)
				db.Create(&ContactDetail{ContactID: contact.ID, ContactTypeID: contactDet[k].ContactTypeID, ContactInfo: contactDet[k].ContactInfo, Active: true})
			}
		}

		response.WriteHeader(200)
		response.Write([]byte("Contact Created"))
	} else {
		response.WriteHeader(400)
		response.Write([]byte("Duplicate contact!!"))
	}

}

// UpdateContactEndpoint get a Contact
func UpdateContactEndpoint(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte("Update a Contact"))
}

// DeleteContactEndpoint get a Contact
func DeleteContactEndpoint(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte("Delete a Contact"))
}

// UploadContactsEndpoint get a Contact
func UploadContactsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte("Upload Contacts"))
}

// DownloadContactsEndpoint get a Contact
func DownloadContactsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte("Download Contacts"))
}
