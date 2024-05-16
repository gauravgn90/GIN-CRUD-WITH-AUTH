package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bxcodec/faker/v3"
)

type User struct {
	Name     string `faker:"name"`
	Username string `faker:"username"`
	Email    string `faker:"email"`
	Password string `faker:"password"`
}

func main() {
	GenerateFakeData()
}

func GenerateFakeData() {
	// Create a new file to write
	file, err := os.Create("users.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Create a JSON encoder for the file
	encoder := json.NewEncoder(file)

	// Generate and write 100,000 records
	for i := 0; i < 100; i++ {
		// Generate fake data for the User struct
		var user User
		err := faker.FakeData(&user)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Write the generated record to the file
		err = encoder.Encode(user)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	// Print a success message
	fmt.Println("Generated 100,000 records and written to users.json successfully.")
}
