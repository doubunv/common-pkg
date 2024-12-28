package commonTool

import (
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	randomString := make([]byte, length)
	for i := 0; i < length; i++ {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}

func TimeToString(timeInt int64) string {
	if timeInt == 0 {
		return ""
	}
	t := time.Unix(timeInt, 0).UTC()
	return t.Format("2006-01-02 15:04:05")
}

func DiffTimeUnix(timeStr1, timeStr2 string) int64 {
	layout := "2006-01-02 15:04:05"
	t1, _ := time.Parse(layout, timeStr1)
	t2, _ := time.Parse(layout, timeStr2)
	return int64(t2.Sub(t1))
}
