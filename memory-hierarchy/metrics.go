package metrics

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type custom uint64

func AverageAge(ages []custom) float64 {
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

func LoadData() []custom {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	users := make([]custom, len(userLines))
	for i, line := range userLines {
		age, _ := strconv.Atoi(line[2])
		users[i] = custom(age)
	}

	return users
}
