package common

import (
	"fmt"
	"hash/fnv"
	"os"
	"sync"

	"github.com/sony/sonyflake"
)

var (
	sf   *sonyflake.Sonyflake
	once sync.Once
)

// ambil machine ID dari hostname hash (unik per pod)
func getMachineID() (uint16, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return 0, err
	}
	hash := fnv.New32()
	hash.Write([]byte(hostname))
	return uint16(hash.Sum32() % 65535), nil
}

// panggil ini saat startup app
func InitSonyflake() {
	once.Do(func() {
		sf = sonyflake.NewSonyflake(sonyflake.Settings{
			MachineID: getMachineID,
		})
		if sf == nil {
			panic("failed to initialize Sonyflake")
		}
	})
}

func NextID() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return id
}

// Generate ID with optional prefix (e.g. "PAN", "INV")
func GenerateID(prefix string) string {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	if prefix == "" {
		return fmt.Sprintf("%d", id)
	}
	return fmt.Sprintf("%s-%d", prefix, id)
}

// Base62 encode (uint64 to string)
func encodeBase62(n uint64) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(chars[n%62]) + result
		n /= 62
	}
	return result
}

// Generate short ID like "INV-aZ2pK9fX1b"
func GenerateShortID(prefix string) string {
	id := NextID()
	return fmt.Sprintf("%s-%s", prefix, encodeBase62(id))
}
