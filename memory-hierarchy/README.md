
# memory-hierarchy

*This exercise was completed as a part of Bradfield's [Computer Science Intensive](https://bradfieldcs.com/courses/languages/)
program.*

## Introduction

memory-hierarchy illustrates the speed ups obtained via cache-friendly software design.

## Context

The performance improvements in optimizing calculations involving average age are as follows:

| calculation                 | ns/op     | improvement |
| --------------------------- | --------- | ----------- |
| baseline                    |   5747016 |             |
| convert to array            |    535163 |         11x |
| remove overflow-safe sum    |     37980 |        151x |
| iteration with index        |     36071 |        159x |
| change from uint64 to uint8 |     30283 |        190x |
| loop unrolling              |     27246 |        211x |
| bound check elimination     |     18240 |        315x |

The performance improvements in optimizing calculations involving average payments are as follows:

| calculation                 | ns/op    | improvement |
|-----------------------------|----------|-------------|
| baseline                    | 29136655 |             |
| convert to array            |  5398410 |          5x |
| remove overflow-safe sum    |  1093093 |         27x |
| change from uint64 to uint8 |  1078355 |         27x |
| loop unrolling              |   511733 |         57x |
| bound check elimination     |   345569 |         84x |

The performance improvements in optimizing calculations involving average payments are as follows:


| calculation                 | ns/op    | improvement |
|-----------------------------|----------|-------------|
| baseline                    | 56252055 |             |
| convert to array            |  6379861 |          9x |
| remove overflow-safe sum    |  2377747 |         24x |
| change from uint64 to uint8 |  2029108 |         28x |
| loop unrolling              |  1313912 |         43x |
| bound check elimination     |   723575 |         78x |

## Baseline

