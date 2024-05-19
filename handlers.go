package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	subscriber := Subscriber{Email: email}

	var mysqlErr *mysql.MySQLError

	if err := DB.Create(&subscriber).Error; err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			http.Error(w, "Email already subscribed", http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Subscription successful"))
}

func CurrencyHandler(w http.ResponseWriter, r *http.Request) {
	rate, err := GetUAHExchangeRate()
	if err != nil {
		http.Error(w, "Failed to get exchange rate", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%.2f", rate)
}
