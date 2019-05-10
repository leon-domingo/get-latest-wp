package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/buger/jsonparser"
)

func main() {
	urlGetVersion := "https://api.wordpress.org/core/version-check/1.7/"
	urlGetWp := "https://es.wordpress.org/latest-es_ES.tar.gz"

	log.Printf("Checking the last version of Wordpress...\n")
	resVersion, err := http.Get(urlGetVersion)
	if err != nil {
		log.Fatal(err)
	}
	defer resVersion.Body.Close()

	dataVersion, err := ioutil.ReadAll(resVersion.Body)
	if err != nil {
		log.Fatal(err)
	}

	version, _, _, err := jsonparser.Get(dataVersion, "offers", "[0]", "current")
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

	fileName := fmt.Sprintf("wordpres-%s-es_ES.tar.gz", version)

	fileWp, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fileWp.Close()

	fileWp.Write(dataWp)
	log.Printf("\"%s\" file was successfully downloaded!\n", fileName)
}
