/**
Author: Inocencio
Description: Check if the current source is a github project and open its project site.
Works on: linux, macOS and Windows.
**/

package main

import (
	"bufio"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/tucnak/store"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

var configFileName = ".ogit.toml"

type Options struct {
	Browser string `short:"b" long:"browser" description:"Set a browser path to ogit's default browser.'"`
}

type Config struct {
	Browser string `toml:"browser"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/**
Get a fragmented string by separator and its side.
 */
func getString(s string, sep string, index int) string {
	tokens := strings.Split(s, sep)

	//if split worked then return element from the index
	if len(tokens) > 1 {
		return strings.TrimSpace(tokens[index])
	}

	return ""
}

func getConfigPath() string {
	fp, err := os.UserConfigDir()
	checkErr(err)
	return path.Join(fp, configFileName)
}

//func getProperties(fp string) *properties.Properties {
//	return properties.MustLoadFile(fp, properties.UTF8)
//}

func getProperties() *Config {
	var config Config

	if err := store.Load(configFileName, &config); err != nil {
		fmt.Println("Err:06 - It was unable to save .ogit file")
		os.Exit(-1)
	}

	return &config
}

func createPropertiesFile() {
	if _, err := os.Stat(getConfigPath()); os.IsNotExist(err) {
		f, err := os.Create(getConfigPath())
		checkErr(err)
		defer f.Close()
	}
}

func saveProperties(config* Config) {
	//create a file to store our properties if it doesn't exist
	createPropertiesFile()

	//write properties content if it exists
	if config != nil {
		//load properties in saveConfig to change its content if needed it
		var saveConfig *Config
		saveConfig = getProperties()

		//check the properties differences and store it if they are different
		if config.Browser != saveConfig.Browser {
			saveConfig.Browser = config.Browser
		}

		if err := store.Save(configFileName, &saveConfig); err != nil {
			fmt.Println("Err:06 - It was unable to save .ogit file")
			os.Exit(-1)
		}

		//var fc string
		//fc += "browser=" + config.Browser
		//
		////write to file
		//f, err := os.Open(fp)
		//checkErr(err)
		//_, err = f.WriteString(fc)
		//checkErr(err)
		//defer f.Close()
	}
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
func openBrowser(filepath string, config* Config) {
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

	if len(config.Browser) > 0 {
		err = exec.Command(config.Browser, url).Start()
		checkErr(err)
	} else {
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
}

func checkConfigFile() *Config {
	var options Options
	var config Config
	configChanged := false

	//parsing the program arguments...
	parser := flags.NewParser(&options, flags.Default)
	_, err := parser.Parse()
	checkErr(err)

	//check if a browser argument is set
	if len(options.Browser) > 0 {
		config.Browser = options.Browser
		configChanged = true
	}

	//save properties to properties file if config is changed
	if configChanged {
		saveProperties(&config)
	}

	//always load the properties file
	conf := getProperties()

	return conf
}

func init() {
	//create a properties file if it doesn't exist
	saveProperties(nil)

	//register user's config path
	configPath, err := os.UserConfigDir()
	checkErr(err)

	store.Init(configPath)
}

func main() {
	//get current home
	pwd, err := os.Getwd()
	checkErr(err)

	//mount a properties file path
	var gitPath = pwd + "/.git"
	fi, err := os.Stat(gitPath)

	//check if gitPath is a valid hidden github folder
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Err:01 - There is no '.git' folder from working dir.")
			os.Exit(-1)
		}
	}

	//check if gitPath is a directory and then try to open a browser
	if fi.IsDir() {
		var config = checkConfigFile()
		openBrowser(gitPath, config)
	} else {
		fmt.Println("Err:02 The '.git' is not a valid GitHub folder, it's a file instead.")
		os.Exit(-1)
	}
}
