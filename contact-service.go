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
	var (
		contactReq ContactRequest
		contactDet []ContactDetail
		contact    Contact
		conCheck   []Contact
	)

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
	var (
		contactReq    ContactRequest
		contactDetReq []ContactDetail
		contact       Contact
	)

	if request.Body == nil {
		http.Error(response, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(request.Body).Decode(&contactReq)
	if err != nil {
		http.Error(response, err.Error(), 400)
		return
	}

	contact.ID = contactReq.ID
	contactDetReq = contactReq.ContactDetail

	// Exist check
	db.First(&contact).Where("id = ?", contactReq.ID)

	if contact.ID > 0 {
		if contactReq.FirstName != "" {
			contact.FirstName = contactReq.FirstName
		}
		if contactReq.LastName != "" {
			contact.LastName = contactReq.LastName
		}

		db.Save(&contact)

		if len(contactDetReq) != 0 {
			for k := range contactDetReq {
				contactDet := ContactDetail{}

				db.First(&contactDet, "id = ?", contactDetReq[k].ID)

				if contactDetReq[k].ContactInfo != "" {
					contactDet.ContactInfo = contactDetReq[k].ContactInfo
				}
				if contactDetReq[k].Active != contactDet.Active {
					contactDet.Active = contactDetReq[k].Active
				}

				db.Save(&contactDet)
			}
		}

		response.WriteHeader(200)
		response.Write([]byte("Contact Updated"))
	} else {
		response.WriteHeader(400)
		response.Write([]byte("User not found"))
	}

}

// DeleteContactEndpoint get a Contact
func DeleteContactEndpoint(response http.ResponseWriter, request *http.Request) {
	var contact Contact
	params := mux.Vars(request)

	// Get Contact to delete
	db.Where("id = ?", params["id"]).First(&contact)

	// Make sure contact exists
	if contact.ID > 0 {
		db.Delete(&contact)
		response.WriteHeader(200)
		response.Write([]byte("Deleted Contact"))
	} else {
		response.WriteHeader(400)
		response.Write([]byte("Contact not found"))
	}

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
