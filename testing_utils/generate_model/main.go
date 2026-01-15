package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/exp/rand"
)

func main() {
	n, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println("Error parsing number of values:", err)
		return
	}

	rand.Seed(uint64(time.Now().UnixNano()))

	// Ensure /data directory exists
	dataDir := "../../data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Println("Error creating /data directory:", err)
		return
	}

	filePath := filepath.Join(dataDir, "data_"+strconv.FormatUint(n, 10)+".bin")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	buf := make([]byte, 8)

	for i := uint64(0); i < n; i++ {
		v := rand.Int63()
		binary.LittleEndian.PutUint64(buf, uint64(v))
		_, err := file.Write(buf)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}
