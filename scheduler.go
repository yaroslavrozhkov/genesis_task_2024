package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

func StartScheduler() {
	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ticker.C:
			SendExchangeRateEmails()
		}
	}
}

func SendExchangeRateEmails() {
	//apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	url := "https://api.exchangerate-api.com/v4/latest/USD?apikey=%s" //fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/USD?apikey=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to get exchange rate:", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return
	}

	// Parse the exchange rate from the response body
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("Failed to parse JSON:", err)
		return
	}

	uahRate := data["rates"].(map[string]interface{})["UAH"].(float64)

	var subscribers []Subscriber
	DB.Find(&subscribers)

	for _, subscriber := range subscribers {
		sendEmail(subscriber.Email, uahRate)
	}
}

func sendEmail(to string, rate float64) {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := []byte(fmt.Sprintf("Subject: Daily USD to UAH Exchange Rate\n\nCurrent exchange rate: 1 USD = %.2f UAH", rate))

	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), auth, from, []string{to}, msg); err != nil {
		log.Println("Failed to send email:", err)
	}
}
