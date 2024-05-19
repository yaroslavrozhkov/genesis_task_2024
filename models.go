package main

import "gorm.io/gorm"

type Subscriber struct {
	gorm.Model
	Email string `gorm:"unique"`
}
