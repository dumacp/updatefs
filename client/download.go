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
	// Create the file
	// out, err := os.Create(filepath)
	// if err != nil {
	// 	return err
	// }
	// defer out.Close()

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

	// // Write the body to file
	// _, err = io.Copy(out, resp.Body)
	// if err != nil {
	// 	log.Printf("error: COPY, %s", err)
	// 	return err
	// }

	// buf := make([]byte, 1024*32)
	// for {
	// 	n, err := resp.Body.Read(buf)
	// 	if err != nil {
	// 		if err != io.EOF && err != io.ErrUnexpectedEOF {
	// 			log.Printf("error: READ, %s", err)
	// 			return err
	// 		}
	// 	}
	// 	if n <= 0 {
	// 		break
	// 	}
	// 	// write a chunk
	// 	if _, err := out.Write(buf[:n]); err != nil {
	// 		log.Printf("error: WRITE, %s", err)
	// 		return err
	// 	}
	// }

	// //Create the file
	// out, err := zip.NewWriter()
	// if err != nil {
	// 	return err
	// }
	// defer out.Close()

	// if body, err := ioutil.ReadAll(resp.Body); err != nil {

	// 	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	for _, zipFile := range zipReader.File {
	// 		unzippedFileBytes, err := readZipFile(zipFile)

	// 	}

	// 	_, err = io.Copy(out, zipReader.)
	// 	if err != nil {
	// 		log.Printf("error: COPY, %s", err)
	// 		return err
	// 	}
	// }

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
