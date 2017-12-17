package main

import (
	"strings"
	"fmt"
	"crypto/rand"
)

func pseudoUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return
	}
	uuid = strings.ToLower(
		fmt.Sprintf("%X-%X-%X-%X-%X",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))
	return
}
