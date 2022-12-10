package main

func SimWriteBack() {
	if PostMemBuffer != -1 {
		RegisterFile[InstructionArray[PostMemBuffer].destReg] = PostMemResult
		PostMemBuffer = -1
		PostMemResult = uint64(0)
	}
	if PostALUBuffer != -1 {
		RegisterFile[InstructionArray[PostALUBuffer].destReg] = PostALUResult
		PostALUBuffer = -1
		PostALUResult = uint64(0)
	}
}
