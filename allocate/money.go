package allocate

import (
	"math"
	"math/rand"
	"time"
)

//var curAllocate int64
//var budget int64
//var total int64
//var moneyLeft int64
//var envelopeLeft int64
var deviation float64
var minAllocate float64
var maxAllocate float64
var flag bool
var rand0 float64
var rand1 float64

func Init() {
	// TODO
	// read from config?
	//budget = 1000000
	//total = 5000
	deviation = 0.25
	//curAllocate = 0
	//moneyLeft = budget
	//envelopeLeft = total
	flag = false
	rand.Seed(time.Now().Unix())
}

func MoneyAllocate(moneyLeft, envelopeLeft int64) (value int) {
	mu := float64(moneyLeft) / float64(envelopeLeft)
	minAllocate = (1 - deviation) * mu
	maxAllocate = (1 + deviation) * mu
	sigma := (mu - minAllocate) / 2
	value = int(standardNormal(mu, sigma))
	if int64(value) > moneyLeft {
		value = int(moneyLeft)
	}
	return
}

func standardNormal(mu float64, sigma float64) (gauss float64) {
	if flag {
		gauss = rand1*sigma + mu
		flag = false
		return
	} else {
		u1 := rand.Float64()
		u2 := rand.Float64()
		rand0 = math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)
		rand1 = math.Sqrt(-2.0*math.Log(u1)) * math.Sin(2*math.Pi*u2)
		flag = true
		return rand0*sigma + mu
	}
}
