package main

import (
	"fmt"
	"os"
	"strconv"
)

// CycleLimit Hard limit cycle count to avoid runaway
const CycleLimit = 10000

// xzr zero register
const xzr uint8 = 0

type MachineState struct {
	// regFile contains registers 00-31
	regFile [32]uint64

	// mem contains simulated memory space
	mem      SimMemory
	pc       int
	breakNow bool
}

// evaluate advances the machine state by one cycle
// Returns New machine state and break
func evaluate(state MachineState) MachineState {
	inst := InstructionArray[(state.pc-96)>>2]

	// branchNow stores whether a branch is taken this cycle
	branched := false

	immediate := uint64(inst.im)
	if inst.im < 0 {
		immediate = -uint64(-inst.im)
	}
	switch inst.op {
	case "B":
		state.pc = int(inst.offset<<2) + state.pc
		branched = true
	case "AND":
		state.regFile[inst.rd] = state.regFile[inst.rm] & state.regFile[inst.rn]
	case "ADD":
		state.regFile[inst.rd] = state.regFile[inst.rm] + state.regFile[inst.rn]
	case "ADDI":
		state.regFile[inst.rd] = state.regFile[inst.rn] + immediate
	case "ORR":
		state.regFile[inst.rd] = state.regFile[inst.rm] | state.regFile[inst.rn]
	case "CBZ":
		if state.regFile[inst.conditional] == 0 {
			state.pc = int(inst.offset<<2) + state.pc
			branched = true
		}
	case "CBNZ":
		if state.regFile[inst.conditional] != 0 {
			state.pc = int(inst.offset<<2) + state.pc
			branched = true
		}
	case "SUB":
		state.regFile[inst.rd] = state.regFile[inst.rn] - state.regFile[inst.rm]
	case "SUBI":
		state.regFile[inst.rd] = state.regFile[inst.rn] - immediate
	case "MOVZ":
		state.regFile[inst.rd] = uint64(inst.field) << inst.shift
	case "MOVK":
		state.regFile[inst.rd] = state.regFile[inst.rd] ^ (state.regFile[inst.rd] & (0xFFFF << inst.shift))
		state.regFile[inst.rd] = state.regFile[inst.rd] | uint64(inst.field)<<inst.shift
	case "LSR":
		state.regFile[inst.rd] = state.regFile[inst.rn] >> uint64(inst.shamt)
	case "LSL":
		state.regFile[inst.rd] = state.regFile[inst.rn] << uint64(inst.shamt)
	case "STUR":
		state.mem = WriteSimMemory(state.mem, int(state.regFile[inst.rn])+(int(inst.address)*4), state.regFile[inst.rt])
	case "LDUR":
		state.regFile[inst.rt] = ReadSimMemory(state.mem, int(state.regFile[inst.rn])+(int(inst.address)*4))
	case "ASR":
		sign := state.regFile[inst.rn] >> 63
		state.regFile[inst.rd] = sign*(uint64(0xFFFFFFFFFFFFFFFF)<<uint64(64-inst.shamt)) | state.regFile[inst.rn]>>inst.shamt
	case "EOR":
		state.regFile[inst.rd] = state.regFile[inst.rm] ^ state.regFile[inst.rn]

	case "BREAK":
		state.breakNow = true
	}

	if branched == false {
		state.pc += 4
	}

	if state.pc >= DataSegmentIndex*4+96 || state.pc < 96 {
		state.breakNow = true
	}

	return state
}

// PrintMachineSim runs evaluate command until break condition is met
// Outputs machine state at each cycle to the specified output file
func PrintMachineSim(output string) {
	file, err := os.Create(output + "_sim.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	simState := MachineState{pc: 96, breakNow: false, mem: InitSimMemory()}
	cycle := 1
	for simState.breakNow == false && simState.pc < DataSegmentIndex*4+96 && cycle < CycleLimit {
		// 20 equal signs and a newline
		lineOut := "====================\n"

		// cycle: [cycle number] [tab] [instruction address] [tab] [instruction string (same as step 3 above)]
		lineOut += "cycle:" + strconv.Itoa(cycle) + "\t" + strconv.Itoa(simState.pc) + "\t"
		lineOut += InstructionArray[(simState.pc-96)>>2].decodedStr + "\n"

		simState = evaluate(simState)

		// [blank line]
		lineOut += "\n"

		// registers:
		lineOut += "registers:\n"

		// print register contents
		for i, reg := range simState.regFile {
			if i == 0 || i == 8 {
				lineOut += "r0" + strconv.Itoa(i)
			} else if i == 16 || i == 24 {
				lineOut += "r" + strconv.Itoa(i)
			}

			lineOut += "\t" + strconv.FormatInt(int64(reg), 10)

			if i%8 == 7 {
				lineOut += "\n"
			}
		}

		// blank line
		lineOut += "\n"

		// data:
		lineOut += "data:"

		// Print memory data
		for i := DataSegmentIndex; i < len(simState.mem.memFile); i++ {
			if i%8 == DataSegmentIndex%8 {
				lineOut += "\n" + strconv.Itoa(i*4+96) + ":"
			}
			lineOut += "\t" + strconv.Itoa(int(simState.mem.memFile[i]))
		}

		lineOut += "\n"

		fmt.Fprint(file, lineOut)

		cycle++
	}

	if cycle >= CycleLimit {
		fmt.Println("Error: Simulation terminated due to runaway protection")
	}
}
