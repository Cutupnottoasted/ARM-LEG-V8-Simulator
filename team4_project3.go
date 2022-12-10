// FILE: team4_project1.go
// DESCRIPTION: A program that takes 32 bit binary sequences and converts them to ARM V8 opcode
// using a structure array to store each opcode and their corresponding format
//
// Instruction structure variables:
// typeofInstruction string: Contains the type of instruction
// rawInstruction string: Contains the 32 bit sequence of binary
// lineValue uint32: Contains the unsigned decimal value of the 32 bit sequence
// programCnt int: Contains corresponding address for each Instruction entry
// opcode uint32: Contains the unsigned decimal value of bits 21 - 31
// op string: Contains the name of the operation (SUB, LDUR, etc)
// rd uint8: Contains the value of the destination register in R, I, IM format
// rn uint8: Contains the value of the base operand in R, D, I format
// rm uint8: Contains the value of the second operand
// shamt uint8: Contains the value of the shift amount in R format
// rt uint8: Contains the value of the target register in D format
// op2 uint8: Contains the value of the 2nd operator in D format
// address uint16: Contains the value of the address [rn + address] in D format
// im int16: Contains the signed value of immediate in I format
// offset int32: Contains the signed value of offset in B, CB format
// conditional uint8: Contains the conditional operator in CB format
// field uint16: Contains the value of a 16 bit pattern in IM format
// shift uint8: Contains the value of the shift code in IM format
//
// Global variables & function variables:
// FileIn "test1_bin.txt": A text file that contains 32 bit sequences
// FileOut "team4_out_dis.txt": A text file that contains the decoded and formatted 32 bit sequences
// InstructionCount: Keeps track of the number of Instructions that have been read and stored.
// DataSegmentIndex: Contains the end address of the Instruction segment and the start of program data
// InstructionArray []Instruction: An array of Instructions that contains each Instruction and corresponding data
// lineOut: Used to string 32 bit binary string into correct the format and append appropriate data
// binValue: Holds the total signed integer of the 32 bit binary string
//
// Functions:
// func readInstructions (input string):
// Pre: Check if text file exists. If the text file does not exist close the stream and output an error.
// Post: Reads entire text file and stores each line in InstructionArray.rawInstruction[index].
// When the BREAK instruction is read, store the index into DataSegmentIndex to mark the end of the data segment.
// If the scanner reads incorrectly output an error.
//
// func decodeInstruction (index int):
// Pre: Convert the InstructionArray.rawInstruction[index] with strconv.ParseUint.
// If the conversion fails, output an error
// Post: Reads InstructionArray.rawInstruction[index] and bitwise isolates the opcode from bit 21 - 31
// It then compares the InstructionArray.opcode[index] to the corresponding instruction format and uses
// bitwise isolation to store each value in the correct data variable. If operators im and offset are
// in 2C then convert else store values as is. Shift operator is multiplied by 16 and then stored.
//
// func printInstructions (output string):
// Pre: Create the output text file. If output file is not created, output an error and close the stream.
// Post: Takes InstructionArray[index] and iterates throughout the entire array from 0 to < InstructionCount &&
// < DataSegmentIndex. Each type of Instruction type is identified and is written in the correct format.
// Once the end of the InstructionArray is reached, use the BREAK function to signal the end of Instructions
// and write the program data 32 bit binary string and associated binValue
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

// Instruction structure variables:
// typeofInstruction string: Contains the type of instruction
// rawInstruction string: Contains the 32 bit sequence of binary
// lineValue uint32: Contains the unsigned decimal value of the 32 bit sequence
// programCnt int: Contains corresponding address for each Instruction entry
// opcode uint32: Contains the unsigned decimal value of bits 21 - 31
// op string: Contains the name of the operation (SUB, LDUR, etc)
// rd uint8: Contains the value of the destination register in R, I, IM format
// rn uint8: Contains the value of the base operand in R, D, I format
// rm uint8: Contains the value of the second operand
// shamt uint8: Contains the value of the shift amount in R format
// rt uint8: Contains the value of the target register in D format
// op2 uint8: Contains the value of the 2nd operator in D format
// address uint16: Contains the value of the address [rn + address] in D format
// im int16: Contains the signed value of immediate in I format
// offset int32: Contains the signed value of offset in B, CB format
// conditional uint8: Contains the conditional operator in CB format
// field uint16: Contains the value of a 16 bit pattern in IM format
// shift uint8: Contains the value of the shift code in IM format
type Instruction struct {
	typeofInstruction string
	rawInstruction    string
	lineValue         uint32
	decodedStr        string
	programCnt        int
	opcode            uint32
	op                string
	rd                uint8
	rn                uint8
	rm                uint8
	shamt             uint8
	rt                uint8
	op2               uint8
	address           uint16
	im                int16
	offset            int32
	conditional       uint8
	field             uint16
	shift             uint8
	destReg           int
	src1Reg           int
	src2Reg           int
}

