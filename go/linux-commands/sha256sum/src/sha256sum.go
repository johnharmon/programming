package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func HashFile(f *os.File, chunkSize int) ([]byte, error) {
	hash := sha256.New()
	buf := make([]byte, chunkSize)

	for {
		bytesRead, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return buf, err
		}
		if bytesRead > 0 {
			hash.Write(buf[:bytesRead])
		} else {
			break
		}
	}
	hash_bytes := hash.Sum(nil)
	return hash_bytes, nil
}

func HashToString(h []byte) string {
	return hex.EncodeToString(h)
}

func HashFileToString(filepath string, chunkSize int) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close() // Close the file when the function returns
	hash, err := HashFile(file, chunkSize)
	if err != nil {
		return "", err
	}
	return HashToString(hash), nil
}

func main() {

	args := os.Args
	if len(args) < 1 {
		fmt.Println("Usage: go run hash.go <file>")
		os.Exit(1)
	}
	filepath := os.Args[0]

	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(2)
	} else {
		defer file.Close()
		hash, err := HashFile(file, 1024)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(3)
		}
		fmt.Printf("%x\n", hash)
	}
}
