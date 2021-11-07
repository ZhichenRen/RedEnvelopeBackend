package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       int64
	CurCount int `json:"cur_count"`
	Amount   int `json:"amount"`
}

// find user by user id
func getUser(uid int64) (user User) {
	dsn := "group9:Group9@haha@tcp(124.238.238.165:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} // Migrate the schema
	db.FirstOrCreate(&user, User{ID: uid})
	return
}

// create a user
func createUser(user User) {
	dsn := "group9:Group9@haha@tcp(124.238.238.165:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Create(&user)
	return
}

// update CurCount, concretely, user grab a red envelope
func updateCurCount(user *User) {
	dsn := "group9:Group9@haha@tcp(124.238.238.165:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	user.CurCount++
	db.Model(&user).Update("cur_count", user.CurCount)
}

// update amount, concretely, user grab a red envelope
func updateAmount(user *User, money int) {
	dsn := "group9:Group9@haha@tcp(124.238.238.165:3306)/red_envelope?charset=utf8&parseTime=True&loc=Local&timeout=10s"
	// connect to mysql
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	user.CurCount++
	user.Amount += money
	db.Model(&user).Update("amount", user.Amount)
}

// combine updateAmount with updateCurCount
// connect to the database twice?
// TODO
// optimize it
// delete two function above?
func updateUser(user *User, money int) {
	updateCurCount(user)
	updateAmount(user, money)
}
