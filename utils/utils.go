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
