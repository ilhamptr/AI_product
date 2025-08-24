package main

import (
	"time"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
	"os"
)

type User struct {
	ID string `gorm:"primaryKey;column:id"`
	Email string `gorm:"column:email;unique;not null"`
	Username string `gorm:"column:username;unique;not null"`
	Password string `gorm:"column:password;not null"`
	Active bool `gorm:"column:active;default:false"`
	CreatedAt time.Time
	
	Jobs []Job `gorm:"foreignKey:AssignBy;references:ID"`
}

type Job struct{
	ID string `gorm:"primaryKey;column:id"`
	JobTitle string `gorm:"column:job_title;not null"`
	Description string `gorm:"column:description;not null"`
	AssignBy string `gorm:"column:assigned_by"`
	Active bool `gorm:"column:active;default:true"`
	CreatedAt time.Time

	Assigner User `gorm:"foreignKey:AssignBy;references:ID"`
}

type Application struct{
	ID string `gorm:"primaryKey;column:id"`
	ApplicantName string `gorm:"column:applicant_name;not null"`
	JobID string `gorm:"column:job_id;not null"`
	PhoneNumber string `gorm:"column:Phone_number;not null"`
	Email string `gorm:"column:email;not null"`
	ResumeFile string `gorm:"column:resume_file;not null"`
	CreatedAt time.Time

	Job Job `gorm:"foreignKey:JobID;references:ID"`

}

type Question struct {
	ID        string `gorm:"primaryKey;column:id"`
	JobID     string `gorm:"column:job_id;not null"`
	Question  string `gorm:"column:question;not null"`
	Required  bool   `gorm:"column:required;default:true"`

	Job Job `gorm:"foreignKey:JobID;references:ID"`
}

type QuestionAnswer struct {
	ID         string `gorm:"primaryKey;column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	Answer     string `gorm:"column:answer;not null"`

	Question Question `gorm:"foreignKey:QuestionID;references:ID"`
}


func dbSession()(*gorm.DB,error){
	godotenv.Load()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
	os.Getenv("DB_HOST"),
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_NAME"),
	os.Getenv("DB_PORT"),)
	db,err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil{
		panic("Failed to connect to database");
	}	
	err = db.AutoMigrate(&User{},&Job{},&Application{},&Question{},&QuestionAnswer{})
	if err != nil{
		return nil,err
	}
	return db,nil
}


