package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
	"github.com/gin-contrib/cors"
	"context"
)


var (
	redisClient *redis.Client
	ctx = context.Background()

)

func main(){
	router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowAllOrigins: true,
        AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:   []string{"Content-Length"},
    }))

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})
	godotenv.Load()

	db,_ := dbSession()

	// protected Endpoints
	protectedEnpoints := router.Group("")
	protectedEnpoints.Use(authentication)
	protectedEnpoints.POST("/protected",protected)
	protectedEnpoints.POST("/add-job",addJob(db))
	protectedEnpoints.GET("/jobs",getJobs(db))
	protectedEnpoints.DELETE("/delete-job/:jobId",deleteJob(db))
	protectedEnpoints.GET("/view-applicants/:jobId",seeApplicants(db))
	protectedEnpoints.POST("/view-top-applicants/:jobId",scoringApplicants(db))
	protectedEnpoints.POST("/view-applicant-evaluation/:jobId/:resumeName",applicantDetails(db))


	router.POST("/apply/:jobId",Apply(db))
	router.GET("/job-details/:jobId",getJobDetails(db))
	router.POST("/forgot-password",forgotPassword(db))
	router.POST("/verify-otp",verifyOtp(db))
	router.GET("/verify-account/:userId",verifyUser(db))
	router.POST("/sign-up",register(db))
	router.POST("/login",login(db))
	
	router.Run("localhost:8080")
}