// GLOBAL DATA VARIABLES

// FileIn "test1_bin.txt": A text file that contains 32 bit sequences
var FileIn = "team4_test1_bin.txt"

// FileOut "team4_out": A text file that contains the decoded and formatted 32 bit sequences
var FileOut = "team4_out"

// InstructionCount Keeps track of the number of Instructions that have been read and stored.
var InstructionCount = 0

// DataSegmentIndex Contains the end address of the instruction segment and the start of program data
var DataSegmentIndex = math.MaxInt32
var DataSegmentAddress = math.MaxInt32

// InstructionArray []Instruction: An array of Instructions that contains each Instruction and corresponding data
var InstructionArray []Instruction

// FUNCTIONS

// func readInstructions (input string):
// Pre: Check if text file exists. If the text file does not exist close the stream and output an error.
// Post: Reads entire text file and stores each line in InstructionArray.rawInstruction[index].
// When the BREAK instruction is read, store the index into DataSegmentIndex to mark the end of the data segment.
// If the scanner reads incorrectly output an error.
func readInstructions(input string) {
	file, err := os.Open(input)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		nextInstruction := Instruction{rawInstruction: scanner.Text()}
		InstructionArray = append(InstructionArray, nextInstruction)
		InstructionCount++

		// If BREAK is hit, mark start of data segment
		if nextInstruction.rawInstruction == "11111110110111101111111111100111" && DataSegmentIndex == math.MaxInt32 {
			DataSegmentIndex = InstructionCount
		}
	}

	// No BREAK, start data segment after last instruction
	if DataSegmentIndex == math.MaxInt32 {
		DataSegmentIndex = InstructionCount
	}
	DataSegmentAddress = DataSegmentIndex*4 + 96

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

