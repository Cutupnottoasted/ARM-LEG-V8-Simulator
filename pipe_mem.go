package main

var PostMemBuffer = -1
var PostMemResult uint64 = 0

func SimMem() {
	if PreMemBuffer[0] != -1 {
		inst := InstructionArray[PreMemBuffer[0]]
		if inst.op == "LDUR" {
			hit, data := CacheRead(int(RegisterFile[inst.rn]) + (int(inst.address) * 4))
			if hit {
				PostMemBuffer = PreMemBuffer[0]
				PostMemResult = data
				PreMemBuffer[0] = PreMemBuffer[1]
				PreMemBuffer[1] = -1
			}
		} else { // if inst.op == "STUR"
			addr := int(RegisterFile[inst.rn]) + (int(inst.address) * 4)
			data := RegisterFile[inst.rt]
			hit := CacheWrite(addr, data)
			if hit {
				//MainMemory = WriteSimMemory(MainMemory, addr, data)
				PreMemBuffer[0] = PreMemBuffer[1]
				PreMemBuffer[1] = -1
			}
		}
	}
}
