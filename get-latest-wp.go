package main

import (
	"flag"
	"fmt"
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

	flag.Parse()

	langCode, ok := countryCodes[countryCode]
	if !ok {
		if langCodef == "" {
			log.Fatalf("Language code for country \"%s\" was not found. You can use the -lang flag in combination with the -country flag to manually indicate the country and the lang like so:\n./get-latest-wp -country xx -lang xx_XX", countryCode)
		}

		langCode = langCodef
	}

	var urlGetWp string
	if countryCode != "en" {
		urlGetWp = fmt.Sprintf("https://%s.wordpress.org/latest-%s.tar.gz", countryCode, langCode)
	} else {
		// en
		urlGetWp = "https://wordpress.org/latest.tar.gz"
	}

	log.Printf("Checking the last version of Wordpress...\n")
	version, err := checkVersion()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Downloading Wordpress %s...\n", version)
	resWp, err := http.Get(urlGetWp)
	if err != nil {
		log.Fatal(err)
	}

	dataWp, err := ioutil.ReadAll(resWp.Body)
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

	fileWp, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fileWp.Close()

	fileWp.Write(dataWp)
	log.Printf("\"%s\" file was successfully downloaded!\n", fileName)
}

// checkVersion request the latest version using an API endpoint provided by Wordpress
func checkVersion() (string, error) {
	urlGetVersion := "https://api.wordpress.org/core/version-check/1.7/"

	resVersion, err := http.Get(urlGetVersion)
	if err != nil {
		log.Fatal(err)
	}

	dataVersion, err := ioutil.ReadAll(resVersion.Body)
	if err != nil {
		return "", err
	}
	defer resVersion.Body.Close()

	version, _, _, err := jsonparser.Get(dataVersion, "offers", "[0]", "current")
	if err != nil {
		return "", err
	}

	return string(version), nil
}
