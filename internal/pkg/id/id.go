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

func IssueMd5ID() string {
	keyID := IssueID()
	id := fmt.Sprintf("%x", md5.Sum([]byte(keyID)))
	return id
}

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}
