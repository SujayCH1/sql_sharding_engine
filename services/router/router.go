package router

import (
	"hash/crc32"
)

func HashShards() {}

func HashShardKey(key string) uint32 {
	hash := crc32.ChecksumIEEE([]byte(key))

	return hash
}
