package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/dumacp/updatefs/loader"
)

const (
	apiversion  = "api/v1"
	apifiledata = "filedata"
)

//NewRequestFilesByDevicename function valid a complete recorrido in metroplo WS
func NewRequestFilesByDevicename(urlin, devicename string, date, limit, skip int) (*[]loader.FileData, error) {

	urlGet := fmt.Sprintf("%s/%s/%s/device/%s", urlin, apiversion, apifiledata, devicename)

	params := url.Values{}
	params.Set("date", fmt.Sprintf("%d", date))
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("skip", fmt.Sprintf("%d", skip))

	urlv, err := url.Parse(urlGet)
	if err != nil {
		return nil, err
	}
	urlv.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", urlv.String(), nil)
	if err != nil {
		return nil, err
	}

	log.Printf("request: %v\n", req)

	tr := loadLocalCert()
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	fmt.Printf("body: %s\n", body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("request error: %v, %s", resp.StatusCode, body)
		return nil, nil
	}

	var store []loader.FileData
	if err := json.Unmarshal(body, &store); err != nil {
		return nil, err
	}

	return &store, nil

	//return http.DefaultClient.Do(req)
}
