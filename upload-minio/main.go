package main

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx := context.Background()
	endpoint := "localhost:9000"
	accessKeyID := "xiUSTzhtefmp6KCulti6"
	secretAccessKey := "sSo4YMpZfjwn5r4jXGF7VUB0ycex5Al5DBucxkcX"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	bucketName := "upload"
	location := "ya-practicum"
	_ = location
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Println("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Println("Successfully created %s\n", bucketName)
	}
	objectName := "mastering-postgresql.pdf"
	filePath := "/home/alexey/GO-PROG/test/upload/upload-minio/mastering-postgresql.pdf"
	contentType := "application/octet-stream"

	info, err := minioClient.FPutObject(ctx,
		bucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}
