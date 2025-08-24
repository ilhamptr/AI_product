// using SendGrid's Go Library
// https://github.com/sendgrid/sendgrid-go
package main

import (
	// "fmt"
	// "log"
	// "os"

	// "github.com/sendgrid/sendgrid-go"
	// "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// func main() {
	// from := mail.NewEmail("Sender", "ilhamptr007@gmail.com")
	// subject := "Sending with SendGrid is Fun"
	// to := mail.NewEmail("Receiver", "iamptr007@gmail.com")
	// plainTextContent := "and easy to do anywhere, even with Go"
	// htmlContent := "<strong>your otp code i 6666</strong>"
	// message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
    //     // client.Request, _ = sendgrid.SetDataResidency(client.Request, "eu")
    //     // uncomment the above line if you are sending mail using a regional EU subuser
	// response, err := client.Send(message)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// }
// }