package commonTool

import (
	"crypto"
	"encoding/hex"
	"fmt"
	uuid "github.com/satori/go.uuid"
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

func GetNowTime() int64 {
	now := time.Now().UTC()
	return now.Unix()
}

func GetTodayZeroTimeInt() int64 {
	now := time.Now()
	startOfYesterday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return startOfYesterday.Unix()
}

func GetYesterdayZeroTimeInt() int64 {
	now := time.Now().Add(time.Hour * -24)
	startOfYesterday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return startOfYesterday.Unix()
}

func GetMonthZeroTimeInt() int64 {
	now := time.Now()
	startOfYesterday := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	return startOfYesterday.Unix()
}

func GetLastMonthZeroTimeInt() int64 {
	now := time.Now()
	// 减去一个月
	lastMonth := now.AddDate(0, -1, 0)

	startOfYesterday := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	return startOfYesterday.Unix()
}

func GenUUID() uuid.UUID {
	v1 := uuid.NewV1()
	return v1
}

func Md5(str string) string {
	h := crypto.MD5.New()
	_, _ = h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 获取当前时区的0点时间戳,东8区就传8
func GetXZeroTodayTimeInt(X int) int64 {
	// 创建时区偏移量
	loc := time.FixedZone(fmt.Sprintf("UTC+%d", X), X*3600)
	// 获取当前时间
	now := time.Now().In(loc)

	// 获取今天的0点时间
	zeroHour := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	// 获取0点时间的时间戳
	return zeroHour.Unix()
}
