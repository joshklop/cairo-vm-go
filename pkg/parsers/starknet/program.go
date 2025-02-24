package vm

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	f "github.com/consensys/gnark-crypto/ecc/stark-curve/fp"
)

type Builtin uint8

const (
	Output Builtin = iota + 1
	RangeCheck
	Pedersen
	ECDSA
	Keccak
	Bitwise
	ECOP
	Poseidon
	SegmentArena
)

func (b Builtin) MarshalJSON() ([]byte, error) {
	switch b {
	case Output:
		return []byte("output"), nil
	case RangeCheck:
		return []byte("range_check"), nil
	case Pedersen:
		return []byte("pedersen"), nil
	case ECDSA:
		return []byte("ecdsa"), nil
	case Keccak:
		return []byte("keccak"), nil
	case Bitwise:
		return []byte("bitwise"), nil
	case ECOP:
		return []byte("ec_op"), nil
	case Poseidon:
		return []byte("poseidon"), nil
	case SegmentArena:
		return []byte("segment_arena"), nil

	}
	return nil, fmt.Errorf("error marshaling builtin with unknow identifer: %d", uint8(b))
}

func (b *Builtin) UnmarshalJSON(data []byte) error {
	builtinName, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("error unmarsahling builtin: %w", err)
	}

	switch builtinName {
	case "output":
		*b = Output
	case "range_check":
		*b = RangeCheck
	case "pedersen":
		*b = Pedersen
	case "ecdsa":
		*b = ECDSA
	case "keccak":
		*b = Keccak
	case "bitwise":
		*b = Bitwise
	case "ec_op":
		*b = ECOP
	case "poseidon":
		*b = Poseidon
	case "segment_arena":
		*b = SegmentArena
	default:
		return fmt.Errorf("error unmarsahling unknwon builtin name: %s", builtinName)
	}
	return nil
}

type EntryPointInfo struct {
	Selector f.Element `json:"selector"`
	Offset   f.Element `json:"offset"`
	Builtins []Builtin `json:"builtins"`
}

type EntryPointByType struct {
	External    []EntryPointInfo `json:"EXTERNAL"`
	L1Handler   []EntryPointInfo `json:"L1_HANDLER"`
	Constructor []EntryPointInfo `json:"CONSTRUCTOR"`
}

type Hints struct {
	Index uint64
	Hints []Hint
}

// Hints are serialized as tuples of (index, []hint)
// https://github.com/starkware-libs/cairo/blob/main/crates/cairo-lang-starknet/src/casm_contract_class.rs#L90
func (hints *Hints) UnmarshalJSON(data []byte) error {
	var rawHints []interface{}
	if err := json.Unmarshal(data, &rawHints); err != nil {
		return err
	}

	index, ok := rawHints[0].(float64)
	if !ok {
		return fmt.Errorf("error unmarshaling hints: index should be uint64")
	}
	hints.Index = uint64(index)

	rest, err := json.Marshal(rawHints[1])
	if err != nil {
		return err
	}

	var h []Hint
	if err := json.Unmarshal(rest, &h); err != nil {
		return err
	}
	hints.Hints = h
	return nil
}

func (hints *Hints) MarshalJSON() ([]byte, error) {
	var rawHints []interface{}
	rawHints = append(rawHints, hints.Index)
	rawHints = append(rawHints, hints.Hints)

	return json.Marshal(rawHints)
}

type Program struct {
	// Prime is fixed to be 0x800000000000011000000000000000000000000000000000000000000000001 and wont fit in a f.Felt
	Bytecode        []f.Element      `json:"bytecode"`
	CompilerVersion string           `json:"compiler_version"`
	EntryPoints     EntryPointByType `json:"entry_points_by_type"`
	Hints           []Hints          `json:"hints" validate:"required"`
}

func ProgramFromFile(pathToFile string) (*Program, error) {
	content, error := os.ReadFile(pathToFile)
	if error != nil {
		return nil, error
	}
	return ProgramFromJSON(content)
}

func ProgramFromJSON(content json.RawMessage) (*Program, error) {
	var program Program
	return &program, json.Unmarshal(content, &program)
}
