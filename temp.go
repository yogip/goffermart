package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := []byte("MyDsdfdsfdsfd1234sfsadfsdfdsafasfdsafdsafdsafarkSecret")
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	fmt.Println(len(string(hashedPassword)), string(hashedPassword))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))
	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	fmt.Println(err) // nil means it is a match
}
