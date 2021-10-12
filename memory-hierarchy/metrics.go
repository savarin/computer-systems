package metrics

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func AverageAge(ages []int) float64 {
	average, count := 0.0, 0.0
	for _, age := range ages {
		count += 1
		average += (float64(age) - average) / count
	}
	return average
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
