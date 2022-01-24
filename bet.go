package bet

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
)

const PANIC_UNMATCHED_TYPE = "unmatched_type"

type Node interface {
	Eval(p map[string]interface{}) bool
	Serialize() ([]byte, error)
}

type BinaryOprator string

const (
	OpAND BinaryOprator = "AND"
	OpOR                = "OR"
	OpNOT               = "NOT"
)

type ComparisonOperator string

const (
	OpLt ComparisonOperator = "Lt"
	OpEq                    = "Eq"
	OpGt                    = "Gt"
)

type BinaryOperation struct {
	Op          BinaryOprator
	Left, Right Node
}

type ComparisonOperation struct {
	Op    ComparisonOperator
	Key   string
	Value interface{}
}

func (n BinaryOperation) Eval(p map[string]interface{}) bool {
	switch n.Op {
	case OpAND:
		return n.Left.Eval(p) && n.Right.Eval(p)
	case OpOR:
		return n.Left.Eval(p) || n.Right.Eval(p)
	case OpNOT:
		return !n.Left.Eval(p)
	default:
		panic("Unknown binary operation " + n.Op)
	}
}

func (n BinaryOperation) Serialize() ([]byte, error) {
	return serialize(n)
}

func (n ComparisonOperation) EvalImpl(p map[string]interface{}, cmpInt func(int, int) bool, cmpString func(string, string) bool) bool {
	if val, found := p[n.Key]; found {
		switch valT := val.(type) {
		case int:
			switch nvalT := n.Value.(type) {
			case int:
				return cmpInt(valT, nvalT)
			default:
				panic(PANIC_UNMATCHED_TYPE)
			}
		case string:
			switch nvalT := n.Value.(type) {
			case string:
				return cmpString(valT, nvalT)
			default:
				panic(PANIC_UNMATCHED_TYPE)
			}
		default:
			panic(PANIC_UNMATCHED_TYPE)
		}
	}
	return false
}

func (n ComparisonOperation) EvalEq(p map[string]interface{}) bool {
	cmpInt := func(a int, b int) bool { return a == b }
	cmpString := func(a string, b string) bool { return a == b }
	return n.EvalImpl(p, cmpInt, cmpString)
}

func (n ComparisonOperation) EvalGt(p map[string]interface{}) bool {
	cmpInt := func(a int, b int) bool { return a > b }
	cmpString := func(a string, b string) bool { return a > b }
	return n.EvalImpl(p, cmpInt, cmpString)
}

func (n ComparisonOperation) EvalLt(p map[string]interface{}) bool {
	cmpInt := func(a int, b int) bool { return a < b }
	cmpString := func(a string, b string) bool { return a < b }
	return n.EvalImpl(p, cmpInt, cmpString)
}

func (n ComparisonOperation) Eval(p map[string]interface{}) bool {
	switch n.Op {
	case OpEq:
		return n.EvalEq(p)
	case OpGt:
		return n.EvalGt(p)
	case OpLt:
		return n.EvalLt(p)
	default:
		panic("Unknown comparison operation " + n.Op)

	}
}

func (n ComparisonOperation) Serialize() ([]byte, error) {
	return serialize(n)
}

func serialize(n Node) ([]byte, error) {
	registerGob()
	buf := bytes.NewBuffer(nil)
	zw := zlib.NewWriter(buf)
	defer zw.Close()
	err := gob.NewEncoder(zw).Encode(&n)
	if err != nil {
		return nil, err
	}
	zw.Flush()
	return buf.Bytes(), nil
}

func Deserialize(data []byte) (Node, error) {
	registerGob()
	var node Node
	buf := bytes.NewBuffer(data)
	zr, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	err = gob.NewDecoder(zr).Decode(&node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func registerGob() {
	gob.Register(BinaryOperation{})
	gob.Register(ComparisonOperation{})
}
