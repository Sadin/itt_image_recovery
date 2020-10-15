package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

func main() {
	fmt.Println("Starting image recovery...")

	// fetch user and host info
	user, err := user.Current()
	check(err)
	host, err := os.Hostname()
	check(err)

	oldName := "test.txt"
	newName := "testing.txt"
	err = os.Rename(oldName, newName)
	check(err)

	// print user/system info to console for awareness
	info := fmt.Sprintf("Running as: %s (id: %s)\nHostname: %s \n\n", user.Username, user.Uid, host)
	fmt.Println(info)

	// obtain slice of image directory
	subdir, err := ioutil.ReadDir("imagefiles")
	check(err)
	// drill down into image directory
	err = os.Chdir("imagefiles")
	check(err)
	fmt.Println("Listing subdirectories... ")
	// loop over slice of patient image directory and do work
	for _, entry := range subdir {
		needsaction := false

		// build path string
		path := fmt.Sprintf("%s/OriginalImages.XVA", entry.Name())

		if _, err := os.Stat(path); os.IsNotExist(err) {
			needsaction = true
		}

		output := fmt.Sprintf("Found %s | contains recovery images: %v", entry.Name(), needsaction)
		fmt.Println(output)

		if needsaction == true {
			// enter patient directory
			fmt.Println("Switching working directory...")
			err = os.Chdir(entry.Name())
			check(err)

			// return
			err = os.Chdir("..")
			check(err)
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
