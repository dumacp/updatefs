package main

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func DownloadFile(client *http.Client, url, filepath string) (string, error) {

	log.Println("Downloading")

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("error: GET, %s", err)
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Received non 200 response code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error: READ, %s", err)
		return "", err
	}

	// Create blank file
	//file, err := os.Create(filepath)
	//file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	return "", err
	//}

	//buff := make([]byte, 4096)

	//for {
	//	nread, err := resp.Body.Read(buff)
	//	if nread == 0 {
	//		if err == nil {
	//			continue
	//		}
	//		if _, err := file.Write(buff); err != nil {
	//			break
	//		}
	//		if err != nil {
	//			break
	//		}
	//	}
	//}
	////size, err := io.Copy(file, resp.Body)
	//defer file.Close()

	////if err != nil {
	////	log.Printf("error: WRITE, %s", err)
	////	return "", err
	////}
	//statsFile, err := file.Stat()
	//log.Printf("Downloaded a file %s with size %d", filepath, statsFile.Size())

	if err := ioutil.WriteFile(filepath, body, 0644); err != nil {
		log.Printf("error: WRITE, %s", err)
		return "", err
	}

	//content, err := ioutil.ReadAll(file)
	//if err != nil {
	//	log.Printf("error: READ, %s", err)
	//	return "", err
	//}

	md5sum := md5.Sum(body)
	md5s := hex.EncodeToString(md5sum[0:])

	return md5s, nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
