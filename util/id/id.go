package id

import (
	"crypto/md5"
	"fmt"

	"github.com/sony/sonyflake"
)

var flake *sonyflake.Sonyflake

func initSonyFlake() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

// IssueID Unique ID generated using Sony's improved twite snowflake algorithm
// https://github.com/sony/sonyflake
func IssueID() string {
	if flake == nil {
		initSonyFlake()
	}
	id, _ := flake.NextID()
	return fmt.Sprintf("%v", id)
}

// IssueMd5ID Deprecated: 请使用 GenerateID
func IssueMd5ID() string {
	return GenerateID()
}

// GenerateID 使用 Sonyflake + MD5 生成唯一 ID，用于请求追踪
func GenerateID() string {
	keyID := IssueID()
	return fmt.Sprintf("%x", md5.Sum([]byte(keyID)))
}

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}
