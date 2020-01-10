package main

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"net/http"
)

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func DownloadFile(client *http.Client, url, filepath string) error {

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("error: GET, %s", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error: READ, %s", err)
		return err
	}

	if err := ioutil.WriteFile(filepath, body, 0644); err != nil {
		log.Printf("error: WRITE, %s", err)
		return err
	}

	return nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
