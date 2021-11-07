package main

type User struct {
	ID       string `gorm:"size:10"`
	CurCount int    `json:"curCount"`
	Amount   int    `json:"amount"`
}

func main() {
	// set the config

	//createUser

	// Create
	//db.Create(&User{userId: "123", curCount: 100, amount: 123})

	/*
		if err != nil {
			panic("Failed, error=" + err.Error())
		}

		//db.Model(&Envelope{}).Where("value = ?", 12).Update("value", 32)
		//envelope := db.Model(&Envelope{}).Where("envelopeId = ?", "123")
		var envelopes []*Envelope
		conditions := map[string]interface{}{
			"envelopeId": "123",
		}
		db.Table(Envelope{}.TableName()).Where(conditions).Find(&envelopes)
		snatchTime := time.Now().UnixNano()
		var envelope Envelope
		envelope = Envelope{envelopeId: "456", opened: 0, value: 10, snatchTime: snatchTime}
		db.Create(&envelope)
	*/
}
