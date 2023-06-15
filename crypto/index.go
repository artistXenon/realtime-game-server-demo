package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"hash/crc32"
	"sort"
	"strings"
)

func GenerateHash(values *[]string) string {
	sort.Strings(*values)
	plain := strings.Join(*values, "")

	hash := sha256.New()
	hash.Write([]byte(plain))
	md := hash.Sum(nil)
	finalHash := base64.StdEncoding.EncodeToString(md)

	return finalHash
}

func GenerateCRCHash(values []byte) uint32 {
	table := crc32.MakeTable(crc32.Castagnoli)
	return crc32.Checksum(values, table)
}
