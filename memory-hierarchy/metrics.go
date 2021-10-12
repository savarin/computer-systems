package metrics

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func AverageAge(ages []int) float64 {
	count := len(ages)
	average := 0
	for i := range ages {
		average += ages[i]
	}
	return float64(average) / float64(count)
}

func LoadData() []int {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	users := make([]int, len(userLines))
	for i, line := range userLines {
		age, _ := strconv.Atoi(line[2])
		users[i] = age
	}

	return users
}
