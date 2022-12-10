package main

import "fmt"

// SimMemory Contains the simulated memory for the virtual machine state
type SimMemory struct {
	// memFile contains the total virtual memory of the simulated machine
	memFile []uint32

	// maxAddr contains the maximum virtual memory address allocated plus one
	maxAddr int
}

// InitSimMemory returns an initialized SimMemory struct with parsed instruction data
func InitSimMemory() SimMemory {
	s := SimMemory{maxAddr: 96}
	for _, instruction := range InstructionArray {
		s.memFile = append(s.memFile, instruction.lineValue)
		s.maxAddr += 4
	}
	return s
}

// ReadSimMemory returns a single word of data from the specified address
func ReadSimMemory(mem SimMemory, adr int) uint64 {
	ret := uint64(0)
	// alias Holds the sub-word offset in bits
	alias := (adr % 4) * 8
	// idx holds base memory index
	idx := (adr - 96) / 4
	if adr >= 96 && idx < len(mem.memFile) {
		ret = uint64(mem.memFile[idx]) << (32 + alias)
		if idx+1 < len(mem.memFile) {
			ret = ret | uint64(mem.memFile[idx+1])<<alias

			if alias != 0 && idx+2 < len(mem.memFile) {
				ret = ret | uint64(mem.memFile[idx+2])>>(32-alias)
			}
		}
	}
	return ret
}

// WriteSimMemory returns an updated SimMemory memory state
func WriteSimMemory(mem SimMemory, adr int, data uint64) SimMemory {
	if adr >= DataSegmentIndex*4+96 {
		// alias Holds the sub-word offset in bits (0, 8, 16, or 24)
		alias := (adr % 4) * 8
		// idx holds base memory index
		idx := (adr - 96) / 4

		for adr+8 > mem.maxAddr {
			mem.memFile = append(mem.memFile, uint32(0))
			mem.maxAddr += 4
		}

		// Map register value to memory word array
		mem.memFile[idx] = mem.memFile[idx] & (0xFFFFFFFF << (32 - alias))
		mem.memFile[idx] = mem.memFile[idx] | uint32(data>>(32+alias))
		mem.memFile[idx+1] = uint32(data >> alias)
		if alias != 0 {
			mem.memFile[idx+2] = mem.memFile[idx+2] & (0xFFFFFFFF >> alias)
			mem.memFile[idx+2] = mem.memFile[idx+2] | uint32(data<<(32-alias))
		}
	} else {
		fmt.Println("ERROR: Blocked write to instruction memory address [", adr, "]:", data)
	}

	return mem
}
