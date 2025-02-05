package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {

	// Make sure the user has provided the correct number of arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: UploadData <file_path> <file_extension>")
		fmt.Println("Example: UploadData /path/to/daily.xml xml")
		os.Exit(1)
	}

	// This is supposed to be the path, I will make sure it is a valid path below
	fileLocation := os.Args[1]
	// This is the file type, it must be either xml or json, I am checking below
	fileType := os.Args[2]
	// This is the bucket name that I already setup in my s3
	var bucketName string = "first-project-cloud-1"

	// Check if the file type is either xml or json
	if fileType != "xml" && fileType != "json" {
		fmt.Println("File type must be xml or json")
		os.Exit(1)
	}

	// Grabbing my local ~/.aws/credentials file
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))

	// If there is an error loading the AWS config, print the error and exit
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		os.Exit(1)
	}

	s3Client := s3.NewFromConfig(config)

	// Trying to open the file at the given path
	file, err := os.Open(fileLocation)

	// If there is an error opening the file, print the error and exit
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	/*
	 * defer will delay the execution of a function until the surrounding function returns.
	 * So in this case, the file.Close() function will be called when the main function returns.
	 */
	defer file.Close()

	// Basically the file name without the path
	baseName := filepath.Base(fileLocation)

	objectKey := baseName

	var contentType *string
	switch fileType {
	case "xml":
		ct := "application/xml"
		contentType = &ct

	case "json":
		ct := "application/json"
		contentType = &ct
	}

	// Create a new S3 client
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &objectKey,
		Body:        file,
		ContentType: contentType,
		Metadata: map[string]string{
			"x-amz-meta-datatype": fileType,
		},
	})

	// If there is an error uploading the file, print the error and exit
	if err != nil {
		fmt.Println("Error uploading file:", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully uploaded %s to s3://%s/%s\n", fileLocation, bucketName, objectKey)
}
