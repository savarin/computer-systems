package metrics

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func AverageAge(ages []uint8) float64 {
	count := len(ages)
	average := uint64(0)
	for i := range ages {
		average += uint64(ages[i])
	}
	return float64(average) / float64(count)
}

func LoadData() []uint8 {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	users := make([]uint8, len(userLines))
	for i, line := range userLines {
		age, _ := strconv.Atoi(line[2])
		users[i] = uint8(age)
	}

	return users
}
