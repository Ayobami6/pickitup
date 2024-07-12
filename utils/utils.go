package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/Ayobami6/pickitup/config"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"gopkg.in/gomail.v2"
)

var rdb *redis.Client
var ctx = context.Background()

var Validate = validator.New()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",        
		DB:       0,               
	})
}


func GenerateAndCacheVerificationCode(email string) (int, error) {
	rand.NewSource(time.Now().UnixNano())

	randomNumber := rand.Intn(9000) + 1000

	numberStr := fmt.Sprintf("%d", randomNumber)

	err := rdb.Set(ctx, email, numberStr, 15*time.Minute).Err()
	if err != nil {
		return 0, err
	}


	return randomNumber, nil
}

func GetCachedVerificationCode(email string) (int, error) {
	val, err := rdb.Get(ctx, email).Result()
	if err != nil {
		return 0, err
	}

	var randomNumber int
	_, err = fmt.Sscanf(string(val), "%d", &randomNumber)
	if err != nil {
		return 0, err
	}

	return randomNumber, nil
}



func WriteJSON(w http.ResponseWriter, status int, status_msg, data any, others ...string) error {
	message := ""
	if len(others) > 0 {
        message = others[0]
    }
	res := map[string]interface{}{
		"status": status_msg,
        "data":  data,
		"message": message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(res)

}


func WriteError(w http.ResponseWriter, status int, err... string) {
	var errMessage string
	if len(err) > 0 {
        errMessage = err[0]
    } else {
		errMessage = "Don't Panic This is From Us!"
	}
	log.Println(err)
    WriteJSON(w, status, "error", nil, errMessage)
}


func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return errors.New("request body is missing")
	}
	return json.NewDecoder(r.Body).Decode(payload)

}

func SendMail(recipient string, subject string, username string, message string) error {
	tmpl, err := os.ReadFile("utils/verification_template.html")
	if err != nil {
		fmt.Println("Error reading template file:", err)
		return err
	}
	
	t, err := template.New("email").Parse(string(tmpl))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}
	
    data := map[string]interface{}{
        "UserName": username,
        "Message":  message,
    }
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}
	
	m := gomail.NewMessage()

	// Set email headers
	m.SetHeader("From", "sainthaywon80@gmail.com")
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)

	// Set the HTML body
	m.SetBody("text/html", body.String())
	smtpHost := config.GetEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := 465
	smtpUser := config.GetEnv("SMTP_USER", "protected@gmail.com")
	smtpPass := config.GetEnv("SMTP_PWD", "protected")

	// Create a new SMTP dialer
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Send the email and handle errors
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	// Success message
	fmt.Println("Email sent successfully!")

	return nil

}


func GetTokenFromRequest(r *http.Request) (string, error) {
	tokenAuth := r.Header.Get("Authorization")
    tokenQuery := r.URL.Query().Get("token")

    if tokenAuth!= "" {
        return tokenAuth, nil
    }

    if tokenQuery!= "" {
        return tokenQuery, nil
    }

    return "", errors.New("token not found in request")
}