package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)
type SignUp struct{
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Login struct{
	Email string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Payload struct{
	ID string
	Email string
	Username string

}



func register(db *gorm.DB)gin.HandlerFunc{
	return func(c *gin.Context){
	var input SignUp
	if err := c.BindJSON(&input); err != nil{
		c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"invalid json"})
		return
	}

	if input.Email == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if input.Username == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}
	if input.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	validPassword := isValidPassword(input.Password)
	if !validPassword{
		c.IndentedJSON(http.StatusBadRequest, gin.H{
		"error": "Password must be at least 8 characters long and contain at least one special character.",})
		return
	}

	var UserExist User
	if err := db.Where("email = ? OR username = ?",input.Email,input.Username).First(&UserExist).Error;err == nil{
		c.IndentedJSON(http.StatusConflict,gin.H{"error":"username or email already exist"})
		return
	}
	id := uuid.New()  
	hashedPassword,_ := bcrypt.GenerateFromPassword([]byte(input.Password),bcrypt.DefaultCost)
	
	user := User{ID: id.String(),Email:input.Email,Username: input.Username,Password: string(hashedPassword)}
	db.Create(&user)

	err := sendVerifLink(user.Email,user.ID)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,gin.H{"error":"can't send verification email, try again"})
		return
	}

	c.IndentedJSON(http.StatusCreated,gin.H{
		"message": "registration success, please check your email",
	})

}
}

func token(data Payload)(string,error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": data.ID,
		"username": data.Username,
		"email": data.Email,
		// this part set the jwt expiration time
		"exp": time.Now().Add(time.Minute * 120).Unix(),
	})
	tokenSting,err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil{
		return "",err
	}
	return tokenSting,nil

}

func sendVerifLink(email,id string)error{

	from := mail.NewEmail("Sender", "ilhamptr007@gmail.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Receiver", email)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := fmt.Sprintf("<p>to verify your account please click this link below:<br> http://127.0.0.1:8080/verify-account/%v</p>",id)
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
	// 	Html:    fmt.Sprintf("<p>to verify your account please click this link below:<br> http://127.0.0.1:8080/verify-account/%v</p>",id),
	// 	Subject: "your verification link",

	// }
	// _,err := resend_client.Emails.Send(params)
	// if err !=nil{
	// 	return fmt.Errorf("email can't be sent")
	// }
	// return nil
} 

func verifyUser(db *gorm.DB)gin.HandlerFunc{
	return func (c *gin.Context)  {
		userId := c.Param("userId")
		var userData User
		if err := db.Where("id = ?", userId).First(&userData).Error; err != nil{
			c.IndentedJSON(http.StatusNotFound,gin.H{"error":"account not found"})
			return
		}
		if userData.Active{
			c.IndentedJSON(http.StatusOK,gin.H{"message":"this user is already verified"})
			return
		}
		userData.Active = true
		db.Save(&userData)
		c.IndentedJSON(http.StatusOK,gin.H{"message":"user has been verified successfully"})
	}
}

func login(db *gorm.DB)gin.HandlerFunc{
	return func(c *gin.Context) {
		var loginInput Login
		if err := c.BindJSON(&loginInput); err != nil{
			c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"invalid json input"})
			return
		} 

		var UserExist User
		if err := db.Where("email = ? OR username = ?",loginInput.Email,loginInput.Username).First(&UserExist).Error; err != nil{
			c.IndentedJSON(http.StatusUnauthorized,gin.H{"error":"account not found"})
			return
		}

		if !UserExist.Active{
			c.IndentedJSON(http.StatusUnauthorized,gin.H{"error":"please verify your account"})
			return
		}

		if err:= bcrypt.CompareHashAndPassword([]byte(UserExist.Password),[]byte(loginInput.Password)); err != nil{
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
    		return
		}

		payloadData := Payload{
			ID: UserExist.ID,
			Email: UserExist.Email,
			Username: UserExist.Username,
		}
		strToken,err := token(payloadData)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError,gin.H{"error":"failed to generate token"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{
        "token": strToken,
   		})

	}
}

func authentication(c *gin.Context){
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader,"Bearer "){
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
		c.Abort() //stop request

		return
	}
	tokenString := strings.TrimPrefix(authHeader,"Bearer ")

	claims,err := verifyToken(tokenString)
	if err != nil{
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort() //stop request
		return
	}

	c.Set("claims",claims)
	c.Next()
}

func protected(c *gin.Context){
	value,exist := c.Get("claims")
	if !exist {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
	claims := value.(jwt.MapClaims)
	
	userID := claims["id"] 
	c.IndentedJSON(http.StatusAccepted,gin.H{"message": userID})
}

func verifyToken(tokenString string)(jwt.MapClaims,error){
	token,err := jwt.Parse(tokenString,func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")),nil
	})
	if err != nil{
		return nil,err
	}
	if !token.Valid{
		return nil,fmt.Errorf("invalid token")
	}
	claims,ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil,fmt.Errorf("invalid token claim")
	}
	return claims,nil
}
