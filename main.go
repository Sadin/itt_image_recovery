package main

import (
	"bytes"
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

	buf1, buf2 bytes.Buffer

	w = io.MultiWriter(&buf1, &buf2)
)

func main() {
	// begin log file
	t := time.Now()
	fileName := fmt.Sprintf("ir_logs_%s.txt", t.Format("2006-01-02_15-04-05"))
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	check(err, f)

	defer f.Close()

	path, err := os.Getwd()
	check(err, f)
	wordPtr := flag.String("path", path, "desired file path where program will execute")

	logOutput("Starting image recovery...", f)

	flag.Parse()

	// fetch user and host info
	user, err := user.Current()
	check(err, f)
	host, err := os.Hostname()
	check(err, f)

	// print user/system info to console for awareness
	info := fmt.Sprintf("Running as: \t%s (id: %s)\nHostname: \t%s\nExecuting in: \t'%s'\n", user.Username, user.Uid, host, *wordPtr)
	logOutput(info, f)

	// time tracking
	defer execTime(time.Now(), f)

	// obtain slice of image directory
	subdir, err := ioutil.ReadDir(*wordPtr)
	check(err, f)
	// drill down into image directory
	err = os.Chdir(*wordPtr)
	check(err, f)
	logOutput("Listing subdirectories... ", f)
	// loop over slice of patient image directory and do work
	for _, entry := range subdir {
		needsaction := false

		// build path string
		path := fmt.Sprintf("%s/%s", entry.Name(), RecoveryDir)

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			needsaction = true
		}

		output := fmt.Sprintf("Found %s | contains recovery images: %v", entry.Name(), needsaction)
		logOutput(output, f)

		if needsaction == true {
			// enter patient directory
			logOutput("Switching working directory...", f)
			err = os.Chdir(entry.Name())
			check(err, f)

			// scan OriginalImages dir for slice of file objects
			imgs, err := ioutil.ReadDir(RecoveryDir)
			check(err, f)

			if len(imgs) > 0 {
				err = os.Chdir(RecoveryDir)
				check(err, f)

				// loop through OriginalImages slice
				for _, entry := range imgs {
					output = fmt.Sprintf("\t%s found...", entry.Name())
					logOutput(output, f)
					// rename
					if strings.Contains(entry.Name(), "Original") {
						fileRename(entry.Name(), f)
					}

				}

				err = os.Chdir("..")
				check(err, f)

				// recheck after files are moved
				imgs, err = ioutil.ReadDir(RecoveryDir)
				check(err, f)
				if len(imgs) == 0 {
					removeDir(RecoveryDir)
				}

			} else {
				removeDir(RecoveryDir)
			}

			// return
			err = os.Chdir("..")
			check(err, f)
		}
	}
}

func check(err error, f *os.File) {
	if err != nil {
		logOutput(err.Error(), f)
	}
}

func fileRename(name string, f *os.File) (string, error) {
	// fix filename, and perform rename + move
	newName := fmt.Sprintf("..\\%s", strings.ReplaceAll(name, " Original", ""))

	fmt.Printf("\t\tRenaming & moving %s --> %s\n", name, newName)
	err := os.Rename(name, newName)
	check(err, f)
	fmt.Println("\t\tSuccess")

	return newName, fmt.Errorf("Error: renaming %s failed", name)
}

func removeDir(name string) (string, error) {
	fmt.Printf("\t%s dir empty, removing...\n", name)
	err := os.Remove(name)
	return "Success", err
}

func execTime(start time.Time, f *os.File) {
	elapsed := time.Since(start)
	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	logOutput(fmt.Sprintf("%s took %s", name, elapsed), f)

}

func logOutput(text string, logfile *os.File) (int, error) {
	r := strings.NewReader(text)

	if _, err := io.Copy(w, r); err != nil {
		log.Fatal(err)
	}

	o := fmt.Sprintf("%s\n", buf2.String())
	fmt.Println(buf1.String())
	return logfile.Write([]byte(o))
}
