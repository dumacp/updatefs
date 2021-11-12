package main

import "flag"

var (
	keycloakurl  = "https://sibus.nebulae.com.co/auth"
	clientid     = "devices"
	clientSecret = "a367c3bc-089e-4b46-8df9-439bf814557e"
	redirecturl  = "https://sibus.nebulae.com.co/updatevoc/*"
	realm        = "DEVICES"
)

func init() {
	flag.StringVar(&keycloakurl, "keycloakUrl", "", "example: \"https://fleet.nebulae.com.co/auth\", keycloak url")
	flag.StringVar(&clientid, "clientID", "", "example: \"devices3\", clientid in realm")
	flag.StringVar(&clientSecret, "clientSecret", "", "example: \"da9bbc28-01d8-43af-8c8a-fb0654937231\", client secret")
	flag.StringVar(&redirecturl, "redirectUrl", "", "example: \"https://fleet-mqtt.nebulae.com.co/\", redirecturl url")
	flag.StringVar(&realm, "realm", "", "example: \"DEVICES\", realm name")
}
