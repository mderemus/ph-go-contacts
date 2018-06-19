# ph-go-contacts

This application was intended to fulfill a need to create, read, update, and delete contacts for the following fields:

- First Name
- Last Name
- Email Address
- Phone Number

The design was modeled after how maybe the Contacts application on your mobile phone would work. I wanted to provide flexibility in that a Contact could have multiple means of communication, i.e. several emails and perhaps more than one phone number. With that that in mind, I designed the database with the details laid out in the next section.

## Database

MS SQL Server was used as the database with the following table structures.

- contact_details: Holds Contact ID and their corresponding Contact Types.
- contact_types: The type of Contacts a person can have, e.g. Phone & Email
- contacts: The main entry for Contacts that stores their First and Last Name

Example Data:
Bob Smith, Phone: 555-999-1098, Phone: 489-222-9933
Chris Sheppard, Phone: 533-939-3398, Email: cshep@go .com

***Note:** Run the database.sql file in MS SQL Server Studio to create the database, tables, and the sql server login necessary for the application. You will then need to update the `var server` variable to the name of your server*

## Testing

I used Postman for testing as well as manually hitting the API address in Chrome and monitoring the Network tab of the Developer Tools. Below are some examples of the API calls you can make and the necessary data structures of the JSON request.

## Create Contact

This creates a Contact and inserts the appropriate Contact Details if present. It also checks for duplicates and responds appropriately.

`http://localhost:12344/contact/create/`

    {
      "firstname": "Bob",
      "lastname": "Smith",
      "contactdetail": [
        {
          "contacttypeid": 1,
          "contactinfo": "lalal@go.com"
        },
        {
            "contacttypeid": 2,
            "contactinfo": "333-554-7834"
        }
      ]
    }

If Successful, it will return a message to the user "Contact Created" as well as a Status 200 OK. If a duplicate is matched, the user will receive a Status 400 as well as a message "Duplicate Contact".

## Read all Contacts

The following endpoint will be used to retrieve all of the contacts

    http://localhost:12344/contact

Which, if successful, the call will return JSON to the user in the following format as well as a Status 200 OK

    {
            "ID": 13,
            "CreatedAt": "2018-06-16T17:22:44.9974771-05:00",
            "UpdatedAt": "2018-06-16T17:22:44.9974771-05:00",
            "DeletedAt": null,
            "FirstName": "Gary",
            "LastName": "Busey",
            "ContactDetails": [
                {
                    "ID": 4,
                    "CreatedAt": "2018-06-16T17:22:45.0094474-05:00",
                    "UpdatedAt": "2018-06-16T17:22:45.0094474-05:00",
                    "DeletedAt": null,
                    "ContactID": 13,
                    "ContactTypeID": 1,
                    "ContactInfo": "wgarye@go.com",
                    "Active": true
                },
                {
                    "ID": 5,
                    "CreatedAt": "2018-06-16T17:22:45.017422-05:00",
                    "UpdatedAt": "2018-06-16T17:22:45.017422-05:00",
                    "DeletedAt": null,
                    "ContactID": 13,
                    "ContactTypeID": 2,
                    "ContactInfo": "333-554-7834",
                    "Active": true
                }
            ]
        }

## Get a specific Contact

The following endpoint will be used to retrieve a specific Contact

    http://localhost:12344/contact/{id}

## Update a Contact

The following endpoint will be used to update a single Contact

    http://localhost:12344/contact/update/

An example of the request data should look like this:

    {
            "ID": 13,
            "FirstName": "Gary",
            "contactdetail": [
                {
                    "ID": 4,
                    "ContactID": 13,
                    "ContactTypeID": 1,
                    "ContactInfo": "garybusey@gmails.com",
                    "Active": true
                },
                {
                    "ID": 5,
                    "ContactID": 13,
                    "ContactTypeID": 2,
                    "ContactInfo": "955-469-1234",
                    "Active": true
                }
            ]
        }

If Successful, it will return a message to the user "Contact Updated" as well as a Status 200 OK

## Delete A Contact

The following endpoint will be used to delete a specific Contact with the matching {id}

    http://localhost:12344/contact/delete/{id}

## To Do

Not necessarily in this order

- Test Scripts
- Complete Upload of CSV
- Complete Download of CSV
