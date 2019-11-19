// 哈希算法有BKDRHash，APHash，DJBHash，JSHash，RSHash，SDBMHash，PJWHash，ELFHash等等,这些都是比较经典的
package hash

// BKDR hash是一种字符哈希算法
func BKDRHash(b []byte) uint32 {
	var seed uint32 = 131 // 31 131 1313 13131 131313 etc...
	var hash uint32 = 0
	for i := 0; i < len(b); i++ {
		hash = hash * seed + uint32(b[i])
	}
	return hash
}

func BKDRHash64(b []byte) uint64 {
	var seed uint64 = 131 // 31 131 1313 13131 131313 etc...
	var hash uint64 = 0
	for i := 0; i < len(b); i++ {
		hash = hash * seed + uint64(b[i])
	}
	return hash
}

// SDBM Hash
func SDBMHash(b []byte) uint32 {
	var hash uint32 = 0
	for i := 0; i < len(b); i++ {
		// equivalent to: hash = 65599*hash + uint32(b[i]);
		hash = uint32(b[i]) + (hash << 6) + (hash << 16) - hash
	}
	return hash
}

func SDBMHash64(b []byte) uint64 {
	var hash uint64 = 0
	for i := 0; i < len(b); i++ {
		// equivalent to: hash = 65599*hash + uint32(b[i])
		hash = uint64(b[i]) + (hash << 6) + (hash << 16) - hash
	}
	return hash
}

func RSHash(b []byte) uint32 {
	var c uint32 = 378551
	var a uint32 = 63689
	var hash uint32 = 0
	for i := 0; i < len(b); i++ {
		hash = hash*a + uint32(b[i])
		a *= c
	}
	return hash
}

// RS Hash
func RSHash64(b []byte) uint64 {
	var c uint64 = 378551
	var a uint64 = 63689
	var hash uint64 = 0
	for i := 0; i < len(b); i++ {
		hash = hash*a + uint64(b[i])
		a *= c
	}
	return hash
}

// JS Hash
func JSHash(b []byte) uint32 {
	var hash uint32 = 1315423911
	for i := 0; i < len(b); i++ {
		hash ^= (hash << 5) + uint32(b[i]) + (hash >> 2)
	}
	return hash
}

// JS Hash
func JSHash64(b []byte) uint64 {
	var hash uint64 = 1315423911
	for i := 0; i < len(b); i++ {
		hash ^= (hash << 5) + uint64(b[i]) + (hash >> 2)
	}
	return hash
}

// P. J. Weinberger Hash
func PJWHash(b []byte) uint32 {
	var BitsInUnignedInt uint32 = 4 * 8
	var ThreeQuarters uint32 = (BitsInUnignedInt * 3) / 4
	var OneEighth uint32 = BitsInUnignedInt / 8
	var HighBits uint32 = (0xFFFFFFFF) << (BitsInUnignedInt - OneEighth)
	var hash uint32 = 0
	var test uint32 = 0
	for i := 0; i < len(b); i++ {
		hash = (hash << OneEighth) + uint32(b[i])
		if test = hash & HighBits; test != 0 {
			hash = (hash ^ (test >> ThreeQuarters)) & (^HighBits + 1)
		}
	}
	return hash
}

// P. J. Weinberger Hash
func PJWHash64(b []byte) uint64 {
	var BitsInUnignedInt uint64 = 4 * 8
	var ThreeQuarters uint64 = (BitsInUnignedInt * 3) / 4
	var OneEighth uint64 = BitsInUnignedInt / 8
	var HighBits uint64 = (0xFFFFFFFFFFFFFFFF) << (BitsInUnignedInt - OneEighth)
	var hash uint64 = 0
	var test uint64 = 0
	for i := 0; i < len(b); i++ {
		hash = (hash << OneEighth) + uint64(b[i])
		if test = hash & HighBits; test != 0 {
			hash = (hash ^ (test >> ThreeQuarters)) & (^HighBits + 1)
		}
	}
	return hash
}

func ELFHash(b []byte) uint32 {
	var hash uint32 = 0
	var x uint32 = 0
	for i := 0; i < len(b); i++ {
		hash = (hash << 4) + uint32(b[i])
		if x = hash & 0xF0000000; x != 0 {
			hash ^= x >> 24
			hash &= ^x + 1
		}
	}
	return hash
}

func ELFHash64(b []byte) uint64 {
	var hash uint64 = 0
	var x uint64 = 0
	for i := 0; i < len(b); i++ {
		hash = (hash << 4) + uint64(b[i])
		if x = hash & 0xF000000000000000; x != 0 {
			hash ^= x >> 24
			hash &= ^x + 1
		}
	}
	return hash
}

func DJBHash(b []byte) uint32 {
	var hash uint32 = 5381
	for i := 0; i < len(b); i++ {
		hash += (hash << 5) + uint32(b[i])
	}
	return hash
}

func DJBHash64(b []byte) uint64 {
	var hash uint64 = 5381
	for i := 0; i < len(b); i++ {
		hash += (hash << 5) + uint64(b[i])
	}
	return hash
}

func APHash(b []byte) uint32 {
	var hash uint32 = 0
	for i := 0; i < len(b); i++ {
		if (i & 1) == 0 {
			hash ^= (hash << 7) ^ uint32(b[i]) ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ uint32(b[i]) ^ (hash >> 5)) + 1
		}
	}
	return hash
}

func APHash64(b []byte) uint64 {
	var hash uint64 = 0
	for i := 0; i < len(b); i++ {
		if (i & 1) == 0 {
			hash ^= (hash << 7) ^ uint64(b[i]) ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ uint64(b[i]) ^ (hash >> 5)) + 1
		}
	}
	return hash
}





