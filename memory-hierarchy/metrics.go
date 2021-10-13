package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

func AverageAge(ages []uint8) float64 {
	count := len(ages)
	limit := count - 3
	average0, average1, average2, average3 := uint64(0), uint64(0), uint64(0), uint64(0)
	i := 0
	for ; i < limit; i += 4 {
		average0 += uint64(ages[i])
		average1 += uint64(ages[i+1])
		average2 += uint64(ages[i+2])
		average3 += uint64(ages[i+3])
	}
	for ; i < count; i++ {
		average0 += uint64(ages[i])
	}
	return float64(average0+average1+average2+average3) / float64(count)
}

func AveragePaymentAmount(payments []float32) float64 {
	count := len(payments)
	average := 0.0
	for _, payment := range payments {
		average += float64(payment)
	}
	return average / float64(count)
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(payments []float32) float64 {
	mean := AveragePaymentAmount(payments)
	count := len(payments)
	squares := 0.0
	for _, payment := range payments {
		amount := float64(payment)
		squares += amount * amount
	}
	return math.Sqrt(squares/float64(count) - mean*mean)
}

func LoadData() ([]uint8, []float32) {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	ages := make([]uint8, len(userLines))
	for i, line := range userLines {
		age, _ := strconv.Atoi(line[2])
		ages[i] = uint8(age)
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	payments := make([]float32, len(paymentLines))
	for i, line := range paymentLines {
		paymentCents, _ := strconv.Atoi(line[0])
		payments[i] = float32(paymentCents) / 100.0
	}

	return ages, payments
}
