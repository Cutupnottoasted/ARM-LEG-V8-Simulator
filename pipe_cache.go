package main

import "fmt"

type Block struct {
	valid int
	dirty int
	tag   int
	word1 int
	word2 int
	LRU   int
}

var CacheSets [4][2]Block
var JustMissedList = []int{-1, -1, -1, -1, -1, -1, -1}

// CacheRead Returns a hit or miss with corresponding data
// If hit, return true and memory value
// If missed, return false and update cache
func CacheRead(addr int) (bool, uint64) {
	hit1 := false
	hit2 := false
	tag := addr >> 5
	set := (addr >> 3) & 0x3
	idx := addr & 0x7
	memVal := ReadSimMemory(MainMemory, addr)

	if addr == 0 {
		return true, 0
	}

	if CacheSets[set][0].valid == 1 && CacheSets[set][0].tag == tag {
		hit1 = true
		if idx == 0 {
			memVal = uint64(CacheSets[set][0].word1)<<32 | uint64(CacheSets[set][0].word2)
			hit2 = true
		} else {

			// Evaluate second cache-line in call
			addr2 := addr + 8
			tag := addr2 >> 5
			set := (addr2 >> 3) & 0x3
			if CacheSets[set][0].valid == 1 && CacheSets[set][0].tag == tag {
				//TODO Aliased read
				hit2 = true
				CacheSets[set][0].LRU = 1
				CacheSets[set][1].LRU = 1
			} else if CacheSets[set][1].valid == 1 && CacheSets[set][1].tag == tag {
				hit2 = true
				CacheSets[set][0].LRU = 0
				CacheSets[set][1].LRU = 0
			} // else { hit2 = false }
		}
	} else if CacheSets[set][1].valid == 1 && CacheSets[set][1].tag == tag {
		hit1 = true
		if idx == 0 {
			memVal = uint64(CacheSets[set][1].word1)<<32 | uint64(CacheSets[set][1].word2)
			hit2 = true
		} else {
			// Evaluate second cache-line in call
			addr2 := addr + 8
			tag := addr2 >> 5
			set := (addr2 >> 3) & 0x3
			if CacheSets[set][0].valid == 1 && CacheSets[set][0].tag == tag {
				hit2 = true
				CacheSets[set][0].LRU = 1
				CacheSets[set][1].LRU = 1
			} else if CacheSets[set][1].valid == 1 && CacheSets[set][1].tag == tag {
				hit2 = true
				CacheSets[set][0].LRU = 0
				CacheSets[set][1].LRU = 0
			} // else { hit = false }
		}
	}
	if hit1 == false {
		for i := 0; i < 4; i++ {
			if JustMissedList[i] == addr^idx {
				break
			} else if JustMissedList[i] == -1 {
				JustMissedList[i] = addr ^ idx
				break
			}
		}
	}
	if idx != 0 && hit2 == false {
		for i := 0; i < 4; i++ {
			if JustMissedList[i] == (addr^idx)+8 {
				break
			} else if JustMissedList[i] == -1 {
				JustMissedList[i] = (addr ^ idx) + 8
				break
			}
		}
	}
	return hit1 && hit2, memVal
}

// CacheUpdatePostCycle updates cache state based on previous cycle misses
func CacheUpdatePostCycle() {
	for i, addr := range JustMissedList {
		if addr == -1 {
			break
		}
		tag := addr >> 5
		set := (addr >> 3) & 0x3

		lru := CacheSets[set][0].LRU

		JustMissedList[i] = -1
		if CacheSets[set][lru].dirty == 1 {
			MainMemory = WriteSimMemory(MainMemory, addr, uint64(CacheSets[set][lru].word1)<<32|uint64(CacheSets[set][lru].word2))
			/*if JustMissedList[i] == addr {
				break
			} else if JustMissedList[i] == -1 {
				JustMissedList[i] = addr
				break
			}*/
		}
		memVal := ReadSimMemory(MainMemory, addr)
		CacheSets[set][lru].word1 = int(memVal >> 32)
		CacheSets[set][lru].word2 = int(memVal)
		CacheSets[set][lru].tag = tag
		CacheSets[set][lru].valid = 1

		if lru == 0 {
			CacheSets[set][0].LRU = 1
			CacheSets[set][1].LRU = 1
		} else {
			CacheSets[set][0].LRU = 0
			CacheSets[set][1].LRU = 0
		}
	}
}

// CacheWrite Returns a hit or miss and writes data to memory
// If hit, return true and update cache
// If missed, return false and evict a cache-line
func CacheWrite(addr int, value uint64) bool {
	hit1 := false
	//hit2 := false
	tag := addr >> 5
	set := (addr >> 3) & 0x3
	idx := addr & 0x7
	memVal := ReadSimMemory(MainMemory, addr^idx)
	lru := CacheSets[set][0].LRU
	writeVal := int(value)

	if addr < DataSegmentAddress {
		fmt.Println("Blocked Illegal write to instruction segment")
		return true
	}

	if CacheSets[set][0].valid == 1 && CacheSets[set][0].tag == tag {
		hit1 = true
		CacheSets[set][0].tag = tag
		CacheSets[set][0].word1 = int(memVal >> 32)
		CacheSets[set][0].word2 = int(memVal)
		CacheSets[set][0].dirty = 1

		if addr%8 == 0 {
			CacheSets[set][0].word1 = writeVal
		} else {
			CacheSets[set][0].word2 = writeVal
		}

		CacheSets[set][0].LRU = 1
		CacheSets[set][1].LRU = 1
	} else if CacheSets[set][1].valid == 1 && CacheSets[set][1].tag == tag {
		hit1 = true
		CacheSets[set][1].tag = tag
		CacheSets[set][1].word1 = int(memVal >> 32)
		CacheSets[set][1].word2 = int(memVal)
		CacheSets[set][1].dirty = 1

		if addr%8 == 0 {
			CacheSets[set][1].word1 = writeVal
		} else {
			CacheSets[set][1].word2 = writeVal
		}

		CacheSets[set][0].LRU = 0
		CacheSets[set][1].LRU = 0
	} else if CacheSets[set][lru].valid == 0 {
		hit1 = true
		CacheSets[set][lru].tag = tag
		CacheSets[set][lru].word1 = int(memVal >> 32)
		CacheSets[set][lru].word2 = int(memVal)
		CacheSets[set][lru].dirty = 1

		if addr%8 == 0 {
			CacheSets[set][lru].word1 = writeVal
		} else {
			CacheSets[set][lru].word2 = writeVal
		}

		CacheSets[set][lru].valid = 1
		CacheSets[set][0].LRU = lru ^ 1
		CacheSets[set][1].LRU = lru ^ 1
	} else {
		for i := 0; i < 4; i++ {
			if JustMissedList[i] == addr^idx {
				break
			} else if JustMissedList[i] == -1 {
				JustMissedList[i] = addr ^ idx
				break
			}
		}
	}

	if hit1 {
		for addr+4 > MainMemory.maxAddr {
			MainMemory.memFile = append(MainMemory.memFile, uint32(0))
			MainMemory.maxAddr += 4
		}
		MainMemory.memFile[(addr-96)/4] = uint32(writeVal)
	}
	return hit1
}

// FlushCache Writes all valid, dirty data into MainMemory
// Call on BREAK after all execution is finished
func FlushCache() {
	for set, v := range CacheSets {
		for _, w := range v {
			if w.dirty == 1 {
				MainMemory = WriteSimMemory(MainMemory, w.tag<<5|set<<3, uint64(w.word1)<<32|uint64(w.word2))
			}
		}
	}
}

func isValidInCache(addr int) {

}
