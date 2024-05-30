package controllers

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediaController struct{}

func (h MediaController) UploadImage(c *gin.Context){
	// err := c.Request.ParseMultipartForm(10 << 20) // max 10 MB
	// if err != nil {
	// 	c.JSON(400, gin.H{"message":"server error"})
	// }
	if c.Request.ContentLength == 0 {
		c.JSON(400, gin.H{"message": "Request body is empty"})
		return
	}
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"message":"server error"})
		return
	}
	if file == nil {
		c.JSON(400, gin.H{"message":"image cannot empty"})
		return
	}
	defer file.Close()
	fileSize := handler.Size
	if fileSize <  10 * 1024 || fileSize > 20 * 1024 * 1024 {
		c.JSON(400, gin.H{"message": "File size must be between 10 KB and 20 MB"})
		return
	}
	fileExt := filepath.Ext(handler.Filename)
	if fileExt != ".jpg" && fileExt != ".jpeg" {
		c.JSON(400, gin.H{"error": "File must be a JPEG image"})
		return
	}

	fileName := uuid.New().String() + fileExt

	S3_REGION := os.Getenv("AWS_REGION")
	S3_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	S3_SECRET_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	S3_BUCKET_NAME := os.Getenv("AWS_S3_BUCKET_NAME")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(S3_ID, S3_SECRET_KEY, ""),
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create AWS session"})
		return
	}

	s3Client := s3.New(sess)
	buffer := make([]byte, handler.Size)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		c.JSON(500, gin.H{"message": "Failed to read file"})
		return
	}
	_, _, err = image.Decode(bytes.NewReader(buffer))
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to decode image"})
		return
	}

	// Upload the file to S3
	// _, err = s3Client.PutObject(&s3.PutObjectInput{
	// 	Bucket: aws.String(S3_BUCKET_NAME),
	// 	Key:    aws.String(fileName),
	// 	Body:   bytes.NewReader(buffer),
	// })
   
    _, err = s3Client.PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(S3_BUCKET_NAME),
        Key:                  aws.String(fileName),
        ACL:                  aws.String("public-read"),
        Body:                 bytes.NewReader(buffer),
    })
	if err != nil {
		 
		c.JSON(500, gin.H{"message": "Failed to upload file to S3"})
		return
	}
	url := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", S3_BUCKET_NAME, S3_REGION, fileName)

	c.JSON(200,gin.H{"message":"success", "data":gin.H{"imageUrl": url}})
}