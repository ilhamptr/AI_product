package main

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"strconv"

)

type jobInput struct{
	JobTitle string `json:"job_title"`
	Description string `json:"description"`
}


func addJob(db *gorm.DB)gin.HandlerFunc{
	return func(c *gin.Context){
		var job_input jobInput
		if err := c.BindJSON(&job_input); err != nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"invalid json"})
		return
		}
		id := uuid.New() 
		value,exist := c.Get("claims")
		if !exist{
			// c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"authentication failed, please try again"})
			return
		} 
		if job_input.JobTitle == "" || job_input.Description == ""{
			c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"please fill the job title and description field"})
			return
		}
		claims := value.(jwt.MapClaims)
		assigner := claims["id"].(string)


		job := Job{ID: id.String(),JobTitle: job_input.JobTitle,Description: job_input.Description,AssignBy: assigner,Assigner: User{ID: assigner,Username: claims["username"].(string),Email: claims["email"].(string)}}
		db.Create(&job)
		c.IndentedJSON(http.StatusCreated,job)
	}
}

func getJobs(db *gorm.DB)gin.HandlerFunc{
	return func(c *gin.Context){
		var jobs []Job
		value,exist := c.Get("claims")
		if !exist{
			c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"authentication failed, please try again"})
			return
		} 
		claims := value.(jwt.MapClaims)
		userId := claims["id"].(string)
		if err := db.Where("assigned_by = ?",userId).Find(&jobs).Error; err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"records not found"})
			return
		}
		c.IndentedJSON(http.StatusOK,jobs)

	}
}

func deleteJob(db *gorm.DB)gin.HandlerFunc{
	return func(c *gin.Context){
		jobId := c.Param("jobId")
		value,exist := c.Get("claims")
		if !exist{
			c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"authentication failed, please try again"})
			return
		} 

		claims := value.(jwt.MapClaims)
		userId := claims["id"].(string)

		var job Job
		if err := db.First(&job,"id = ?",jobId).Error; err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"record not found"})
			return 
		}

		if userId != job.AssignBy{
			c.IndentedJSON(http.StatusForbidden,gin.H{"error":"you are not allowed to delete this record"})
			return 
		}

		if err:= db.Delete(&job).Error; err!=nil{
			c.IndentedJSON(http.StatusInternalServerError,gin.H{"error":"failed to delete record"})
		}
		c.IndentedJSON(http.StatusOK,gin.H{"message":"record deleted successfully"})
	}
}

func seeApplicants(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		jobId := c.Param("jobId")
		var job Job

		val,exist := authenticated(c)
		claims := val.(jwt.MapClaims)
		userId := claims["id"].(string)

		if !exist{
			c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"authentication failed, please try again"})
			return
		} 
		if err := db.First(&job,"id = ?",jobId).Error; err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"record not found"})
			return 
		}

		authorized := JobAuthorization(userId,job)

		if !authorized{
			c.IndentedJSON(http.StatusForbidden,gin.H{"error":"you are not allowed to access this record"})
			return 
		}

		var applicants []Application
		if err := db.Where("job_id",jobId).Find(&applicants).Error; err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"records not found"})
			return
		}
		c.IndentedJSON(http.StatusOK,applicants)
	}
}

func scoringApplicants(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		jobId := c.Param("jobId")
		token := c.GetHeader("Authorization")
		query:=c.DefaultQuery("topNum","1")
		topNumber,_ := strconv.Atoi(query)
		var job Job

		val,exist := authenticated(c)
		claims := val.(jwt.MapClaims)
		userId := claims["id"].(string)

		if !exist{
			c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"authentication failed, please try again"})
			return
		} 

		if err := db.First(&job,"id = ?",jobId).Error; err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"record not found"})
			return 
		}

		authorized := JobAuthorization(userId,job)
		if !authorized{
			c.IndentedJSON(http.StatusForbidden,gin.H{"error":"you are not allowed to access this record"})
			return 
		}

		val,err := getScoring(token,job.Description,jobId,topNumber)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError,gin.H{"error":err})
			return 
		}
		c.IndentedJSON(http.StatusOK,gin.H{"result":val})
	}
}

func applicantDetails(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		jobId := c.Param("jobId")
		resumeName := c.Param("resumeName")

		val,exist := authenticated(c)
		var job Job

		claims := val.(jwt.MapClaims)
		userId := claims["id"].(string)

		if !exist{
			c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"authentication failed, please try again"})
			return
		} 
		
		if err := db.First(&job,"id = ?",jobId).Error; err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"record not found"})
			return 
		}

		authorized := JobAuthorization(userId,job)
		if !authorized{
			c.IndentedJSON(http.StatusForbidden,gin.H{"error":"you are not allowed to access this record"})
			return 
		}

		val,err := scoringDetails(token,job.Description,resumeName)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError,gin.H{"error":err})
			return 
		}
		c.IndentedJSON(http.StatusOK,gin.H{"result":val})
	}
}
