package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

var (
	// RecoveryDir directory containing images for recovery
	RecoveryDir = "OriginalImages.XVA"
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
		needsaction := false

		// build path string
		path := fmt.Sprintf("%s/%s", entry.Name(), RecoveryDir)

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			needsaction = true
		}

		output := fmt.Sprintf("Found %s | contains recovery images: %v", entry.Name(), needsaction)
		fmt.Println(output)

		if needsaction == true {
			// enter patient directory
			fmt.Println("Switching working directory...")
			err = os.Chdir(entry.Name())
			check(err)

			// scan OriginalImages dir for slice of file objects
			imgs, err := ioutil.ReadDir(RecoveryDir)
			check(err)

			if len(imgs) > 0 {
				err = os.Chdir(RecoveryDir)
				check(err)

				// loop through OriginalImages slice
				for _, entry := range imgs {
					fmt.Printf("\t%s found...\n", entry.Name())

					// rename
					if strings.Contains(entry.Name(), "Original") {
						fileRename(entry.Name())
					}

				}

				err = os.Chdir("..")
				check(err)

				// recheck after files are moved
				imgs, err = ioutil.ReadDir(RecoveryDir)
				check(err)
				if len(imgs) == 0 {
					removeDir(RecoveryDir)
				}

			} else {
				removeDir(RecoveryDir)
			}

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

func fileRename(name string) (string, error) {
	// fix filename, and perform rename + move
	newName := fmt.Sprintf("..\\%s", strings.Replace(name, " Original", "", -1))

	fmt.Printf("\t\tRenaming & moving %s --> %s\n", name, newName)
	err := os.Rename(name, newName)
	check(err)
	fmt.Println("\t\tSuccess")

	return newName, fmt.Errorf("Error: renaming %s failed", name)
}

func removeDir(name string) (string, error) {
	fmt.Printf("\t%s dir empty, removing...\n", name)
	err := os.Remove(name)
	return "Success", err
}
