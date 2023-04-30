package main

import (
	"flag"
	"fmt"
	"io"
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

const (
	GET_VERSION_URL   = "https://api.wordpress.org/core/version-check/1.7/"
	BASE_DOWNLOAD_URL = "https://download.wordpress.org"
)

func main() {
	var countryCode string
	flag.StringVar(&countryCode, "country", "en", "The two-letter country code for the version you want to download (en-english, es-spanish, de-german, fr-french, ...)")
	countryCode = strings.ToLower(countryCode)

	var langCodeFlag string
	flag.StringVar(&langCodeFlag, "lang", "", "The language code in the format xx_YY")

	var version string
	flag.StringVar(&version, "version", "", "The version of Wordpress you want to download in the format X.Y.Z")

	flag.Parse()

	langCode, ok := countryCodes[countryCode]
	if !ok {
		if langCodeFlag == "" {
			log.Fatalf("Language code for country \"%s\" was not found. You can use the -lang flag in combination with the -country flag to manually indicate the country and the lang like so:\n./get-latest-wp -country xx -lang xx_XX", countryCode)
		}
		langCode = langCodeFlag
	}

	if version == "" {
		log.Println("Checking the latest version of Wordpress...")
		var err error
		if version, err = checkVersion(); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Downloading Wordpress %s...\n", version)
	resWp, err := http.Get(getDownloadWordpressURL(countryCode, langCode, version))
	if err != nil {
		log.Fatal(err)
	} else if resWp.StatusCode == http.StatusNotFound {
		log.Fatalf("Version %s does not exist", version)
	} else if resWp.StatusCode != http.StatusOK {
		log.Fatalf("An error ocurred while downloading Wordpress (%d)", resWp.StatusCode)
	}

	fileName := fmt.Sprintf("wordpress-%s.tar.gz", version)
	if countryCode != "en" {
		fileName = fmt.Sprintf("wordpress-%s-%s.tar.gz", version, langCode)
	}

	downloadedFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer downloadedFile.Close()

	_, err = io.Copy(downloadedFile, resWp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s file was successfully downloaded!\n", fileName)
}

func getDownloadWordpressURL(countryCode, langCode string, version string) string {
	wordpressURL := fmt.Sprintf("%s/wordpress-%s.tar.gz", BASE_DOWNLOAD_URL, version)
	if countryCode != "en" {
		wordpressURL = fmt.Sprintf("%s/wordpress-%s-%s.tar.gz", BASE_DOWNLOAD_URL, version, langCode)
	}

	return wordpressURL
}

// checkVersion requests the latest version using an API endpoint provided by Wordpress
func checkVersion() (string, error) {
	resVersion, err := http.Get(GET_VERSION_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resVersion.Body.Close()

	dataVersion, err := io.ReadAll(resVersion.Body)
	if err != nil {
		return "", err
	}

	version, err := jsonparser.GetString(dataVersion, "offers", "[0]", "current")
	if err != nil {
		return "", err
	}

	return version, nil
}
