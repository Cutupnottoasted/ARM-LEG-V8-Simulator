package main

import (
	"fmt"
	"os"
	"strconv"
)

var RegisterFile [32]uint64
var MainMemory SimMemory
var ProgramCounter = 96

func SimulatePipeline(output string) {
	file, err := os.Create(output + "_pipeline.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	MainMemory = InitSimMemory()
	pipeBreak = false
	cycle := 0
	// TODO Remove cycle limit after testing is done!
	for cycle = 0; (pipeBreak == false || pipelineIsEmpty() == false) && cycle < CycleLimit; cycle++ {
		//for cycle := 0; cycle < 4; cycle++ {
		SimWriteBack()
		SimALU()
		SimMem()
		SimIssue()
		SimFetch()
		if pipeBreak && pipelineIsEmpty() {
			FlushCache()
		}
		fmt.Fprintln(file, PipelineStateToString(cycle))
		CacheUpdatePostCycle()
	}
}

func pipelineIsEmpty() bool {
	preIssueIsEmpty := PreIssueBuffer[0] == -1
	postIssueIsEmpty := PreALUBuffer[0] == -1 && PreMemBuffer[0] == -1
	preWBIsEmpty := PostALUBuffer == -1 && PostMemBuffer == -1
	return preIssueIsEmpty && postIssueIsEmpty && preWBIsEmpty
}

func PipelineStateToString(cycle int) string {
	strOut := "--------------------\n"
	strOut += "Cycle:" + strconv.Itoa(cycle) + "\n\n"

	strOut += "Pre-Issue Buffer:\n"
	for i, v := range PreIssueBuffer {
		if v != -1 {
			strOut += "\tEntry " + strconv.Itoa(i) + ":\t[" + InstructionArray[v].decodedStr + "]\n"
		} else {
			strOut += "\tEntry " + strconv.Itoa(i) + ":\n"
		}
	}
	strOut += "Pre_ALU Queue:\n"
	for i, v := range PreALUBuffer {
		if v != -1 {
			strOut += "\tEntry " + strconv.Itoa(i) + ":\t[" + InstructionArray[v].decodedStr + "]\n"
		} else {
			strOut += "\tEntry " + strconv.Itoa(i) + ":\n"
		}
	}
	strOut += "Post_ALU Queue:\n"
	if PostALUBuffer != -1 {
		strOut += "\tEntry 0:\t[" + InstructionArray[PostALUBuffer].decodedStr + "]\n"
	} else {
		strOut += "\tEntry 0:\n"
	}
	strOut += "Pre_MEM Queue:\n"
	for i, v := range PreMemBuffer {
		if v != -1 {
			strOut += "\tEntry " + strconv.Itoa(i) + ":\t[" + InstructionArray[v].decodedStr + "]\n"
		} else {
			strOut += "\tEntry " + strconv.Itoa(i) + ":\n"
		}
	}
	strOut += "Post_MEM Queue:\n"
	if PostMemBuffer != -1 {
		strOut += "\tEntry 0:\t[" + InstructionArray[PostMemBuffer].decodedStr + "]\n"
	} else {
		strOut += "\tEntry 0:\n"
	}

	strOut += "\nRegisters\n"
	for i, reg := range RegisterFile {
		if i == 0 || i == 8 {
			strOut += "r0" + strconv.Itoa(i) + ":"
		} else if i == 16 || i == 24 {
			strOut += "r" + strconv.Itoa(i) + ":"
		}

		strOut += "\t" + strconv.FormatInt(int64(reg), 10)

		if i%8 == 7 {
			strOut += "\n"
		}
	}

	strOut += "\nCache\n"
	for i, v := range CacheSets {
		strOut += "Set " + strconv.Itoa(i) + ": LRU=" + strconv.Itoa(CacheSets[i][0].LRU) + "\n"
		for j, q := range v {
			strOut += "\tEntry " + strconv.Itoa(j) + ": ["
			strOut += "(" + strconv.Itoa(q.valid) + "," + strconv.Itoa(q.dirty) + "," + strconv.Itoa(q.tag) + ")"
			strOut += "<" + printCacheWord(q.word1) + "," + printCacheWord(q.word2) + ">]\n"
		}

	}

	strOut += "\nData"

	for i := DataSegmentIndex; i < len(MainMemory.memFile); i++ {
		if i%8 == DataSegmentIndex%8 {
			strOut += "\n" + strconv.Itoa(i*4+96) + ":"
		}
		strOut += "\t" + strconv.Itoa(int(MainMemory.memFile[i]))
	}

	strOut += "\n#########################  END OF CYCLE #########################"
	return strOut
}

// printCacheWord
func printCacheWord(val int) string {
	strOut := strconv.FormatInt(int64(uint32(val)), 2)
	for len(strOut) < 32 {
		strOut = "0" + strOut
	}
	return strOut
}

// AccumulateDeps flags and returns additional dependent registers in the specified instruction array
func AccumulateDeps(instructions []int) uint32 {
	deps := uint32(0)
	for _, v := range instructions {
		if v != -1 {
			inst := InstructionArray[v]
			switch inst.typeofInstruction {
			case "R":
				deps = 1 << inst.rd //| 1<<inst.rm | 1<<inst.rn
			case "IM":
				deps = 1 << inst.rd
			case "I":
				deps = 1 << inst.rd //| 1<<inst.rn
			case "D":
				if inst.op == "LDUR" {
					deps = 1 << inst.rt
				}
			}
		}
	}
	return deps
}

// RegisterDeps accumulated the register dependencies for an instruction about to be executed
func RegisterDeps(instruction int) uint32 {
	inst := InstructionArray[instruction]
	deps := uint32(0)
	switch inst.typeofInstruction {
	case "R":
		deps = 1<<inst.rd | 1<<inst.rm | 1<<inst.rn
	case "IM":
		deps = 1 << inst.rd
	case "I":
		deps = 1<<inst.rd | 1<<inst.rn
	case "D":
		deps = 1<<inst.rt | 1<<inst.rn
	case "CB":
		deps = 1 << inst.conditional
	}
	return deps
}
