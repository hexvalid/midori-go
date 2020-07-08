package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"regexp"
	"strings"
)

func FindInString(str, start, end string) (string, error) {
	var match []byte
	index := strings.Index(str, start)
	if index == -1 {
		return "", errors.New("string is not found in target string")
	}
	index += len(start)
	for {
		char := str[index]
		if strings.HasPrefix(str[index:index+len(match)], end) {
			break
		}
		match = append(match, char)
		index++
	}
	return string(match), nil
}

func RegexString(str, start, end string) [][]string {
	re := regexp.MustCompile(start + "(.*?)" + end)
	rm := re.FindAllStringSubmatch(str, -1)
	return rm
}

func SHA256(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func MD5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CRC32(f float64) string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, f)
	hasher := crc32.NewIEEE()
	hasher.Write([]byte(buf.Bytes()))
	return hex.EncodeToString(hasher.Sum(nil))
}
