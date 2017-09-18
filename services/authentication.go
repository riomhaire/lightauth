package services

import (
	"crypto/sha256"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

type User struct {
	User     string   `json:"username"`
	Password string   `json:"password"` // Salted
	Enabled  bool     `json:"enabled"`
	Roles    []string `json:"roles"`
}

type AuthenticationArgs struct {
	User     string `json:"username"`
	Password string `json:"password"`
}

type AuthenticationReply struct {
	Token string `json:"token"`
}

type AuthenticationService int

func (t *AuthenticationService) Authenticate(r *http.Request, args *AuthenticationArgs, result *AuthenticationReply) error {
	log.Printf("Authentication Authenticate Called for '%v'\n", args.User)

	if val, ok := users[args.User]; ok {
		//do something here TODO should look up user and get SALT etc
		hash := HashPassword(args.Password, fmt.Sprintf("%v%v", args.User, args.Password)) // I know each user should have on salt which is not the user
		//log.Printf("'%s' == '%s'\n", hash, val.Password)
		if hash != val.Password {
			return errors.New("401 Not Allowed")
		}
		if val.Enabled == false {
			return errors.New("403 User '" + args.User + "' Disabled")
		}

		sessionKey, err := ApplicationSessionService.NewSession(val)
		if err != nil {
			return errors.New("500 " + err.Error())
		}
		*result = AuthenticationReply{sessionKey}
	} else {
		return errors.New("401 User Not Allowed")
	}
	return nil
}

// Initiaizes data structues - IE Read user DB
func LoadUsers(filename string) error {
	csvfile, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer csvfile.Close()
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("Reading User Database %s\n", filename)
	// Create user map
	for index, row := range records {
		if index > 0 && len(row) > 0 {
			user := User{}
			user.User = row[0]
			user.Password = row[1]

			v, _ := strconv.ParseBool(row[2])
			user.Enabled = v

			roles := strings.Split(row[3], ":")
			user.Roles = roles
			// Add
			users[user.User] = user
			log.Printf("\tUser %v, Enabled = %v\n", user.User, user.Enabled)
		}
	}
	log.Printf("#Number of users = %v\n", len(users))
	return nil
}

func HashPassword(password, salt string) string {
	bpassword := []byte(password)
	bsalt := []byte(salt)
	v := pbkdf2.Key(bpassword, bsalt, 4096, sha256.Size, sha256.New)
	hash := fmt.Sprintf("%x", v) // I know each user should have on salt
	//log.Printf("%s -> %s\n", password, hash)
	return hash
}

// Get version
func (t *AuthenticationService) Version(r *http.Request, args *string, result *string) error {
	*result = Version
	return nil
}
