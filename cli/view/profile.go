package View

import (
	"context"
	rd "crypto/rand"
	"fmt"
	"log"

	Db "github.com/PARTHIB-DEB/TermoTalks/db"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database = Db.MongoConnect()

func CreateProfileView(fname string, lname string, email string, pwd string, pwd2 string) {
	// This function is intended to create a profile view.
	// The implementation details will be added later.

	// Validation
	if fname == "" {
		log.Fatal("First name is required to create a profile")
	}
	if lname == "" {
		log.Fatal("Last name is required to create a profile")
	}
	if email == "" {
		log.Fatal("Email is required to create a profile")
	}
	if pwd == "" {
		log.Fatal("Password is required to create a profile")
	}
	if pwd2 == "" {
		log.Fatal("Confirm Password is required to create a profile")
	}
	if pwd != pwd2 {
		log.Fatal("Passwords do not match")
	}

	// Username Generation
	uname := ""
	uname += rd.Text()[:10]
	for i := 0; i <= len(fname); i = i + 2 {
		uname += fname[i : i*2]
		uname += lname[i : i*2]
	}
	uname += rd.Text()[:10]

	// Creation
	creds := map[string]string{
		"username":  uname,
		"firstname": fname,
		"lastname":  lname,
		"email":     email,
		"password":  pwd,
	}
	_, err := db.Collection("profiles").InsertOne(context.Background(), creds)
	if err != nil {
		log.Fatal("Error creating profile:", err)
	}
}

func UpdateProfileView(fname string, lname string, email string, pwd string, uname string) {
	if uname == "" {
		log.Fatal("Username is required to update a profile")
	} else {

		// Fetch That Tuple
		type profile struct {
			firstname string
			lastname  string
			email     string
			pwd       string
		}
		var fetched_prof, new_prof profile
		tuple := db.Collection("profiles").FindOne(context.Background(), map[string]string{"username": uname})
		if tuple.Err() != nil {
			log.Fatal("Profile not found for username:", uname)
		}
		err := tuple.Decode(&fetched_prof)
		if err != nil {
			log.Fatal(err.Error())
		}

		// logic to validate fields
		if fname != fetched_prof.firstname && fname != "" {
			new_prof.firstname = fname
		}

		if lname != fetched_prof.lastname && lname != "" {
			new_prof.lastname = lname
		}

		if email != fetched_prof.email && email != "" {
			new_prof.email = email
		}

		if pwd != fetched_prof.pwd && pwd != "" {
			var pwd2 string
			fmt.Print("Enter Password Again :")
			fmt.Scanf("%s", &pwd2)
			if pwd2 != pwd {
				log.Fatal("Passwords Should match")
			}
			new_prof.pwd = pwd
		}

		_, err = db.Collection("profiles").UpdateOne(context.Background(), &fetched_prof, &new_prof)
		if err != nil {
			log.Fatal("Error updating profile:", err)
		}
	}
}

func DeleteProfileView(uname string) {
	if uname == "" {
		log.Fatal("Username is required to Delete a profile")
	} else {
		type unameS struct {
			username string
		}
		_, err := db.Collection("profiles").DeleteOne(context.Background(), &unameS{username: uname})
		if err != nil {
			log.Fatal(err)
		}
	}
}