We start with the following program (also viewable [here](https://github.com/savarin/computer-systems/blob/5bf88d681aa56e0f8650c00756ae18dd68a9dd0f/memory-hierarchy/metrics.go)).

```go

package metrics

import (
    "encoding/csv"
    "log"
    // "math"
    "os"
    "strconv"
    "time"
)

type UserId int
type UserMap map[UserId]*User

type Address struct {
    fullAddress string
    zip         int
}

type DollarAmount struct {
    dollars, cents uint64
}

type Payment struct {
    amount DollarAmount
    time   time.Time
}

type User struct {
    id       UserId
    name     string
    age      int
    address  Address
    payments []Payment
}

func AverageAge(users UserMap) float64 {
    average, count := 0.0, 0.0
    for _, u := range users {
        count += 1
        average += (float64(u.age) - average) / count
    }
    return average
}

func LoadData() UserMap {
    f, err := os.Open("users.csv")
    if err != nil {
        log.Fatalln("Unable to read users.csv", err)
    }
    reader := csv.NewReader(f)
    userLines, err := reader.ReadAll()
    if err != nil {
        log.Fatalln("Unable to parse users.csv as csv", err)
    }

    users := make(UserMap, len(userLines))
    for _, line := range userLines {
        id, _ := strconv.Atoi(line[0])
        name := line[1]
        age, _ := strconv.Atoi(line[2])
        address := line[3]
        zip, _ := strconv.Atoi(line[3])
        users[UserId(id)] = &User{UserId(id), name, age, Address{address, zip}, []Payment{}}
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

    for _, line := range paymentLines {
        userId, _ := strconv.Atoi(line[2])
        paymentCents, _ := strconv.Atoi(line[0])
        datetime, _ := time.Parse(time.RFC3339, line[1])
        users[UserId(userId)].payments = append(users[UserId(userId)].payments, Payment{
            DollarAmount{uint64(paymentCents / 100), uint64(paymentCents % 100)},
            datetime,
        })
    }

    return users
}

```

## Improvements

**Step 1** The first change involves replacing the array of structs (representing each user) with an array
of ints (representing the age of each user). This allows the cache line to be 'packed' with the
values needed for the calculation of the mean, instead of user attributes that are not relevant.

For ease of review, we focus only on the age calculations. The changes described can be reviewed
in a PR format [here](https://github.com/savarin/computer-systems/commit/d25677e18a31621262e6c151a56079964f45d6cc?branch=d25677e18a31621262e6c151a56079964f45d6cc&diff=split).

| calculation                 | ns/op     | improvement |
|-----------------------------|-----------|-------------|
| baseline                    |   5747016 |             |
| convert to array            |    535163 |         11x |

```go
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

```

**Step 2** Next we remove the in-loop calculation that makes the sum overflow-safe. In particular, the division
operation is much more expensive compared to addition (PR format [here](https://github.com/savarin/computer-systems/commit/136c2af1bf958c12f26544047f4c8ac5c9877ef3?branch=136c2af1bf958c12f26544047f4c8ac5c9877ef3&diff=split)).

| calculation                 | ns/op     | improvement |
|-----------------------------|-----------|-------------|
| baseline                    |   5747016 |             |
| convert to array            |    535163 |         11x |
| remove overflow-safe sum    |     37980 |        151x |

```go
func AverageAge(ages []int) float64 {
    count := len(ages)
    average := 0
    for _, age := range ages {
        average += age
    }
    return float64(average) / float64(count)
}

```

**Step 3** Iterating through the array with the index speeds things up as iterating through the values requires
a copy of the array (PR format [here](https://github.com/savarin/computer-systems/commit/78dea3280bbccbd3da61c3c7f1baf1d02a89e1f3?branch=78dea3280bbccbd3da61c3c7f1baf1d02a89e1f3&diff=split)).

| calculation                 | ns/op     | improvement |
|-----------------------------|-----------|-------------|
| baseline                    |   5747016 |             |
| convert to array            |    535163 |         11x |
| remove overflow-safe sum    |     37980 |        151x |
| iteration with index        |     36071 |        159x |

```go
func AverageAge(ages []int) float64 {
    count := len(ages)
    average := 0
    for i := range ages {
        average += ages[i]
    }
    return float64(average) / float64(count)
}
```

**Step 4** We can further pack the cache line by using `uint8` instead of `uint64` (PR format [here](https://github.com/savarin/computer-systems/commit/bbff04d1aaf6e545170384dedbf37ac6a4529872?branch=bbff04d1aaf6e545170384dedbf37ac6a4529872&diff=split)).

| calculation                 | ns/op     | improvement |
|-----------------------------|-----------|-------------|
| baseline                    |   5747016 |             |
| convert to array            |    535163 |         11x |
| remove overflow-safe sum    |     37980 |        151x |
| iteration with index        |     36071 |        159x |
| change from uint64 to uint8 |     30283 |        190x |

```go
func AverageAge(ages []uint8) float64 {
    count := len(ages)
    average := uint64(0)
    for i := range ages {
        average += uint64(ages[i])
    }
    return float64(average) / float64(count)
}
```

**Step 5** We introduce loop unrolling, which allows for more efficient pipelining (PR format [here](https://github.com/savarin/computer-systems/commit/b181e048307c5b6a947bdde545ab8e8d11ccb4ff?branch=b181e048307c5b6a947bdde545ab8e8d11ccb4ff&diff=split)).

| calculation                 | ns/op     | improvement |
|-----------------------------|-----------|-------------|
| baseline                    |   5747016 |             |
| convert to array            |    535163 |         11x |
| remove overflow-safe sum    |     37980 |        151x |
| iteration with index        |     36071 |        159x |
| change from uint64 to uint8 |     30283 |        190x |
| loop unrolling              |     27246 |        211x |

```go
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
```

**Step 6** Finally, we remove array bound checks by running `go test -gcflags="-B" -bench=.`.

| calculation                 | ns/op     | improvement |
|-----------------------------|-----------|-------------|
| baseline                    |   5747016 |             |
| convert to array            |    535163 |         11x |
| remove overflow-safe sum    |     37980 |        151x |
| iteration with index        |     36071 |        159x |
| change from uint64 to uint8 |     30283 |        190x |
| loop unrolling              |     27246 |        211x |
| bound check elimination     |     18240 |        315x |
