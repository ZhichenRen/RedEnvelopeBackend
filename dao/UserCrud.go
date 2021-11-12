package dao

import "fmt"

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

func UpdateCurCount(uid int64) (int, error) {
	tx := _db.Begin()
	//defer func() {
	//	if r := recover(); r != nil {
	//		tx.Rollback()
	//	}
	//}()
	if err := tx.Error; err != nil {
		return 0, err
	}

	user := User{}
	err := tx.First(&user, User{ID: uid}).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, uid).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	user.CurCount++
	tx.Model(&user).Update("cur_count", user.CurCount)
	fmt.Println("Current count of user ", uid, ": ", user.CurCount)
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return user.CurCount, nil
}

// update amount, concretely, user grab a red envelope
func updateAmount(user *User, money int) {
	user.Amount += money
	_db.Model(&user).Update("amount", user.Amount)
}

// combine updateAmount with updateCurCount
func updateUser(user *User, money int) {
	updateCurCount(user)
	updateAmount(user, money)
}
