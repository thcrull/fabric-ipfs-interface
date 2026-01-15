package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"golang.org/x/exp/rand"
)

func main() {
	n, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println("Error parsing number of files: ", err)
		return
	}

	rand.Seed(uint64(time.Now().UnixNano()))

	file, _ := os.Create("data_" + strconv.FormatUint(n, 10) + ".bin")
	defer file.Close()

	buf := make([]byte, 8)

	for i := uint64(0); i < n; i++ {
		v := rand.Float64()
		binary.LittleEndian.PutUint64(buf, math.Float64bits(v))
		file.Write(buf)
	}

	return
}
