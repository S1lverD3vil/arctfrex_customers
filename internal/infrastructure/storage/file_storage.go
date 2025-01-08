package infrastructure

import (
	"arctfrex-customers/internal/common"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewFileStorage(engine *gin.Engine) {

	endpoint := common.MINIO_ENDPOINT    // e.g. "localhost:9000"
	accessKey := common.MINIO_ACCESS_KEY // e.g. "admin"
	secretKey := common.MINIO_SECRET_KEY // e.g. "admin123"
	bucket := common.MINIO_BUCKET_NAME
	// Initialize MinIO client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Set this to true if you're using HTTPS
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Create a bucket if it doesn't exist.
	err = createBucket(minioClient, bucket)
	if err != nil {
		log.Fatalln(err)
	}
}

// Create a bucket if it doesn't exist
func createBucket(client *minio.Client, bucketName string) error {
	ctx := context.Background()

	// Check if the bucket already exists.
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("Bucket %s already exists.\n", bucketName)
		return nil
	}

	// Create a new bucket.
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		return fmt.Errorf("could not create bucket: %v", err)
	}
	fmt.Printf("Successfully created bucket %s.\n", bucketName)
	return nil
}

// Upload a file to MinIO
func uploadFile(client *minio.Client, bucketName, objectName string, content []byte) error {
	ctx := context.Background()

	// Convert the byte slice to an io.Reader
	reader := bytes.NewReader(content)

	// Upload the file
	uploadInfo, err := client.PutObject(ctx, bucketName, objectName, reader, int64(len(content)), minio.PutObjectOptions{
		ContentType: "application/text",
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	fmt.Printf("Uploaded %s of size %d successfully.\n", objectName, uploadInfo.Size)
	return nil
}

// Retrieve a file from MinIO
func retrieveFile(client *minio.Client, bucketName, objectName string) ([]byte, error) {
	ctx := context.Background()

	// Retrieve the object
	object, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file: %v", err)
	}
	defer object.Close()

	// Read the object content
	content, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %v", err)
	}

	return content, nil
}
