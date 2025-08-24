package main

import (
	"net/http"
	"os"
	"time"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	// "github.com/resend/resend-go/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type ForgotPassword struct{
	Email string `json:"email"`
}

type VerifyOTP struct{
	Email string `json:"email"`
	Username string `json:"username"`
	Otp string `json:"otp"`
	NewPassword string `json:"new_password"`
}

func sendOtp(email,otp string)error{
	from := mail.NewEmail("Example User", "ilhamptr007@gmail.com")
	subject := "Your OTP"
	to := mail.NewEmail("Example User", email)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := fmt.Sprintf("<strong>your otp is: %v</strong>",otp)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
        // client.Request, _ = sendgrid.SetDataResidency(client.Request, "eu")
        // uncomment the above line if you are sending mail using a regional EU subuser
	_, err := client.Send(message)
	if err !=nil{
		return fmt.Errorf("email can't be sent")
	}
	return nil

	// resend_client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	// params := &resend.SendEmailRequest{
    //     From: "Acme <onboarding@resend.dev>",
	// 	To: []string{email},
	// 	Html:    fmt.Sprintf("<strong>your otp is: %v</strong>",otp),
	// 	Subject: "your otp",

	// }
	// _,err := resend_client.Emails.Send(params)
	// if err !=nil{
	// 	return fmt.Errorf("email can't be sent")
	// }
	// return nil
} 



func forgotPassword(db *gorm.DB)gin.HandlerFunc{
	
	return func(c *gin.Context)  {
		var input ForgotPassword
		var user User
		if err := c.BindJSON(&input); err != nil{
			c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"invalid json input"})
			return
		} 
		if input.Email == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "email is required"})
			return
		}

		result := db.Where("email = ?",input.Email).First(&user)

		if result.RowsAffected < 1{
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if result.Error != nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		key := "otp:" + input.Email
		otp,err := generateOTP()
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to generate otp"})
			return 
		}
		expiration := 5*time.Minute
		err  = redisClient.Set(ctx,key,otp,expiration).Err()
		if err !=nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to store otp"})
			return
		}
		err = sendOtp(input.Email,otp)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to send otp"})
			return
		}
		c.IndentedJSON(http.StatusOK,gin.H{"message":"otp sent successfully"})
		
	}
}

func verifyOtp(db *gorm.DB)gin.HandlerFunc{
	return func(c *gin.Context){
		var verify VerifyOTP

		if err := c.BindJSON(&verify); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		otpKey := "otp:" + verify.Email
		val,err:= redisClient.Get(ctx,otpKey).Result()
		if err == redis.Nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "OTP not found or expired"})
			return
		} else if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		if val != verify.Otp{
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "invalid otp"})
			return
		}

		validPassword := isValidPassword(verify.NewPassword)
		if !validPassword{
			c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters long and contain at least one special character.",})
			return
		}

		var userData User
		if err := db.Where("email = ? OR username = ?",verify.Email,verify.Username).First(&userData).Error; err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"account not found"})
			return
		}

		hashedPassword,err := bcrypt.GenerateFromPassword([]byte(verify.NewPassword),bcrypt.DefaultCost)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError,gin.H{"error":"failed to hash password"})
			return
		}
		
		userData.Password = string(hashedPassword)
		db.Save(&userData)

		c.IndentedJSON(http.StatusOK, gin.H{"message": "password has been updated"})
	}
}