// func decodeInstruction (index int):
// Pre: Convert the InstructionArray.rawInstruction[index] with strconv.ParseUint.
// If the conversion fails, output an error
// Post: Reads InstructionArray.rawInstruction[index] and bitwise isolates the opcode from bit 21 - 31
// It then compares the InstructionArray.opcode[index] to the corresponding instruction format and uses
// bitwise isolation to store each value in the correct data variable. If operators im and offset are
// in 2C then convert else store values as is. Shift operator is multiplied by 16 and then stored.
func decodeInstruction(index int) {

	lVal, err := strconv.ParseUint(InstructionArray[index].rawInstruction, 2, 32)
	if err != nil {
		fmt.Println(err)
	}

	InstructionArray[index].lineValue = uint32(lVal)

	InstructionArray[index].programCnt = 96 + (index * 4)

	InstructionArray[index].opcode = InstructionArray[index].lineValue >> 21

	if 160 <= InstructionArray[index].opcode && InstructionArray[index].opcode <= 191 {
		InstructionArray[index].typeofInstruction = "B"
		InstructionArray[index].op = "B"
		InstructionArray[index].offset = int32((InstructionArray[index].lineValue >> 0) & 0x3FFFFFF)
		if InstructionArray[index].offset&0x2000000 != 0 {
			InstructionArray[index].offset = -((InstructionArray[index].offset ^ 0x3FFFFFF) + 1)
		}
	} else if 1104 == InstructionArray[index].opcode {
		InstructionArray[index].typeofInstruction = "R"
		InstructionArray[index].op = "AND"
	} else if 1112 == InstructionArray[index].opcode {
		InstructionArray[index].typeofInstruction = "R"
		InstructionArray[index].op = "ADD"
	} else if 1160 == InstructionArray[index].opcode || 1161 == InstructionArray[index].opcode {
		InstructionArray[index].typeofInstruction = "I"
		InstructionArray[index].op = "ADDI"
	} else if 1360 == InstructionArray[index].opcode {
		InstructionArray[index].typeofInstruction = "R"
		InstructionArray[index].op = "ORR"
	} else if 1440 <= InstructionArray[index].opcode && InstructionArray[index].opcode <= 1447 {
		InstructionArray[index].typeofInstruction = "CB"
		InstructionArray[index].op = "CBZ"
	} else if 1448 <= InstructionArray[index].opcode && InstructionArray[index].opcode <= 1455 {
		InstructionArray[index].typeofInstruction = "CB"
		InstructionArray[index].op = "CBNZ"
	} else if InstructionArray[index].opcode == 1624 {
		InstructionArray[index].typeofInstruction = "R"
		InstructionArray[index].op = "SUB"
	} else if 1672 == InstructionArray[index].opcode || InstructionArray[index].opcode == 1673 {
		InstructionArray[index].typeofInstruction = "I"
		InstructionArray[index].op = "SUBI"
	} else if 1684 <= InstructionArray[index].opcode && InstructionArray[index].opcode <= 1687 {
		InstructionArray[index].typeofInstruction = "IM"
		InstructionArray[index].op = "MOVZ"
	} else if 1940 <= InstructionArray[index].opcode && InstructionArray[index].opcode <= 1943 {
		InstructionArray[index].typeofInstruction = "IM"
		InstructionArray[index].op = "MOVK"
	} else if InstructionArray[index].opcode == 1690 {
		InstructionArray[index].typeofInstruction = "R"
		InstructionArray[index].op = "LSR"
	} else if InstructionArray[index].opcode == 1691 {
		InstructionArray[index].typeofInstruction = "R"
		InstructionArray[index].op = "LSL"
	} else if InstructionArray[index].opcode == 1984 {
		InstructionArray[index].typeofInstruction = "D"
		InstructionArray[index].op = "STUR"
	} else if InstructionArray[index].opcode == 1986 {
		InstructionArray[index].typeofInstruction = "D"
		InstructionArray[index].op = "LDUR"
	} else if InstructionArray[index].opcode == 1692 {
		InstructionArray[index].typeofInstruction = "R"
		InstructionArray[index].op = "ASR"
	} else if InstructionArray[index].opcode == 0 {
		InstructionArray[index].typeofInstruction = "N/A"
		InstructionArray[index].op = "NOP"
	} else if InstructionArray[index].opcode == 1872 {
		InstructionArray[index].typeofInstruction = "R"
		InstructionArray[index].op = "EOR"
	} else if InstructionArray[index].opcode == 2038 {
		InstructionArray[index].typeofInstruction = "BREAK"
		InstructionArray[index].op = "BREAK"
	} else {
		InstructionArray[index].typeofInstruction = "UNKNOWN"
		InstructionArray[index].op = "UNKNOWN"
	}

	if InstructionArray[index].typeofInstruction == "R" {
		InstructionArray[index].rd = uint8((InstructionArray[index].lineValue >> 0) & 0x1F)
		InstructionArray[index].rn = uint8((InstructionArray[index].lineValue >> 5) & 0x1F)
		InstructionArray[index].shamt = uint8((InstructionArray[index].lineValue >> 10) & 0x3F)
		InstructionArray[index].rm = uint8((InstructionArray[index].lineValue >> 16) & 0x1F)
	} else if InstructionArray[index].typeofInstruction == "D" {
		InstructionArray[index].rt = uint8((InstructionArray[index].lineValue >> 0) & 0x1F)
		InstructionArray[index].rn = uint8((InstructionArray[index].lineValue >> 5) & 0x1F)
		InstructionArray[index].op2 = uint8((InstructionArray[index].lineValue >> 10) & 0x3)
		InstructionArray[index].address = uint16((InstructionArray[index].lineValue >> 12) & 0x1FF)
	} else if InstructionArray[index].typeofInstruction == "CB" {
		InstructionArray[index].conditional = uint8((InstructionArray[index].lineValue >> 0) & 0x1F)
		InstructionArray[index].offset = int32((InstructionArray[index].lineValue >> 5) & 0x7FFFF)
		if InstructionArray[index].offset&0x40000 != 0 {
			InstructionArray[index].offset = -((InstructionArray[index].offset ^ 0x7FFFF) + 1)
		}
	} else if InstructionArray[index].typeofInstruction == "I" {
		InstructionArray[index].rd = uint8((InstructionArray[index].lineValue >> 0) & 0x1F)
		InstructionArray[index].rn = uint8((InstructionArray[index].lineValue >> 5) & 0x1F)
		InstructionArray[index].im = int16((InstructionArray[index].lineValue >> 10) & 0xFFF)
		// Project 2 spec states that immediate value is positive
		//if InstructionArray[index].im&0x800 != 0 {
		//	InstructionArray[index].im = -((InstructionArray[index].im ^ 0xFFF) + 1)
		//}
	} else if InstructionArray[index].typeofInstruction == "IM" {
		InstructionArray[index].rd = uint8((InstructionArray[index].lineValue >> 0) & 0x1F)
		InstructionArray[index].field = uint16((InstructionArray[index].lineValue >> 5) & 0xFFFF)
		InstructionArray[index].shift = uint8((InstructionArray[index].lineValue>>21)&0x3) * 16
	}
	InstructionArray[index].destReg = int(InstructionArray[index].rd)
	InstructionArray[index].src1Reg = int(InstructionArray[index].rn)
	InstructionArray[index].src2Reg = int(InstructionArray[index].rm)
	if InstructionArray[index].op == "LDUR" {
		InstructionArray[index].destReg = int(InstructionArray[index].rt)
	}

	// TODO: placeholder value, overwrite at print stage
	InstructionArray[index].decodedStr = strconv.Itoa(int(InstructionArray[index].lineValue))
}

