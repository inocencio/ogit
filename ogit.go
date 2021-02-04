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

func getString(s string, sep string, index int) string {
	tokens := strings.Split(s, sep)

	//if split worked then return element from the index
	if len(tokens) > 1 {
		return strings.TrimSpace(tokens[index])
	}

	return ""
}

/**
Try to get an URL from .git/config file from url parameter.
 */
func getURL(filename string) string {
	file, err := os.Open(filename)
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		tokens := strings.Split(line, "=")

		//check if there is an URL and further if it's a site or e-mail
		if len(tokens) > 1 {
			//is there a URL?
			if strings.TrimSpace(strings.ToLower(tokens[0])) == "url" {
				addr := strings.TrimSpace(tokens[1])

				//check if address starts with https then return it
				if strings.HasPrefix(addr, "https") {
					return addr
				} else {
					//check if addr contains 'git@'
					addr = getString(addr, "git@", 1)

					//is an e-mail?
					if len(addr) > 0 {
						return "https://" + strings.Replace(addr, ":", "/", 1)
					}
				}
			} //url
		} //=
	} //scan

	return ""
}

/**
Open a browser based on current OS and defined browser as default.
 */
func openBrowser(filepath string) {
	var err error
	var url string

	//accomplish the full path with configure file.
	filepath += "/config"

	if _, err = os.Stat(filepath); os.IsNotExist(err) {
		fmt.Println("Err:03 - It was unable to retrieve its GitHub site because there isn't a config file inside the '.git' folder")
		os.Exit(-1)
	}

	//try to get the URL
	url = getURL(filepath)

	if len(url) == 0 {
		fmt.Println("Err:04 - It was unable to parse URL.")
		os.Exit(-1)
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
	pwd, err := os.Getwd()
	checkErr(err)

	var path = pwd + "/.git"
	fi, err := os.Stat(path)

	//check if path is a valid hidden github folder
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Err:01 - There is no '.git' folder from working dir.")
			os.Exit(-1)
		}
	}

	//check if path is a directory and then try to open a browser
	if fi.IsDir() {
		openBrowser(path)
	} else {
		fmt.Println("Err:02 The '.git' is not a valid GitHub folder, it's a file instead.")
		os.Exit(-1)
	}
}
