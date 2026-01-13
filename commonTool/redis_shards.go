package commonTool

import (
	"fmt"
	"github.com/cespare/xxhash/v2"
	"strconv"
)

func ShardIndexXX(userID int64, n uint64) uint64 {
	// xxhash 直接 hash 字符串 bytes
	b := []byte(strconv.FormatInt(userID, 10))
	return xxhash.Sum64(b) % n
}

func ShardTagXX(userID int64, n uint64) string {
	idx := ShardIndexXX(userID, n)
	return fmt.Sprintf("{s%03d}", idx)
}
