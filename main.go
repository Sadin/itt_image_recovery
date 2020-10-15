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
		// drill down into patient directory
		err = os.Chdir(entry.Name())
		check(err)
		// no recovery images unless proven otherwise
		y := false
		if _, err := os.Stat("OriginalImages.XVA"); os.IsNotExist(err) {
			y = true
		}
		output := fmt.Sprintf("Found %s | contains recovery images: %v", entry.Name(), y)
		fmt.Println(output)

		// return up to
		err = os.Chdir("..")
		check(err)
	}

	oldName := "test.txt"
	newName := "testing.txt"
	err = os.Rename(oldName, newName)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
