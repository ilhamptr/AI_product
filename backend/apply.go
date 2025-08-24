package main

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"fmt"
	"gorm.io/gorm"
)

type Candidate struct {
    Name       string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Email string  `json:"email"`
}

func getJobDetails(db *gorm.DB)gin.HandlerFunc{
	return func(c *gin.Context) {
		var job Job
		jobId := c.Param("jobId")
		if err:= db.Where("id = ?",jobId).First(&job).Error; err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"failed to find the job"})
			return
		}
		c.IndentedJSON(http.StatusOK,job)
	}
}

func Apply(db *gorm.DB)gin.HandlerFunc{
	return func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		phone := c.PostForm("phone_number")
		file,_ := c.FormFile("file")
		var job Job
		if name == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
			return
		}
		if phone == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
			return
		}
		if email == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}
		if file.Size == 0  {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Resume file is required"})
			return
		}
		if file.Size > 3*1024*1024{
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is too large. max size is 3MB."})
    		return
		}

		jobId := c.Param("jobId")
		if err:= db.Where("id = ?",jobId).First(&job).Error; err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"failed to find the job"})
			return
		}
		if !job.Active{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"job offer is currently inactive"})
			return
		}
		
		ApplicantID := uuid.New()
		FileName := uuid.New()
		filename := fmt.Sprintf("%s.pdf", FileName.String())
		fileBytes,err := FileProcessing(c,filename)
		if err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"failed to upload"})
			return
		}

		// later we can add the scoring algorithm function in this function ( i need more money for this unfortunately :( ))
		// resumeURL,err := generateSignedURL(filename)
		_,err =  submitResume(fileBytes,jobId,name,filename)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError,gin.H{"error":"something went wrong, try again"})
			return		
		}
		candidate := Application{ID: ApplicantID.String(),ApplicantName: name,JobID:jobId,PhoneNumber: phone,Email: email,ResumeFile: filename}
		db.Create(&candidate)
		c.IndentedJSON(http.StatusCreated,gin.H{"id":ApplicantID.String(),"name":name,"phone_number":phone,"email":email})
	}
}