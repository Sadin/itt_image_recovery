package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
)

func main() {
	fmt.Println("Starting image recovery...")

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	// print user/system info to console for awareness
	info := fmt.Sprintf("Running as: %s (id: %s)\nHostname: %s \n\n", user.Username, user.Uid, host)
	fmt.Println(info)

	oldName := "test.txt"
	newName := "testing.txt"
	err = os.Rename(oldName, newName)
	if err != nil {
		log.Fatal(err)
	}
}
