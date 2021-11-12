package dao

import (
	"fmt"
	"gorm.io/gorm/clause"
)

type User struct {
	ID       int64
	CurCount int `json:"cur_count"`
	Amount   int `json:"amount"`
}

// find user by user id

func GetUser(uid int64) (user User, err error) {
	err = _db.First(&user, User{ID: uid}).Error
	return
}

// create a user
func createUser(user User) {
	_db.Create(&user)
	return
}

// update CurCount, concretely, user grab a red envelope
func updateCurCount(user *User) {
	user.CurCount++
	_db.Model(&user).Update("cur_count", user.CurCount)
}

func UpdateCurCount(uid int64) error {
	tx := _db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	user := User{}
	err := tx.First(&user, User{ID: uid}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&user)
	user.CurCount++
	if err := tx.Model(&user).Update("cur_count", user.CurCount).Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}
	fmt.Println("Current count of user ", uid, ": ", user.CurCount)
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// update amount, concretely, user grab a red envelope
func updateAmount(user *User, money int) {
	user.Amount += money
	_db.Model(&user).Update("amount", user.Amount)
}

func UpdateAmount(uid int64, money int) error {
	tx := _db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	user := User{}
	err := tx.First(&user, User{ID: uid}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&user)
	if err := tx.Model(&user).Update("amount", user.Amount + money).Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}
	fmt.Println("Current amount of user ", uid, ": ", user.Amount)
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// combine updateAmount with updateCurCount
func updateUser(user *User, money int) {
	updateCurCount(user)
	updateAmount(user, money)
}
