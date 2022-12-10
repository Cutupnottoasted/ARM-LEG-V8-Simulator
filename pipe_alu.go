package main

var PostALUBuffer = -1
var PostALUResult uint64 = 0

func SimALU() {
	if PreALUBuffer[0] != -1 {
		PostALUBuffer = PreALUBuffer[0]
		inst := InstructionArray[PreALUBuffer[0]]
		PreALUBuffer[0] = PreALUBuffer[1]
		PreALUBuffer[1] = -1

		immediate := uint64(inst.im)
		if inst.im < 0 {
			immediate = -uint64(-inst.im)
		}
		switch inst.op {
		case "AND":
			PostALUResult = RegisterFile[inst.rm] & RegisterFile[inst.rn]
		case "ADD":
			PostALUResult = RegisterFile[inst.rm] + RegisterFile[inst.rn]
		case "ADDI":
			PostALUResult = RegisterFile[inst.rn] + immediate
		case "ORR":
			PostALUResult = RegisterFile[inst.rm] | RegisterFile[inst.rn]
		case "SUB":
			PostALUResult = RegisterFile[inst.rn] - RegisterFile[inst.rm]
		case "SUBI":
			PostALUResult = RegisterFile[inst.rn] - immediate
		case "MOVZ":
			PostALUResult = uint64(inst.field) << inst.shift
		case "MOVK":
			PostALUResult = RegisterFile[inst.rd] ^ (RegisterFile[inst.rd] & (0xFFFF << inst.shift))
			PostALUResult = RegisterFile[inst.rd] | uint64(inst.field)<<inst.shift
		case "LSR":
			PostALUResult = RegisterFile[inst.rn] >> uint64(inst.shamt)
		case "LSL":
			PostALUResult = RegisterFile[inst.rn] << uint64(inst.shamt)
		case "ASR":
			sign := RegisterFile[inst.rn] >> 63
			PostALUResult = sign*(uint64(0xFFFFFFFFFFFFFFFF)<<uint64(64-inst.shamt)) | RegisterFile[inst.rn]>>inst.shamt
		case "EOR":
			PostALUResult = RegisterFile[inst.rm] ^ RegisterFile[inst.rn]
		}

	}
}
