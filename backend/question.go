package main

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

type QuestionInput struct{
	Question string `json:"question"`
	Required bool `json:"required"`
}


// func addQuestions(db *gorm.DB)gin.HandlerFunc{
// 	return func(c *gin.Context){
// 		id := uuid.New()
// 		jobId := c.Param("jobId")
// 		value,exist := c.Get("claims")
// 		if !exist{
// 			c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"authentication failed, please try again"})
// 			return
// 		}
// 		var question QuestionInput
// 		if err := c.BindJSON(&question); err != nil{
// 			c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"invalid json"})
// 			return
// 		}
// 		questionData := Question{ID: id.String(),Question: question.Question,Required: question.Required}
// 		db.Create(&questionData)

// 		c.IndentedJSON(http.StatusCreated,gin.H{"message":"created"})
// 	}

// }