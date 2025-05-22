package core

import (
	"github.com/romitou/insatutorat/database/models"
	"math"
	"time"
)

func AvailabilityScore(slotsA models.Slots, slotsB models.Slots) float64 {
	score := 0.0
	for day := time.Monday; day <= time.Friday; day++ {
		dayA := slotsA[day]
		dayB := slotsB[day]
		if dayA == nil || dayB == nil {
			continue
		}
		for slot := M1; slot <= A4; slot++ {
			slotA := dayA[slot]
			slotB := dayB[slot]

			availA := availabilityValue(slotA)
			availB := availabilityValue(slotB)

			// c.f. rapport de TIP pour explications
			score += availA * availB
		}
	}
	return score
}

// c.f. rapport de TIP pour explications
func availabilityValue(slot int) float64 {
	return math.Exp(-1.0 / 8.0 * float64(slot)) // exp(-1/8 * slot)
}
