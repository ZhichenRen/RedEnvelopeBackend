package dao

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