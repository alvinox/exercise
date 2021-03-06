package code

import (
    "fmt"
    "bytes"
    "encoding/binary"
)

type Instructions []byte

type Opcode byte

const (
    OpConstant Opcode = iota
    OpPop
    OpAdd
    OpSub
    OpMul
    OpDiv
    OpTrue
    OpFalse
    OpNull
    OpEqual
    OpNotEqual
    OpGreaterThan
    OpMinus
    OpBang
    OpIndex
    OpJumpNotTruthy
    OpJump
    OpGetGlobal
    OpSetGlobal
    OpGetLocal
    OpSetLocal
    OpArray
    OpHash
    OpCall
    OpReturnValue
    OpReturn
    OpGetBuiltin
    OpClosure
    OpGetFree
)

const (
    ConstWidth         int = 2 // const pool max size is 65536
    GlobalWidth        int = 2 // max number of global variables is 65536
    LocalWidth         int = 1 // max number of local variables is 256
    FreeWidth          int = 1 // max number of free variables is 256
    BuiltinWidth       int = 1 // max number of builtin functions is 256
    InstructionWidth   int = 2 // max number of instructions is 65536
    CallParamWidth     int = 1 // max number of parameters of each function is 256
)

type Definition struct {
    Name          string
    OperandWidths []int
}

var definitions = map[Opcode] *Definition {
    OpConstant:     {"OpConstant",      []int{ConstWidth}},
    OpPop:          {"OpPop",           []int{}},
    OpAdd:          {"OpAdd",           []int{}},
    OpSub:          {"OpSub",           []int{}},
    OpMul:          {"OpMul",           []int{}},
    OpDiv:          {"OpDiv",           []int{}},
    OpTrue:         {"OpTrue",          []int{}},
    OpFalse:        {"OpFlase",         []int{}},
    OpNull:         {"OpNull",          []int{}},
    OpEqual:        {"OpEqual",         []int{}},
    OpNotEqual:     {"OpNotEqual",      []int{}},
    OpGreaterThan:  {"OpGreaterThan",   []int{}},
    OpMinus:        {"OpMinus",         []int{}},
    OpBang:         {"OpBang",          []int{}},
    OpIndex:        {"OpIndex",         []int{}},
    OpJumpNotTruthy:{"OpJumpNotTruthy", []int{InstructionWidth}},
    OpJump:         {"OpJump",          []int{InstructionWidth}},
    OpGetGlobal:    {"OpGetGlobal",     []int{GlobalWidth}},
    OpSetGlobal:    {"OpSetGlobal",     []int{GlobalWidth}},
    OpGetLocal:     {"OpGetLocal",      []int{LocalWidth}},
    OpSetLocal:     {"OpSetLocal",      []int{LocalWidth}},
    OpArray:        {"OpArray",         []int{GlobalWidth}},
    OpHash:         {"OpHash",          []int{GlobalWidth}},
    OpCall:         {"OpCall",          []int{CallParamWidth}},
    OpReturnValue:  {"OpReturnValue",   []int{}},
    OpReturn:       {"OpReturn",        []int{}},
    OpGetBuiltin:   {"OpGetBuiltin",    []int{BuiltinWidth}},
    OpClosure:      {"OpClosure",       []int{ConstWidth, FreeWidth}},
    OpGetFree:      {"OpGetFree",       []int{FreeWidth}},
}

func Lookup(op byte) (*Definition, error) {
    def, ok := definitions[Opcode(op)]
    if !ok {
        return nil, fmt.Errorf("opcode %d undefined", op)
    }

    return def, nil
}

func Make(op Opcode, operands ...int) []byte {
    def, ok := definitions[op]
    if !ok {
        return []byte{}
    }

    instructionsLen := 1
    for _, w := range def.OperandWidths {
        instructionsLen += w
    }

    instructions := make([]byte, instructionsLen)
    instructions[0] = byte(op)

    offset := 1
    for i, o := range operands {
        width := def.OperandWidths[i]
        switch width {
            case 2:
                binary.BigEndian.PutUint16(instructions[offset:], uint16(o))
            case 1:
                instructions[offset] = byte(o)
            
        }
        offset += width
    }

    return instructions
}

func (ins Instructions) String() string {
    var out bytes.Buffer

    i := 0
    for i < len(ins) {
        def, err := Lookup(ins[i])
        if err != nil {
            fmt.Fprintf(&out, "Error: %s\n", err)
            continue
        }

        operands, read := ReadOperands(def, ins[i+1:])

        fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

        i += 1 + read
    }

    return out.String()
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
    operands := make([]int, len(def.OperandWidths))
    offset := 0

    for i, width := range def.OperandWidths {
        switch width {
            case 2:
                operands[i] = int(ReadUint16(ins[offset:]))
            case 1:
                operands[i] = int(ReadUint8(ins[offset:]))
        }

        offset += width
    }

    return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
    return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 {
    return uint8(ins[0])
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
    operandCount := len(def.OperandWidths)

    if len(operands) != operandCount {
        return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
            len(operands), operandCount)
    }

    switch operandCount {
        case 0:
            return def.Name
        case 1:
            return fmt.Sprintf("%s %d", def.Name, operands[0])
        case 2:
            return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
    }

    return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}