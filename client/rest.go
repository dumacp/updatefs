package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"

	"github.com/dumacp/updatefs/loader"
)

const (
	apiversion = "updatevoc/api/v2"
	apidevice  = "files/device"
	apiupdate  = "updates"
)

//NewRequestFilesByDevicename function valid a complete recorrido in metroplo WS
func NewRequestFilesByDevicename(client *http.Client, urlin, devicename string, date, limit, skip int) ([]*loader.FileData, error) {

	urlGet := fmt.Sprintf("%s/%s/%s/%s", urlin, apiversion, apidevice, devicename)

	params := url.Values{}
	params.Set("date", fmt.Sprintf("%d", date))
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("skip", fmt.Sprintf("%d", skip))

	urlv, err := url.Parse(urlGet)
	if err != nil {
		return nil, err
	}
	urlv.RawQuery = params.Encode()

	// req, err := http.NewRequest("GET", urlv.String(), nil)
	// if err != nil {
	// 	return nil, err
	// }

	// log.Printf("request: %v\n", req)

	// tr := loadLocalCert()
	// client := &http.Client{Transport: tr}
	// defer client.CloseIdleConnections()

	resp, err := client.Get(urlv.String())
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

	store := make([]loader.FileData, 0)
	if err := json.Unmarshal(body, &store); err != nil {
		return nil, err
	}

	storeMap := make(map[int]*loader.FileData)
	storeSort := make([]*loader.FileData, 0)

	for _, v := range store {
		storeMap[int(v.Date)] = &v
	}

	keys := make([]int, 0, len(storeMap))
	for k := range storeMap {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		storeSort = append(storeSort, storeMap[k])
	}

	return storeSort, nil

	//return http.DefaultClient.Do(req)
}

//NewUpdateByDevicename function valid a complete recorrido in metroplo WS
func NewUpdateByDevicename(client *http.Client, urlin, devicename, filemd5 string, filedata string, date int) error {

	urlPost := fmt.Sprintf("%s/%s/%s", urlin, apiversion, apiupdate)

	params := url.Values{}
	params.Set("date", fmt.Sprintf("%d", date))
	params.Set("devicename", devicename)
	params.Set("filemd5", filemd5)
	params.Set("filedata", filedata)

	// tr := loadLocalCert()
	// client := &http.Client{Transport: tr}
	// defer client.CloseIdleConnections()

	resp, err := client.PostForm(urlPost, params)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	fmt.Printf("body: %s\n", body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("request error: %v, %s", resp.StatusCode, body)
		return nil
	}
	return nil

	//return http.DefaultClient.Do(req)
}
