package driver

import (
	"context"
	"io"
	"strings"

	"cloud.google.com/go/storage"
)

const (
	projectID = "lyra-proj"
)

func auth() (*storage.Client, error) {
	ctx := context.Background()
	return storage.NewClient(ctx)
}

func createDocument(client *storage.Client, bucketName string, docName string, json string) error {
	_, _, err := upload(context.Background(),
		strings.NewReader(json),
		projectID,
		bucketName,
		docName,
		false)
	return err
}

// lifted from https://github.com/GoogleCloudPlatform/golang-samples/blob/master/storage/gcsupload/gcsupload.go
func upload(ctx context.Context, r io.Reader, projectID, bucket, name string, public bool) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	bh := client.Bucket(bucket)

	// Next check if the bucket exists
	if _, err = bh.Attrs(ctx); err != nil {
		return nil, nil, err
	}

	obj := bh.Object(name)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, r); err != nil {
		return nil, nil, err
	}
	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if public {
		if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return nil, nil, err
		}
	}

	attrs, err := obj.Attrs(ctx)
	return obj, attrs, err
}
