package id

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"hash/fnv"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/sony/sonyflake"
)

var flake *sonyflake.Sonyflake
var fallbackSeq uint64

// IssueID 生成唯一 ID (基于 Sony 改进的 Snowflake 算法)
// https://github.com/sony/sonyflake
func IssueID() string {
	if flake == nil {
		return fallbackID()
	}

	next, err := flake.NextID()
	if err != nil {
		return fallbackID()
	}
	return strconv.FormatUint(next, 10)
}

// GenerateID 生成唯一 ID (基于 Sonyflake + MD5)
func GenerateID() string {
	keyID := IssueID()
	id := fmt.Sprintf("%x", md5.Sum([]byte(keyID)))
	return id
}

func init() {
	flake = newSonyflake()
}

// Sonyflake 默认会从本机私有 IPv4 生成 machine id；
// 在仅有 loopback/仅 IPv6 等环境下可能返回 nil，进而触发空指针 panic。
// 这里显式提供 MachineID 计算逻辑，保证初始化与生成 ID 的稳定性。
func newSonyflake() *sonyflake.Sonyflake {
	sf, err := sonyflake.New(sonyflake.Settings{
		MachineID: func() (uint16, error) {
			return resolveMachineID(), nil
		},
	})
	if err == nil && sf != nil {
		return sf
	}

	// 理论上不会走到这里（MachineID 不返回 error），但保留一个不会失败的兜底。
	sf, _ = sonyflake.New(sonyflake.Settings{
		MachineID: func() (uint16, error) {
			return 0, nil
		},
	})
	return sf
}

func resolveMachineID() uint16 {
	if id, ok := machineIDFromIPv4(); ok {
		return id
	}
	if id, ok := machineIDFromMAC(); ok {
		return id
	}
	if id, ok := machineIDFromHostname(); ok {
		return id
	}
	return randomUint16()
}

func machineIDFromIPv4() (uint16, bool) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0, false
	}

	var firstNonLoopback net.IP
	var firstLoopback net.IP

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP == nil {
			continue
		}

		ip4 := ipnet.IP.To4()
		if ip4 == nil {
			continue
		}

		if ip4.IsLoopback() {
			if firstLoopback == nil {
				firstLoopback = ip4
			}
			continue
		}

		if isPrivateIPv4(ip4) {
			return lower16(ip4), true
		}
		if firstNonLoopback == nil {
			firstNonLoopback = ip4
		}
	}

	if firstNonLoopback != nil {
		return lower16(firstNonLoopback), true
	}
	if firstLoopback != nil {
		return lower16(firstLoopback), true
	}
	return 0, false
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 ||
			ip[0] == 172 && ip[1] >= 16 && ip[1] < 32 ||
			ip[0] == 192 && ip[1] == 168 ||
			ip[0] == 169 && ip[1] == 254)
}

func lower16(ip net.IP) uint16 {
	return uint16(ip[2])<<8 | uint16(ip[3])
}

func machineIDFromMAC() (uint16, bool) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return 0, false
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		hw := iface.HardwareAddr
		if len(hw) < 2 {
			continue
		}
		return uint16(hw[len(hw)-2])<<8 | uint16(hw[len(hw)-1]), true
	}
	return 0, false
}

func machineIDFromHostname() (uint16, bool) {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		return 0, false
	}

	h := fnv.New32a()
	_, _ = h.Write([]byte(hostname))
	return uint16(h.Sum32() & 0xFFFF), true
}

func randomUint16() uint16 {
	var b [2]byte
	if _, err := rand.Read(b[:]); err == nil {
		return uint16(b[0])<<8 | uint16(b[1])
	}
	return uint16(uint64(time.Now().UnixNano()) & 0xFFFF)
}

func fallbackID() string {
	nowMs := uint64(time.Now().UnixMilli())
	seq := atomic.AddUint64(&fallbackSeq, 1) & 0xFFFF
	next := (nowMs << 16) | seq
	return strconv.FormatUint(next, 10)
}
