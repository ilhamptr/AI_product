package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

func upload(fileName string, fileData []byte) error {
    projectID := os.Getenv("PROJECT_ID")
    bucket := os.Getenv("BUCKET") 
    token := os.Getenv("SUPABASE_TOKEN")

    uploadURL := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/%s/%s", projectID, bucket, fileName)

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("file", fileName)
    if err != nil {
        return err
    }
    part.Write(fileData)
    writer.Close()

    req, err := http.NewRequest("POST", uploadURL, body)
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        b, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("failed to upload: %s", b)
    }

    return nil
}


func FileProcessing(c *gin.Context,fileName string) ([]byte,error){
    file, err := c.FormFile("file")
    if err != nil {
        return nil,err
    }

    openedFile, err := file.Open()
    if err != nil {
        return  nil,err
    }
    defer openedFile.Close()

    fileBytes, err := io.ReadAll(openedFile)
    if err != nil {
        return  nil,err
    }
	
    err = upload(fileName, fileBytes)
    if err != nil {
        return nil,err
    }

	return fileBytes,nil
}