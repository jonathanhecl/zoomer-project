package main

import (
	"os"
	"strconv"
)

func isValidFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func isValidPath(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isValidPort(port string) bool {
	const maxPort = 65535
	const minPort = 1
	n, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	if n < minPort || n > maxPort {
		return false
	}
	return true
}
