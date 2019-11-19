package usermgmt

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/firestore"

	firebase "firebase.google.com/go"
)

type userDocument struct {
	Workspaces map[string]string `firestore:"workspaces"`
}

type workspaceDocument struct {
	Name string `firestore:"name"`
}

type userRecord struct {
	Email        string                 `json:"email"`
	UID          string                 `json:"uid"`
	CustomClaims map[string]interface{} `json:"customClaims"`

	// Supported fields
	//
	// disabled
	// displayName
	// email
	// emailVerified
	// metadata
	// passwordHash
	// passwordSalt
	// phoneNumber
	// photoURL
	// providerData
	// tenantId
	// tokensValidAfterTime
	// uid

}

//
// This function generates a random workspace ID
//
// We have 3 requirements:
// * We want the IDs to be short because we want to fit as many as possible into a JWT custom claim, which is limited to 1000 bytes
// * We want the ID space to be as large as possible so that we don't run out if we have millions of workspaces
// * Workspace IDs are also used as bucket labels in GCS and so must be lower case (ruling out base 64 encoding)
//
// Our solution is to use 25-bit IDs encoded as 5 base-32 characters
//
func generateWorkspaceID() (string, error) {
	// We generate a cryptographically random number into a byte array
	// We only need 25 bits but that still means an array of length 4 to fit it all in
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Encode the whole 32 bits into base-32
	encoded := base32.StdEncoding.EncodeToString(b)
	// We only need 25 bits so let's take the first 5 base-32 chars and discard the rest
	encoded = encoded[:5]
	// Base-32 is encoded in upper case, so we convert to lower case
	encoded = strings.ToLower(encoded)
	// We're done!
	return encoded, err
}

// UserCreated responds to creation of Identity Platform users
func UserCreated(ctx context.Context, user userRecord) error {

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	// Get the Firestore client
	project := os.Getenv("GCP_PROJECT")
	if len(project) == 0 {
		project = firestore.DetectProjectID
	}
	fc, err := firestore.NewClient(ctx, project)
	if err != nil {
		return err
	}

	// Check for a "user" document in Firestore
	log.Printf("Checking for user document ...")
	_, err = fc.Collection("users").Doc(user.UID).Get(ctx)
	if err == nil {
		// No error therefore the user exists and we don't need to create a default workspace for them
		return nil
	}
	// If this error is not just indicating the doc doesn't exist yet then there is a problem
	if status.Code(err) != codes.NotFound {
		return err
	}

	// Create a default workspace
	workspace, err := generateWorkspaceID()
	if err != nil {
		return fmt.Errorf("error generating workspace ID: %v", err)
	}
	workspaceName := "Default Workspace for " + user.Email

	// Write the "user" document to Firestore
	userDocument := userDocument{
		Workspaces: map[string]string{workspace: workspaceName},
	}
	log.Printf("Writing user document: %v", userDocument)
	_, err = fc.Collection("users").Doc(user.UID).Set(ctx, userDocument)
	if err != nil {
		return err
	}

	// Write the "workspace" document to Firestore
	log.Printf("Writing workspace document ...")
	_, err = fc.Collection("ws").Doc(workspace).Set(ctx, workspaceDocument{Name: workspaceName})
	if err != nil {
		return err
	}

	// Add this workspace to the user's custom claim
	if user.CustomClaims == nil {
		user.CustomClaims = map[string]interface{}{}
	}
	if user.CustomClaims["owner"] == nil {
		user.CustomClaims["owner"] = []string{}
	}
	ownerWorkspaces := user.CustomClaims["owner"].([]string)
	ownerWorkspaces = append(ownerWorkspaces, workspace)
	user.CustomClaims["owner"] = ownerWorkspaces

	// Update the claim in Identity Platform
	log.Printf("Updating claims in Identity Platform")
	client, err := app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting auth client: %v", err)
	}
	err = client.SetCustomUserClaims(ctx, user.UID, user.CustomClaims)
	if err != nil {
		return fmt.Errorf("error setting custom claims: %v", err)
	}

	return nil
}
