package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	localCertDir = "/usr/local/share/ca-certificates/"
)

func loadLocalCert() *http.Transport {

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := ioutil.ReadDir(localCertDir)
	if err != nil {
		log.Printf("Failed to append %q to RootCAs: %v", localCertDir, err)
	} else {
		for _, cert := range certs {
			file, err := ioutil.ReadFile(localCertDir + cert.Name())
			if err != nil {
				log.Fatalf("Failed to append %q to RootCAs: %v", cert, err)
			}
			// Append our cert to the system pool
			if ok := rootCAs.AppendCertsFromPEM(file); !ok {
				log.Println("No certs appended, using system certs only")
			}
		}
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		//InsecureSkipVerify: *insecure,
		RootCAs: rootCAs,
	}
	tr := &http.Transport{}
	tr.TLSClientConfig = config
	return tr
	/**
	client := &http.Client{Transport: tr}

	// Uses local self-signed cert
	req := http.NewRequest(http.MethodGet, "https://localhost/api/version", nil)
	resp, err := client.Do(req)
	// Handle resp and err

	// Still works with host-trusted CAs!
	req = http.NewRequest(http.MethodGet, "https://example.com/", nil)
	resp, err = client.Do(req)
	// Handle resp and err
	/**/
}
