package main

import "fmt"

var PreIssueBuffer = []int{-1, -1, -1, -1}
var pipeBreak = false

func SimFetch() {
	if PreIssueBuffer[3] == -1 && pipeBreak == false {
		if ProgramCounter >= DataSegmentAddress || ProgramCounter < 96 {
			pipeBreak = true
			fmt.Println("ERROR: Program Counter escaped instruction segment")
			return
		} else if pipeBreak {
			return
		}

		hit, _ := CacheRead(ProgramCounter ^ (ProgramCounter & 0x7))
		if hit {
			nextInst := InstructionArray[(ProgramCounter-96)/4]
			isBreak, pcOffset := evaluateNext(nextInst, ProgramCounter)

			pipeBreak = isBreak
			ProgramCounter += pcOffset
			// TODO Handle PC instruction segment bounds
			if ProgramCounter >= DataSegmentAddress || ProgramCounter < 96 {
				pipeBreak = true
				fmt.Println("ERROR: Program Counter escaped instruction segment")
				return
			} else if PreIssueBuffer[3] == -1 && pipeBreak == false {
				//Cache testing only
				//hit, _ = CacheRead(ProgramCounter)
				hit, _ = CacheRead(ProgramCounter ^ (ProgramCounter & 0x7))
				if hit {
					nextInst = InstructionArray[(ProgramCounter-96)/4]
					pipeBreak, pcOffset = evaluateNext(nextInst, ProgramCounter)
					ProgramCounter += pcOffset
				}
			}
		}
	}
}

// evaluateNext checks the provided
func evaluateNext(inst Instruction, pc int) (bool, int) {
	isBreak := inst.op == "BREAK"
	offset := 0

	if isBreak == false {
		if inst.typeofInstruction == "B" {
			offset = int(inst.offset << 2)
		} else if inst.typeofInstruction == "CB" {
			// Accumulate and compare downstream dependencies
			deps := AccumulateDeps(PreALUBuffer) | AccumulateDeps([]int{PostALUBuffer})
			deps = deps | AccumulateDeps(PreMemBuffer) | AccumulateDeps([]int{PostMemBuffer})
			deps = deps | AccumulateDeps(PreIssueBuffer)
			if 1<<inst.conditional&deps == uint32(0) {
				if inst.typeofInstruction == "CBZ" {
					if RegisterFile[inst.conditional] == 0 {
						offset = int(inst.offset << 2)
					} else {
						offset = 4
					}
				} else { // CBNZ
					if RegisterFile[inst.conditional] != 0 {
						offset = int(inst.offset << 2)
					} else {
						offset = 4
					}
				}
			}
		} else if inst.typeofInstruction == "NOP" {
			offset = 4
		} else {
			pushToPreIssue((pc - 96) / 4)
			offset = 4
		}
	}

	return isBreak, offset
}

// pushToPreIssue adds one entry to the PreIssueBuffer in order
func pushToPreIssue(idx int) {
	for i, v := range PreIssueBuffer {
		if v == -1 {
			PreIssueBuffer[i] = idx
			break
		}
	}
}
