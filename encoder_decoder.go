package main

import (
	"fmt"
	"strconv"
	"strings"
)

var HashTable = []string{"N", "v", "E", "O", "q", "c", "f", "l", "p", "G"}
var HashSalt = int64(1666)

func EncodeID(userID int64) string {
	userID += HashSalt
	out := fmt.Sprintf("%d", userID)
	for i := range HashTable {
		out = strings.Replace(out, strconv.Itoa(i), HashTable[i], -1)
	}
	return out
}

func DecodeID(str string) int64 {
	for i := range HashTable {
		str = strings.Replace(str, HashTable[i], strconv.Itoa(i), -1)
	}
	out, _ := strconv.ParseInt(str, 10, 32)
	if out > 0 {
		out -= HashSalt
	}
	return out
}

func makeUserHash(sender int64, receiver int64) string {
	return EncodeID((sender / 14) + (receiver / 12) + HashSalt)
}
