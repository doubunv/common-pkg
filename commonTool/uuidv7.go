package commonTool

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

// 返回36位长度字符串
func NewUUidV7() (string, error) {
	var u [16]byte

	// 1. 写入 48-bit 毫秒时间戳（big endian）
	ts := uint64(time.Now().UnixMilli())
	u[0] = byte(ts >> 40)
	u[1] = byte(ts >> 32)
	u[2] = byte(ts >> 24)
	u[3] = byte(ts >> 16)
	u[4] = byte(ts >> 8)
	u[5] = byte(ts)

	// 2. 随机填充剩余字节
	if _, err := rand.Read(u[6:]); err != nil {
		return "", err
	}

	// 3. 设置版本号（v7）
	u[6] = (u[6] & 0x0f) | 0x70

	// 4. 设置 variant（RFC 4122）
	u[8] = (u[8] & 0x3f) | 0x80

	// 5. 转字符串
	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		binary.BigEndian.Uint32(u[0:4]),
		binary.BigEndian.Uint16(u[4:6]),
		binary.BigEndian.Uint16(u[6:8]),
		binary.BigEndian.Uint16(u[8:10]),
		u[10:16],
	), nil
}
