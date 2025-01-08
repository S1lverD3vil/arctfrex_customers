package storage

import (
	"arctfrex-customers/internal/common"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client     *minio.Client
	bucketName string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucketName string) (*MinioClient, error) {

	// log.Println(accessKey)
	// log.Println(secretKey)
	minioEndpointSecured, _ := strconv.ParseBool(os.Getenv(common.MINIO_ENDPOINT_SECURED))
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: minioEndpointSecured, //false, // set to `false` if you don't want HTTPS
	})
	if err != nil {
		return nil, err
	}

	// Create a bucket if it doesn't exist.
	err = createBucket(minioClient, bucketName)
	if err != nil {
		log.Fatalln(err)
	}

	return &MinioClient{
		client:     minioClient,
		bucketName: bucketName,
	}, nil
}

// func (mc *MinioClient) UploadFile(contentType string, file multipart.File, fileName string, fileSize int64) error {
// 	_, err := mc.client.PutObject(context.Background(), mc.bucketName, fileName, file, -1, minio.PutObjectOptions{ContentType: contentType})
// 	// return err
// 	// Max retry count
// 	retries := 3
// 	for i := 0; i < retries; i++ {
// 		// Attempt to upload
// 		_, err := mc.client.FPutObject(context.Background(), mc.bucketName, fileName, filePath, minio.PutObjectOptions{})
// 		if err == nil {
// 			return nil // Upload successful
// 		}

// 		// Handle multipart error
// 		if err.Error() == "The specified multipart upload does not exist" {
// 			// Check and abort incomplete uploads
// 			err := abortIncompleteUploads(ctx, minioClient, bucketName, objectName)
// 			if err != nil {
// 				return fmt.Errorf("failed to abort incomplete uploads: %v", err)
// 			}
// 		} else {
// 			return fmt.Errorf("upload error: %v", err)
// 		}
// 	}
// 	return fmt.Errorf("max retries reached for upload")
// }

func (mc *MinioClient) UploadFile(contentType string, file multipart.File, fileName string, fileSize int64) error {
	_, err := mc.client.PutObject(context.Background(), mc.bucketName, fileName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (mc *MinioClient) DownloadFile(fileName string) ([]byte, error) {
	// Get the object
	obj, err := mc.client.GetObject(context.Background(), mc.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	// Get object info to determine the size
	objInfo, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	if objInfo.Size == 0 {
		return nil, fmt.Errorf("object is empty")
	}

	// Read the object data into a byte slice
	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (mc *MinioClient) PreviewFile(fileName string) ([]byte, error) {
	// Get the object
	obj, err := mc.client.GetObject(context.Background(), mc.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	// Get object info to determine the size
	objInfo, err := obj.Stat()
	if err != nil {
		return nil, err
	}

	if objInfo.Size == 0 {
		return nil, fmt.Errorf("object is empty")
	}

	data := make([]byte, objInfo.Size)
	_, err = io.ReadFull(obj, data) // Use ReadFull to read the entire object
	if err != nil {
		return nil, err
	}

	return data, nil
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

// func  abortIncompleteUploads(ctx context.Context, minioClient *minio.Client, bucketName, objectName string) error {
// 	incompleteUploads, err := mc.client.ListIncompleteUploads(ctx, bucketName, objectName, "", true)
// 	if err != nil {
// 		return err
// 	}

// 	// Abort each incomplete upload
// 	for _, upload := range incompleteUploads {
// 		err = minioClient.AbortMultipartUpload(ctx, bucketName, upload.Key, upload.UploadID)
// 		if err != nil {
// 			return fmt.Errorf("failed to abort upload %s: %v", upload.UploadID, err)
// 		}
// 		fmt.Printf("Aborted incomplete upload: %s\n", upload.UploadID)
// 	}
// 	return nil
// }
