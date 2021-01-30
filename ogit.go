/**
Author: Inocêncio T. de Oliveira

Check if the current source code is a github project and then open its github location site.

Just type ogit on terminal and correlative github site will be opened.
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

/**
Get github URL from /.git/config file.
 */
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

	filepath += "/config"

	if _, err = os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Error: this is not a valid GitHub Project!")
		}
	}

	url = getURL(filepath)

	if len(url) == 0 {
		panic("It was unable to parse URL.")
		return
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
		err = fmt.Errorf("Unsupported platform")
	}

	checkErr(err)
}

func main() {
	cDir, err := os.Getwd()
	checkErr(err)

	//check if path is a valid hidden github folder
	var path = cDir + "/.git"
	fi, err := os.Stat(cDir)

	if os.IsNotExist(err) {
		log.Fatal("Error: this is not a valid GitHub Project!")
	}

	if fi.IsDir() {
		openBrowser(path)
	}
}