// PrintInstructions (output string):
// Pre: Create the output text file. If output file is not created, output an error and close the stream.
// Post: Takes InstructionArray[index] and iterates throughout the entire array from 0 to < InstructionCount &&
// < DataSegmentIndex. Each type of Instruction type is identified and is written in the correct format.
// Once the end of the InstructionArray is reached, use the BREAK function to signal the end of Instructions
// and write the program data 32 bit binary string and associated binValue
func PrintInstructions(output string) {
	file, err := os.Create(output + "_dis.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	for i := 0; i < InstructionCount && i < DataSegmentIndex; i++ {
		var lineOut = ""
		if InstructionArray[i].typeofInstruction == "R" {
			lineOut += InstructionArray[i].rawInstruction[0:11] + " "
			lineOut += InstructionArray[i].rawInstruction[11:16] + " "
			lineOut += InstructionArray[i].rawInstruction[16:22] + " "
			lineOut += InstructionArray[i].rawInstruction[22:27] + " "
			lineOut += InstructionArray[i].rawInstruction[27:32] + "\t"
			lineOut += strconv.Itoa(InstructionArray[i].programCnt) + "\t"

			InstructionArray[i].decodedStr = InstructionArray[i].op + "\t"

			InstructionArray[i].decodedStr += "R" + strconv.Itoa(int(InstructionArray[i].rd)) + ", "

			if InstructionArray[i].op[1] == 'S' {
				InstructionArray[i].decodedStr += "R" + strconv.Itoa(int(InstructionArray[i].rn)) + ", "
				InstructionArray[i].decodedStr += "#" + strconv.Itoa(int(InstructionArray[i].shamt))
			} else {
				InstructionArray[i].decodedStr += "R" + strconv.Itoa(int(InstructionArray[i].rn)) + ", "
				InstructionArray[i].decodedStr += "R" + strconv.Itoa(int(InstructionArray[i].rm))
			}

			lineOut += InstructionArray[i].decodedStr
		} else if InstructionArray[i].typeofInstruction == "D" {
			lineOut += InstructionArray[i].rawInstruction[0:11] + " "
			lineOut += InstructionArray[i].rawInstruction[11:20] + " "
			lineOut += InstructionArray[i].rawInstruction[20:22] + " "
			lineOut += InstructionArray[i].rawInstruction[22:27] + " "
			lineOut += InstructionArray[i].rawInstruction[27:32] + "\t"
			lineOut += strconv.Itoa(InstructionArray[i].programCnt) + "\t"

			InstructionArray[i].decodedStr = InstructionArray[i].op + "\t"

			InstructionArray[i].decodedStr += "R" + strconv.Itoa(int(InstructionArray[i].rt)) + ", ["
			InstructionArray[i].decodedStr += "R" + strconv.Itoa(int(InstructionArray[i].rn)) + ", "
			InstructionArray[i].decodedStr += "#" + strconv.Itoa(int(InstructionArray[i].address)) + "]"

			lineOut += InstructionArray[i].decodedStr
		} else if InstructionArray[i].typeofInstruction == "I" {
			lineOut += InstructionArray[i].rawInstruction[0:10] + " "
			lineOut += InstructionArray[i].rawInstruction[10:22] + " "
			lineOut += InstructionArray[i].rawInstruction[22:27] + " "
			lineOut += InstructionArray[i].rawInstruction[27:32] + "\t"
			lineOut += strconv.Itoa(InstructionArray[i].programCnt) + "\t"

			InstructionArray[i].decodedStr = InstructionArray[i].op + "\t"

			InstructionArray[i].decodedStr += "R" + strconv.Itoa(int(InstructionArray[i].rd)) + ", "
			InstructionArray[i].decodedStr += "R" + strconv.Itoa(int(InstructionArray[i].rn)) + ", "
			InstructionArray[i].decodedStr += "#" + strconv.Itoa(int(InstructionArray[i].im))

			lineOut += InstructionArray[i].decodedStr
		} else if InstructionArray[i].typeofInstruction == "B" {
			lineOut += InstructionArray[i].rawInstruction[0:6] + " "
			lineOut += InstructionArray[i].rawInstruction[6:32] + "\t"
			lineOut += strconv.Itoa(InstructionArray[i].programCnt) + "\t"

			InstructionArray[i].decodedStr = InstructionArray[i].op + "\t"

			InstructionArray[i].decodedStr += "#" + strconv.Itoa(int(InstructionArray[i].offset))

			lineOut += InstructionArray[i].decodedStr
		} else if InstructionArray[i].typeofInstruction == "CB" {
			lineOut += InstructionArray[i].rawInstruction[0:8] + " "
			lineOut += InstructionArray[i].rawInstruction[8:27] + " "
			lineOut += InstructionArray[i].rawInstruction[27:32] + "\t"
			lineOut += strconv.Itoa(InstructionArray[i].programCnt) + "\t"

			InstructionArray[i].decodedStr = InstructionArray[i].op + "\t"

			InstructionArray[i].decodedStr += "R" + strconv.Itoa(int(InstructionArray[i].rd)) + ", "
			InstructionArray[i].decodedStr += "#" + strconv.Itoa(int(InstructionArray[i].offset))

			lineOut += InstructionArray[i].decodedStr
		} else if InstructionArray[i].typeofInstruction == "IM" {
			lineOut += InstructionArray[i].rawInstruction[0:9] + " "
			lineOut += InstructionArray[i].rawInstruction[9:11] + " "
			lineOut += InstructionArray[i].rawInstruction[11:27] + " "
			lineOut += InstructionArray[i].rawInstruction[27:32] + "\t"
			lineOut += strconv.Itoa(InstructionArray[i].programCnt) + "\t"

			InstructionArray[i].decodedStr = InstructionArray[i].op + "\t"

			InstructionArray[i].decodedStr += "R" + strconv.Itoa(int(InstructionArray[i].rd)) + ", "
			InstructionArray[i].decodedStr += strconv.Itoa(int(InstructionArray[i].field)) + ", "
			InstructionArray[i].decodedStr += "LSL " + strconv.Itoa(int(InstructionArray[i].shift))

			lineOut += InstructionArray[i].decodedStr
		} else if InstructionArray[i].typeofInstruction == "BREAK" {
			lineOut += "\n"
			lineOut += InstructionArray[i].rawInstruction[0:8]
			lineOut += InstructionArray[i].rawInstruction[8:11]
			lineOut += InstructionArray[i].rawInstruction[11:16]
			lineOut += InstructionArray[i].rawInstruction[16:21]
			lineOut += InstructionArray[i].rawInstruction[21:26]
			lineOut += InstructionArray[i].rawInstruction[26:32] + "\t"
			lineOut += strconv.Itoa(InstructionArray[i].programCnt) + "\t"

			InstructionArray[i].decodedStr = InstructionArray[i].op + "\t"

			lineOut += InstructionArray[i].decodedStr
		} else {
			lineOut += InstructionArray[i].rawInstruction[0:32] + "\t"
			lineOut += strconv.Itoa(InstructionArray[i].programCnt) + "\t"
			InstructionArray[i].decodedStr = InstructionArray[i].op + "\t"

			lineOut += InstructionArray[i].decodedStr
		}

		fmt.Fprintln(file, lineOut)
	}

	for i := DataSegmentIndex; i < InstructionCount; i++ {
		var lineOut = InstructionArray[i].rawInstruction + "\t"
		lineOut += strconv.Itoa(InstructionArray[i].programCnt) + "\t"

		var binValue = int64(InstructionArray[i].lineValue)
		if binValue&0x80000000 != 0 {
			binValue = -((binValue ^ 0xffffffff) + 1)
		}
		lineOut += strconv.FormatInt(binValue, 10)

		fmt.Fprintln(file, lineOut)
	}
}

func main() {

	for i, arg := range os.Args[1:] {
		if arg == "-i" {
			i++
			FileIn = os.Args[i+1]
		} else if arg == "-o" {
			i++
			FileOut = os.Args[i+1]
		}
	}

	//fmt.Println("Input:", FileIn)
	//fmt.Println("Output:", FileOut)

	readInstructions(FileIn)

	for i := 0; i < InstructionCount; i++ {
		decodeInstruction(i)
	}

	PrintInstructions(FileOut)
	PrintMachineSim(FileOut)
	SimulatePipeline(FileOut)
}
