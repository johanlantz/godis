package utils

import (
	"fmt"
	"strings"
)

func MarshalToResp(cmd string) []byte {
	segments := strings.Fields(cmd)
	bytes := []byte("")
	bytes = fmt.Append(bytes, "*", len(segments), "\r\n")

	for _, v := range segments {
		bytes = fmt.Append(bytes, "$", len(v), "\r\n", v, "\r\n")
	}
	return bytes
}

func Keys[K comparable, V any](m map[K]V) []K {
	var s []K
	for key := range m {
		s = append(s, key)
	}
	return s
}
