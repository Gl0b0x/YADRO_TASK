package pkg

import (
	"math"
	"time"
)

type Computer struct {
	Number    int
	IsBusy    bool
	TimeStart time.Time
	TotalTime time.Time
	Profit    int
}

func (com *Computer) CountProfit(timeEnd time.Time, price int) {
	com.IsBusy = false
	duration := timeEnd.Sub(com.TimeStart)
	durationHours := math.Ceil(duration.Hours())
	com.Profit += int(durationHours) * price
	com.TotalTime = com.TotalTime.Add(duration)
}
