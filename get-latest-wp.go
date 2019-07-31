package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/buger/jsonparser"
)

var countryCodes = map[string]string{
	"en": "",
	"es": "es_ES",
	"de": "de_DE",
	"fr": "fr_FR",
	"it": "it_IT",
}

func main() {

	var countryCode string
	flag.StringVar(&countryCode, "country", "en", "The two-letter country code for the version you want to download (en-english, es-spanish, de-german, fr-french, ...)")

	countryCode = strings.ToLower(countryCode)

	var langCodef string
	flag.StringVar(&langCodef, "lang", "", "The language code in the format xx_YY")

	// explicit version
	var version string
	flag.StringVar(&version, "version", "", "The version of Wordpress you want to download")

	flag.Parse()

	langCode, ok := countryCodes[countryCode]
	if !ok {
		if langCodef == "" {
			log.Fatalf("Language code for country \"%s\" was not found. You can use the -lang flag in combination with the -country flag to manually indicate the country and the lang like so:\n./get-latest-wp -country xx -lang xx_XX", countryCode)
		}

		langCode = langCodef
	}

	var getWpURL string
	if countryCode != "en" {
		getWpURL = fmt.Sprintf("https://%s.wordpress.org/latest-%s.tar.gz", countryCode, langCode)
	} else {
		// en
		getWpURL = "https://wordpress.org/latest.tar.gz"
	}

	if version == "" {
		log.Printf("Checking the last version of Wordpress...\n")
		var err error
		if version, err = checkVersion(); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Downloading Wordpress %s...\n", version)
	resWp, err := http.Get(getWpURL)
	if err != nil {
		log.Fatal(err)
	}

	var fileName string
	if countryCode != "en" {
		fileName = fmt.Sprintf("wordpress-%s-%s.tar.gz", version, langCode)
	} else {
		// en
		fileName = fmt.Sprintf("wordpress-%s.tar.gz", version)
	}

	fileWp, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fileWp.Close()

	// read from the reader (resWp.Body) and copy to the writer (fileWp)
	_, err = io.Copy(fileWp, resWp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\"%s\" file was successfully downloaded!\n", fileName)
}

// checkVersion request the latest version using an API endpoint provided by Wordpress
func checkVersion() (string, error) {

	const getVersionURL = "https://api.wordpress.org/core/version-check/1.7/"

	resVersion, err := http.Get(getVersionURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resVersion.Body.Close()

	dataVersion, err := ioutil.ReadAll(resVersion.Body)
	if err != nil {
		return "", err
	}

	version, err := jsonparser.GetString(dataVersion, "offers", "[0]", "current")
	if err != nil {
		return "", err
	}

	return version, nil
}
