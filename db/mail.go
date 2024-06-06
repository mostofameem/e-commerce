package db

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"math/big"
	"net/smtp"
)

// Sends a verification email with the secret code
func sendVerificationEmail(email, secretCode string) error {
	from := "mostofameem@gmail.com"
	password := "rrhbslzsveidepmb"
	to := email
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Read and parse the HTML template
	tmpl, err := template.ParseFiles("/home/mostofa/ecommerce/mail.html")
	if err != nil {
		return err
	}

	// Create a buffer to hold the executed template
	var body bytes.Buffer
	err = tmpl.Execute(&body, map[string]string{"SecretCode": secretCode})
	if err != nil {
		return err
	}

	message := []byte("Subject: Email Verification\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		body.String())

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
}

// Generates a random secret code
func generateSecretCode() string {
	const min = 10000
	const max = 99999

	// Calculate the range size
	rangeSize := max - min + 1

	// Generate a random number in the range 0 to (rangeSize - 1)
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(rangeSize)))
	if err != nil {
		// Handle error
		fmt.Println("Error generating random number:", err)
		return ""
	}

	// Add the min value to shift the range to [min, max]
	code := int(nBig.Int64()) + min

	return fmt.Sprintf("%05d", code)
}
