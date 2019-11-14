package usermgmt

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"

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

func generateWorkspaceID() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
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
