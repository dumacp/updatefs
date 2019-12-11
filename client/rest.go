package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

//NewRequestValidarRecorrido function valid a complete recorrido in metroplo WS
func NewRequestFiles(url0, devicename string) (*http.Response, error) {

	urlGet := fmt.Sprintf("%s/%s", url0, devicename)

	urlorigin, err := url.Parse(urlGet)
	if err != nil {
		return nil, err
	}

	urlorigin.Query().Add("limit", "1")

	req, err := http.NewRequest("GET", urlorigin.RawQuery, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("request: %v", req)

	tr := loadLocalCert()
	client := &http.Client{Transport: tr}

	return client.Do(req)
	//return http.DefaultClient.Do(req)
}
