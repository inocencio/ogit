/**
Author: Inocencio
Description: Check if the current source is a github project and open its project site.
**/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getURL(filename string) string {
	file, err := os.Open(filename)
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		tokens := strings.Split(line, "=")

		if len(tokens) > 1 {
			if strings.TrimSpace(strings.ToLower(tokens[0])) == "url" {
				return strings.TrimSpace(tokens[1])
			}
		}

	}

	return ""
}

func openBrowser(filepath string) {
	var err error
	var url string

	//acomplish the full path with configure file.
	filepath += "/config"

	_, err = os.Stat(filepath)

	if os.IsNotExist(err) {
		log.Fatal("It was unable to retrieve its GitHub site from this path.")
	}

	//try to get the URL
	url = getURL(filepath)

	if len(url) == 0 {
		log.Fatal("It was unable to parse URL.")
	}

	fmt.Printf("Opening %s...\n", url)

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Unsupported platform\n")
	}

	checkErr(err)
}

func main() {
	cDir, err := os.Getwd()
	checkErr(err)

	//check if path is a valid hidden github folder
	var path = cDir + "/.git"
	fi, err := os.Stat(cDir)
	checkErr(err)

	if fi.IsDir() {
		openBrowser(path)
	}
}
