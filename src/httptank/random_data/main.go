package random_data

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

var data = make([]string, 0)

func Init(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	return scanner.Err()
}

func GetRandArg() string {
	if len(data) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return data[rand.Intn(len(data)-1)]
}
