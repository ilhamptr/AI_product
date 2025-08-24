package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"encoding/json"

)

type ScoringInput struct{
	JobDesc string `json:"jobDescription" binding:"required"`
	TopNumber int `json:"topNumber" binding:"required"`
}

type EvaluationInput struct{
	JobDesc string `json:"jobDescription" binding:"required"`
	ResumeName string `json:"resumeName" binding:"required"`
}

func submitResume(fileBytes []byte,jobID,applicantName,filename string) ([]byte,error) {
	// change the domain in production
	endpoint := fmt.Sprintf("http://127.0.0.1:8000/resume-processing/%v",jobID)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileField,err := writer.CreateFormFile("file",filename)
	if err != nil{
		return nil,err
	}
	if _, err := fileField.Write(fileBytes); err != nil {
        return nil,err
    }
	_ = writer.WriteField("applicantName",applicantName)
	_ = writer.WriteField("resumeFile",filename)
	writer.Close()
	req,_ := http.NewRequest("POST",endpoint,body)
	req.Header.Set("Content-Type",writer.FormDataContentType())
	client := &http.Client{}
	resp,err := client.Do(req)
	if err != nil{
		return nil,err
	}
	defer resp.Body.Close()
	respBody,_ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated{
		return nil, fmt.Errorf("upload failed: %s", resp.Status)
	}
	return respBody,nil

}

func getScoring(token,jobDesc,jobId string, topNumber int) (map[string]interface{},error){
	// change the domain in production
	endpoint := fmt.Sprintf("http://127.0.0.1:8000/get-applicants/%v",jobId)

	data := ScoringInput{
		JobDesc: jobDesc,
		TopNumber: topNumber,
	}

	jsonData,_ := json.Marshal(data)
	req,_ := http.NewRequest("POST",endpoint,bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	resp,err:= client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes,_ := io.ReadAll(resp.Body)
	var result map[string]interface{}

	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("upload failed: %s", resp.Status)
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}

	return result,nil
}

func scoringDetails(token,jobDesc,resumeName string) (map[string]interface{},error){
	// change the domain in production
	endpoint := "http://127.0.0.1:8000/matching-process/"

	data := EvaluationInput{
		JobDesc: jobDesc,
		ResumeName: resumeName,
	}

	jsonData,_ := json.Marshal(data)
	req,_ := http.NewRequest("POST",endpoint,bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	resp,err:= client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes,_ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("upload failed: %s", resp.Status)
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}

	return result,nil
}