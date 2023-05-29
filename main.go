package main

import (
	"finbox/org/backend/cloud-drop/uploader"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	// "path/filepath"
	// "github.com/charmbracelet/bubbles/spinner"
	// tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
)

type Result struct {
	count   int
	files   []string
	folders []string
}

type UploadPath struct {
	urls []string
}

func moveProcessedFiles(files []string, folders []string) {
	uploadRoot := []string{"./uploads/"}
	uploadCompleteRoot := []string{"./uploads-complete/"}

	for _, folder := range folders {
		err := os.MkdirAll(fmt.Sprintf("./uploads-complete/%s/", folder), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, file := range files {
		sourcePath := append(uploadRoot, file)
		destinationPath := append(uploadCompleteRoot, file)
		err := os.Rename(strings.Join(sourcePath, ""), strings.Join(destinationPath, ""))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanupDirs(folders []string) {
	for _, folder := range folders {
		err := os.Remove(fmt.Sprintf("./uploads/%s/", folder))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func processDirFiles(dir fs.DirEntry, profileName string, bucket string) (Result, UploadPath, error) {
	var builder strings.Builder
	builder.WriteString("./uploads/")
	builder.WriteString(dir.Name())
	path := builder.String()

	fmt.Printf("Processing %s: \n", path)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var result Result
	var s3Url UploadPath

	for _, file := range files {
		var filePath []string
		var uploadFilePath []string
		filePath = append(filePath, "./uploads/", dir.Name(), "/", file.Name())
		uploadFilePath = append(uploadFilePath, dir.Name(), "/", file.Name())

		uploaderStruct := &uploader.S3Client{
			Bucket:         bucket,
			OsFilePath:     strings.Join(filePath, ""),
			UploadFilePath: strings.Join(uploadFilePath, ""),
		}

		uploaderStruct.CreateClient(profileName)
		ok, err := uploaderStruct.UploadToS3()
		if err != nil {
			fmt.Println(err)
		}

		if ok {
			s3Url.urls = append(s3Url.urls, fmt.Sprintf("https://%s.s3.ap-south-1.amazonaws.com/%s/%s", bucket, dir.Name(), file.Name()))
			result.count++
			result.files = append(result.files, strings.Join(uploadFilePath, ""))

		}
	}

	return result, s3Url, nil
}

func main() {
	fmt.Println("Cloud Drop - Initialized...")
	cleanup := flag.Bool("cleanup", false, "Set true to delete uploaded files and folders")
	profile := flag.String("profile", "test-profile", "Set the aws profile to use")
	bucket := flag.String("bucket", "test-bucket", "Bucket name to uploade your assets to")

	flag.Parse()
	fmt.Printf("Using %s AWS Profile \n", *profile)

	var result Result
	var uploadPaths UploadPath
	dirs, err := os.ReadDir("./uploads")
	if err != nil {
		log.Fatal("Cannot read uploads directory. Please check if directory with the name 'uploads' exists")
	}

	if len(dirs) == 0 {
		fmt.Println("No directories found inside 'uploads'")
		fmt.Println("Cloud Drop - Done")
		return
	}

	for _, dir := range dirs {
		currentResult, currentUploadPaths, err := processDirFiles(dir, *profile, *bucket)
		if err != nil {
			fmt.Println(err)
		}

		result.count = result.count + currentResult.count
		result.files = append(result.files, currentResult.files...)
		result.folders = append(result.folders, dir.Name())
		uploadPaths.urls = append(uploadPaths.urls, currentUploadPaths.urls...)
	}
	fmt.Printf("Uploaded %d files successfully!\n", result.count)
	fmt.Printf("Uploaded Files: \n%s", strings.Join(uploadPaths.urls, "\n"))

	fmt.Println("Cloud Drop - Moving Processed files...")
	moveProcessedFiles(result.files, result.folders)
	if *cleanup {
		fmt.Println("Cloud Drop - Cleaning up empty upload directories")
		cleanupDirs(result.folders)
	}

	fmt.Println("Cloud Drop - Done")
}
