package inst

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var (
	INST_TERM string = ";"
	ARG_SEP   string = ","
	ELEM_SEP  string = "."
)

/*
Decode an Instruction arg.

example:
arg := DecodeArg("4.size")
arg == "size" = true

:return: string
*/
func DecodeArg(encoded_arg string) (string, error) {
	var arg string
	var err error

	if strings.Index(encoded_arg, ELEM_SEP) == -1 {
		return arg, errors.New("Instruction element separator not found.")
	}

	splited_arg := strings.Split(encoded_arg, ELEM_SEP)

	length, err := strconv.Atoi(splited_arg[0])
	if err != nil {
		return arg, errors.New("Instruction element contains invalid length.")
	}

	if len(splited_arg[1]) != length {
		return arg, fmt.Errorf("Instruction arg (%s) has invalid length.", encoded_arg)
	}

	return splited_arg[1], nil
}

/*
Encode argument to be sent in a valid Guacamole instruction.

example:
arg = EncodeArg('size')
arg == '4.size' = true

:return: string

*/
func EncodeArg(arg string) string {
	arg_length := len(arg)
	arg_length_str := strconv.Itoa(arg_length)
	return strings.Join([]string{arg_length_str, arg}, ELEM_SEP)
}

type Instruction struct {
	OpCode string
	Args   []string
}

func NewInstruction(opcode string, args ...string) *Instruction {
	instruction := Instruction{OpCode: opcode, Args: args}
	return &instruction
}

func LoadInstruction(inst string) (*Instruction, error) {
	var instruction Instruction
	var args []string

	if !strings.HasSuffix(inst, INST_TERM) {
		return &instruction, errors.New("Instruction termination not found.")
	}

	elems := strings.Split(string(inst[0:len(inst)-1]), ARG_SEP)

	for index := range elems {
		decoded_arg, err := DecodeArg(elems[index])
		if err != nil {
			log.Println(err)
		} else {
			args = append(args, decoded_arg)
		}
	}

	instruction.OpCode = args[0]
	instruction.Args = args[1:]

	return &instruction, nil
}

func (this *Instruction) Encode() string {
	instruction_iter := make([]string, len(this.Args)+1)

	instruction_iter[0] = this.OpCode

	for index := range this.Args {
		instruction_iter[index+1] = this.Args[index]
	}

	for index, arg := range instruction_iter {
		instruction_iter[index] = EncodeArg(arg)
	}

	elems := strings.Join(instruction_iter, ARG_SEP)

	return elems + INST_TERM
}

func (this *Instruction) String() string {
	return this.Encode()
}

/*
func main() {
	ins, err := LoadInstruction("4.size,4.1024,3.768;")
	fmt.Println(ins.OpCode, ins.Args)
	fmt.Println(err)
	fmt.Println(ins.Encode())
}
*/
