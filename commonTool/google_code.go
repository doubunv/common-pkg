package commonTool

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"strings"
	"time"
)

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}

func GenGoogleAuthSecret() (string, error) {
	dictionary := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	for k, v := range bytes {
		bytes[k] = dictionary[int(v)%len(dictionary)]
	}
	return strings.ToUpper(string(bytes)), nil
}

func BuildGoogleQR(secret string, account string, issuer string) string {
	const (
		algorithm = "SHA1"
		digit     = 6
		period    = 30
	)
	// generate QR code url: https://zh.qr-code-generator.com/
	// Refer: https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=%s&digits=%d&period=%d", issuer, account, secret, issuer, algorithm, digit, period)
}

func CheckGoogleAuthCode(code uint32, secret string) error {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return err
	}

	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix()
	pwd := oneTimePassword(key, toBytes(epochSeconds/30))
	if pwd != code {
		return fmt.Errorf("mismatch google auth code")
	}
	return nil
}
