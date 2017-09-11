package main

import (
	"flag"
	"fmt"

	"github.com/riomhaire/lightauth/services"
)

func main() {
	// Command lines
	username := flag.String("user", "anonymous", "Username associated with the token")
	password := flag.String("password", "", "Password to use - cannot be empty")
	roles := flag.String("roles", "guest:public", "List of roles separated by ':'")

	flag.Parse()
	if len(*password) == 0 {
		fmt.Println("Password Parameter cannot be empty")
		return
	}
	hash := services.HashPassword(*password, fmt.Sprintf("%v%v", *username, *password))
	fmt.Printf("%v,%v,true,%v\n", *username, hash, *roles)
}
