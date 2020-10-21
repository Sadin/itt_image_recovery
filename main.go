package main

import (
	"flag"
	"fmt"
	"io"
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

	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	// create log file
	t := time.Now()
	fileName := fmt.Sprintf("ir_logs_%s.txt", t.Format("2006-01-02_15-04-05"))
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal(err)
	}

	// setup multiwriter for file & console logging
	w := io.MultiWriter(file, os.Stdout)

	// init loggers
	InfoLogger = log.New(w, "INFO:", log.LstdFlags|log.Lshortfile)
	WarnLogger = log.New(w, "WARN:", log.LstdFlags|log.Lshortfile)
	ErrorLogger = log.New(w, "ERRO:", log.LstdFlags|log.Lshortfile)
}

func main() {

	// time tracking
	defer execTime(time.Now())

	path, err := os.Getwd()
	// check error
	wordPtr := flag.String("path", path, "desired file path where program will execute")

	InfoLogger.Println("Starting image recovery...")

	flag.Parse()

	// fetch user and host info
	user, err := user.Current()
	if err != nil {
		WarnLogger.Println("Unable to find user information.")
	}
	host, err := os.Hostname()
	if err != nil {
		WarnLogger.Println("Unable to find hostname information.")
	}

	// print user/system info to console for awareness
	info := fmt.Sprintf("Running as: %s (id: %s)\tHostname: %s\tExecuting in: '%s'\t", user.Username, user.Uid, host, *wordPtr)
	InfoLogger.Println(info)

	// obtain slice of image directory
	subdir, err := ioutil.ReadDir(*wordPtr)
	if err != nil {
		ErrorLogger.Fatalln("Unable to obtain slice of image directory, exiting.")
	}
	// drill down into image directory
	err = os.Chdir(*wordPtr)
	if err != nil {
		ErrorLogger.Fatalln("Unable to change working dir to image directory, exiting.")
	}
	// loop over slice of patient image directory and do work
	for _, entry := range subdir {
		needsaction := false

		// build path string
		path := fmt.Sprintf("%s/%s", entry.Name(), RecoveryDir)

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			needsaction = true
		}

		output := fmt.Sprintf("Found %s | contains recovery images: %v", entry.Name(), needsaction)
		InfoLogger.Println(output)

		if needsaction {
			// enter patient directory
			InfoLogger.Println("Switching working directory...")
			err = os.Chdir(entry.Name())
			if err != nil {
				ErrorLogger.Fatalln("Unable to change working directory, exiting.")
			}

			// scan OriginalImages dir for slice of file objects
			imgs, err := ioutil.ReadDir(RecoveryDir)
			if err != nil {
				ErrorLogger.Fatalln("Unable to obtain slice of recovery dir, exiting.")
			}

			if len(imgs) > 0 {
				err = os.Chdir(RecoveryDir)
				if err != nil {
					ErrorLogger.Fatalln("Unable to change working directory, exiting.")
				}

				// loop through OriginalImages slice
				for _, entry := range imgs {
					output = fmt.Sprintf("\t%s found...", entry.Name())
					InfoLogger.Println(output)
					// rename
					if strings.Contains(entry.Name(), "Original") {
						fileRename(entry.Name())
					}

				}

				err = os.Chdir("..")
				if err != nil {
					ErrorLogger.Fatalln("Unable to change working directory, exiting.")
				}

				// recheck after files are moved
				imgs, err = ioutil.ReadDir(RecoveryDir)
				if err != nil {
					ErrorLogger.Fatalln("Unable to obtain slice of recovery dir, exiting.")
				}

				if len(imgs) == 0 {
					if removeDir(RecoveryDir) != nil {
						WarnLogger.Println("Unable to remove recovery directory.")
					}
				}

			} else {
				removeDir(RecoveryDir)
			}

			// return
			err = os.Chdir("..")
			if err != nil {
				ErrorLogger.Fatalln("Unable to return to patient directory, exiting.")
			}
		}
	}
}

func fileRename(name string) (string, error) {
	// fix filename, and perform rename + move
	newName := fmt.Sprintf("..\\%s", strings.ReplaceAll(name, " Original", ""))

	InfoLogger.Printf("\t\tRenaming & moving %s --> %s\n", name, newName)
	err := os.Rename(name, newName)
	if err != nil {
		ErrorLogger.Fatalln("Unable to rename image file, exiting.")
	}
	InfoLogger.Println("\t\tSuccess")

	return newName, fmt.Errorf("Error: renaming %s failed", name)
}

func removeDir(name string) error {
	InfoLogger.Printf("\t%s dir empty, removing...\n", name)
	err := os.Remove(name)
	return err
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

	InfoLogger.Printf("%s took %s", name, elapsed)

}
