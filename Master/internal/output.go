package internal

import (
	"Master/constants"
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func SaveOutputOnS3() {
	var config constants.Config
	constants.ReadJsonConfig(&config)
	bucketRegion := config.Region
	bucketName := config.Bucket

	// Name of output folder
	outputFolder := "output"
	// Name of created zip file
	// ROME time zone (CET)
	romeLocation, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		log.Fatalf("Error during time zone acquisition: %v", err)
		return
	}

	currentTime := time.Now().In(romeLocation)

	// Converts timestamp to string format
	timeString := currentTime.Format("2006-01-02_15-04-05")
	zipFileName := fmt.Sprintf("output_%s.zip", timeString)

	// Zip output folder
	err = zipFolder(outputFolder, zipFileName)
	if err != nil {
		log.Fatalf("Error during file zipping: %v", err)
		return
	}

	// Save zip file on S3
	err = uploadToS3(zipFileName, bucketName, bucketRegion)
	if err != nil {
		log.Fatalf("Error during s3 upload: %v", err)
		return
	}

	log.Println("Completed!")
}

func zipFolder(sourceFolder, zipFileName string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			log.Fatalf("Error during file closing: %v", err)
			return
		}
	}(zipFile)

	archive := zip.NewWriter(zipFile)
	defer func(archive *zip.Writer) {
		err := archive.Close()
		if err != nil {
			log.Fatalf("Error during archive closing: %v", err)
			return
		}
	}(archive)

	err = filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		relPath, err := filepath.Rel(sourceFolder, path)
		if err != nil {
			return err
		}

		zipEntry, err := archive.Create(relPath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("Error during file closing: ", err)
				return
			}
		}(file)

		_, err = io.Copy(zipEntry, file)
		return err
	})

	return err
}

func uploadToS3(fileName, bucketName, bucketRegion string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(bucketRegion),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error during s3 upload: %v", err)
			return
		}
	}(file)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	return err
}
