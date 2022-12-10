package main

var PreMemBuffer = []int{-1, -1}
var PreALUBuffer = []int{-1, -1}

func SimIssue() {
	// Issue 1
	issuePass()
	issuePass()
	// Issue 2
	//issuePassMEM()
}

func issuePass() {
	// Has an instruction been issued this pass
	issued := false
	// Contains 32 flags for each of the upstream register dependencies
	deps := AccumulateDeps(PreALUBuffer) | AccumulateDeps([]int{PostALUBuffer})
	deps = deps | AccumulateDeps(PreMemBuffer) | AccumulateDeps([]int{PostMemBuffer})
	// If a stur has previously been encountered this pass
	sturHazard := false
	for i, v := range PreIssueBuffer {
		if v == -1 {
			break
		} else if issued {
			PreIssueBuffer[i-1] = PreIssueBuffer[i]
			PreIssueBuffer[i] = -1
		} else if InstructionArray[v].typeofInstruction == "D" {
			//TODO Handle Memory instruction hazards
			if RegisterDeps(v)&deps == 0 {
				if InstructionArray[v].op == "STUR" && sturHazard {
					continue
				} else if InstructionArray[v].op == "STUR" {
					sturHazard = true
				}
				addrDep := false
				for _, j := range PreMemBuffer {
					addrDep = addrDep || addressDependency(v, j)
				}
				for idx := 0; idx < i; i++ {
					addrDep = addrDep || addressDependency(v, PreIssueBuffer[idx])
				}
				if PreMemBuffer[0] == -1 {
					PreMemBuffer[0] = v
					issued = true
					PreIssueBuffer[i] = -1
				} else if PreMemBuffer[1] == -1 {
					PreMemBuffer[1] = v
					issued = true
					PreIssueBuffer[i] = -1
				}

			}

		} else {
			// Check hazards
			if RegisterDeps(v)&deps == 0 {
				if PreALUBuffer[0] == -1 {
					PreALUBuffer[0] = v
					issued = true
					PreIssueBuffer[i] = -1
				} else if PreALUBuffer[1] == -1 {
					PreALUBuffer[1] = v
					issued = true
					PreIssueBuffer[i] = -1
				}
			}
			deps = deps | AccumulateDeps([]int{v})
		}
	}
}

// addressDependency conflict with downstream address?
func addressDependency(i1 int, i2 int) bool {
	if i1 == -1 || i2 == -1 {
		return false
	} else {
		inst1 := InstructionArray[i1]
		inst2 := InstructionArray[i2]

		return int(RegisterFile[inst1.rn])+(int(inst1.address)*4) == int(RegisterFile[inst2.rn])+(int(inst2.address)*4)
	}
}
