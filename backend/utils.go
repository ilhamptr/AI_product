package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"github.com/gin-gonic/gin"
)

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	specialCharRegex := regexp.MustCompile(`[!@#~$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
	return specialCharRegex.MatchString(password)
}


func generateOTP() (string,error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "",err
	}
	return fmt.Sprintf("%06d", n.Int64()),err
}

func authenticated(context *gin.Context) (interface{},bool){
	val,exist := context.Get("claims")
	if !exist{
		return val,false
	}
	return val,true
}

func JobAuthorization(userId string, job Job) bool{
		return userId == job.AssignBy
}