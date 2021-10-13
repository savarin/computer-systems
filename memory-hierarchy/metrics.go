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
	limit := count - 3
	average0, average1, average2, average3 := float64(0), float64(0), float64(0), float64(0)
	i := 0
	for ; i < limit; i += 4 {
		average0 += float64(payments[i])
		average1 += float64(payments[i+1])
		average2 += float64(payments[i+2])
		average3 += float64(payments[i+3])
	}
	for ; i < count; i++ {
		average0 += float64(payments[i])
	}
	return (average0 + average1 + average2 + average3) / float64(count)
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(payments []float32) float64 {
	mean := AveragePaymentAmount(payments)
	count := len(payments)
	limit := count - 3
	squares0, squares1, squares2, squares3 := float64(0), float64(0), float64(0), float64(0)
	i := 0
	for ; i < limit; i += 4 {
		squares0 += float64(payments[i] * payments[i])
		squares1 += float64(payments[i+1] * payments[i+1])
		squares2 += float64(payments[i+2] * payments[i+2])
		squares3 += float64(payments[i+3] * payments[i+3])
	}
	for ; i < count; i++ {
		squares0 += float64(payments[i] * payments[i])
	}
	return math.Sqrt((squares0+squares1+squares2+squares3)/float64(count) - mean*mean)
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
