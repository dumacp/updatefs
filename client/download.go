package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func DownloadFile(client *http.Client, url, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// // Get the data
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return err
	// }

	// transport := loadLocalCert()
	// client := http.Client{Transport: transport}

	// resp, err := client.Do(req)
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("error: GET, %s", err)
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("error: COPY, %s", err)
		return err
	}

	return nil
}
