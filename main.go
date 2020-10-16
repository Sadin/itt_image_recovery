package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	// RecoveryDir directory containing images for recovery
	RecoveryDir = "OriginalImages.XVA"

	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func main() {

	path, err := os.Getwd()
	check(err)
	wordPtr := flag.String("path", path, "desired file path where program will execute")

	fmt.Println("Starting image recovery...")

	flag.Parse()

	// fetch user and host info
	user, err := user.Current()
	check(err)
	host, err := os.Hostname()
	check(err)

	// print user/system info to console for awareness
	info := fmt.Sprintf("Running as: \t%s (id: %s)\nHostname: \t%s\nExecuting in: \t'%s'\n", user.Username, user.Uid, host, *wordPtr)
	fmt.Println(info)

	// time tracking
	defer execTime(time.Now())

	// obtain slice of image directory
	subdir, err := ioutil.ReadDir(*wordPtr)
	check(err)
	// drill down into image directory
	err = os.Chdir(*wordPtr)
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
	newName := fmt.Sprintf("..\\%s", strings.ReplaceAll(name, " Original", ""))

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

func execTime(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.Println(fmt.Sprintf("%s took %s", name, elapsed))
}
