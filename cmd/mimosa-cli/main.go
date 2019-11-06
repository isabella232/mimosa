package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
)

// Code to attach custom claims to Identity Platform users
//
// To run you must specifying GOOGLE_CLOUD_PROJECT and GOOGLE_APPLICATION_CREDENTIALS. GOOGLE_APPLICATION_CREDENTIALS must be a service account
//
// export GOOGLE_CLOUD_PROJECT=mimosa-256008
// export GOOGLE_APPLICATION_CREDENTIALS=/Users/scott/Desktop/mimosa-256008-76ba4ff12eee-compute-engine-default.json

func main() {

	// User IDs taken from: https://console.cloud.google.com/customer-identity/users?project=mimosa-256008
	// FIMXE Should really be parameterized via command line args
	users := map[string]map[string]interface{}{
		//"alice@example.com":
		"OdudPCfFz4TvOjuhPEDGl8IAv6s2": map[string]interface{}{
			"org": "org001",
			"cap": []string{"run"},
		},
		//"bob@example.com":
		"t7c0YM9OvAgvodcCcO5w9ClfKcF2": map[string]interface{}{
			"org": "org001",
			"cap": []string{},
		},
		// "charlie@example.com":
		"G5sDnBYL14MRDkkdEUxu3XTHWUG2": map[string]interface{}{
			"org": "org002",
			"cap": []string{"run"},
		},
		// "dervla@example.com":
		"8AaMoo1d11RpNYvAP5XtO02iZCX2": map[string]interface{}{
			"org": "org002",
			"cap": []string{},
		},
	}

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	for user, claims := range users {
		err = client.SetCustomUserClaims(ctx, user, claims)
		if err != nil {
			log.Fatalf("error setting custom claims %v\n", err)
		}
	}

}
