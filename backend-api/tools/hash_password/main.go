// tools/hash_password/main.go
// Utility to generate bcrypt password hashes for seeding the admins table.
//
// Usage:
//   go run ./tools/hash_password -password admin123
//   go run ./tools/hash_password -password mySecret -cost 12

package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := flag.String("password", "", "Plain-text password to hash (required)")
	cost := flag.Int("cost", 10, "bcrypt cost factor (10–14 recommended)")
	flag.Parse()

	if *password == "" {
		log.Fatal("Usage: go run ./tools/hash_password -password <your_password>")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), *cost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	fmt.Printf("Password : %s\n", *password)
	fmt.Printf("Hash     : %s\n", string(hash))
	fmt.Printf("Cost     : %d\n", *cost)
}
