package bet

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
)

const PANIC_UNMATCHED_TYPE = "unmatched_type"

type NodeEncoded interface {
	Decode() Node
}
type Node interface {
	Eval(p map[string]interface{}) bool
	Encode() NodeEncoded
	Serialize() ([]byte, error)
}

type And struct {
	Left, Right Node
}

type AndEnc struct {
	Left, Right NodeEncoded
}

func (n And) Eval(p map[string]interface{}) bool {
	return n.Left.Eval(p) && n.Right.Eval(p)
}

func (n And) Encode() NodeEncoded {
	return AndEnc{
		Left:  n.Left.Encode(),
		Right: n.Right.Encode(),
	}
}

func (n And) Serialize() ([]byte, error) {
	return serialize(n)
}

func (n AndEnc) Decode() Node {
	return And{
		Left:  n.Left.Decode(),
		Right: n.Right.Decode(),
	}
}

type Or struct {
	Left, Right Node
}

type OrEnc struct {
	Left, Right NodeEncoded
}

func (n Or) Eval(p map[string]interface{}) bool {
	return n.Left.Eval(p) || n.Right.Eval(p)
}

func (n Or) Encode() NodeEncoded {
	return OrEnc{
		Left:  n.Left.Encode(),
		Right: n.Right.Encode(),
	}
}

func (n Or) Serialize() ([]byte, error) {
	return serialize(n)
}

func (n OrEnc) Decode() Node {
	return Or{
		Left:  n.Left.Decode(),
		Right: n.Right.Decode(),
	}
}

type Not struct {
	Child Node
}

type NotEnc struct {
	Child NodeEncoded
}

func (n Not) Eval(p map[string]interface{}) bool {
	return !n.Child.Eval(p)
}

func (n Not) Encode() NodeEncoded {
	return NotEnc{
		Child: n.Child.Encode(),
	}
}

func (n Not) Serialize() ([]byte, error) {
	return serialize(n)
}

func (n NotEnc) Decode() Node {
	return Not{
		Child: n.Child.Decode(),
	}
}

type NodeValue struct {
	Key   string
	Value interface{}
}

type NodeValueInt struct {
	Key   string
	Value int
}

type NodeValueString struct {
	Key   string
	Value string
}

type NodeValueEncoded interface {
	Decode() NodeValue
}

func (n NodeValue) Eval(p map[string]interface{}, cmpInt func(int, int) bool, cmpString func(string, string) bool) bool {
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

func (n NodeValue) Encode() NodeValueEncoded {
	switch valT := n.Value.(type) {
	case int:
		return NodeValueInt{
			Key:   n.Key,
			Value: valT,
		}
	case string:
		return NodeValueString{
			Key:   n.Key,
			Value: valT,
		}
	default:
		panic(PANIC_UNMATCHED_TYPE)
	}
}

func (n NodeValueInt) Decode() NodeValue {
	return NodeValue{
		Key:   n.Key,
		Value: n.Value,
	}
}

func (n NodeValueString) Decode() NodeValue {
	return NodeValue{
		Key:   n.Key,
		Value: n.Value,
	}
}

type Eq struct {
	NodeValue
}

type EqString struct {
	NodeValueString
}

type EqInt struct {
	NodeValueInt
}

func (n Eq) Eval(p map[string]interface{}) bool {
	cmpInt := func(a int, b int) bool { return a == b }
	cmpString := func(a string, b string) bool { return a == b }
	return n.NodeValue.Eval(p, cmpInt, cmpString)
}

func (n Eq) Encode() NodeEncoded {
	switch valT := n.NodeValue.Encode().(type) {
	case NodeValueInt:
		return EqInt{
			valT,
		}
	case NodeValueString:
		return EqString{
			valT,
		}
	default:
		panic(PANIC_UNMATCHED_TYPE)
	}
}

func (n Eq) Serialize() ([]byte, error) {
	return serialize(n)
}

func (n EqString) Decode() Node {
	return Eq{
		n.NodeValueString.Decode(),
	}
}

func (n EqInt) Decode() Node {
	return Eq{
		n.NodeValueInt.Decode(),
	}
}

type Gt struct {
	NodeValue
}

type GtString struct {
	NodeValueString
}

type GtInt struct {
	NodeValueInt
}

func (n Gt) Eval(p map[string]interface{}) bool {
	cmpInt := func(a int, b int) bool { return a > b }
	cmpString := func(a string, b string) bool { return a > b }
	return n.NodeValue.Eval(p, cmpInt, cmpString)
}

func (n Gt) Encode() NodeEncoded {
	switch valT := n.NodeValue.Encode().(type) {
	case NodeValueInt:
		return GtInt{
			valT,
		}
	case NodeValueString:
		return GtString{
			valT,
		}
	default:
		panic(PANIC_UNMATCHED_TYPE)
	}
}

func (n Gt) Serialize() ([]byte, error) {
	return serialize(n)
}

func (n GtString) Decode() Node {
	return Gt{
		n.NodeValueString.Decode(),
	}
}

func (n GtInt) Decode() Node {
	return Gt{
		n.NodeValueInt.Decode(),
	}
}

type Lt struct {
	NodeValue
}

type LtString struct {
	NodeValueString
}

type LtInt struct {
	NodeValueInt
}

func (n Lt) Eval(p map[string]interface{}) bool {
	cmpInt := func(a int, b int) bool { return a < b }
	cmpString := func(a string, b string) bool { return a < b }
	return n.NodeValue.Eval(p, cmpInt, cmpString)
}

func (n Lt) Encode() NodeEncoded {
	switch valT := n.NodeValue.Encode().(type) {
	case NodeValueInt:
		return LtInt{
			valT,
		}
	case NodeValueString:
		return LtString{
			valT,
		}
	default:
		panic(PANIC_UNMATCHED_TYPE)
	}
}

func (n Lt) Serialize() ([]byte, error) {
	return serialize(n)
}

func (n LtString) Decode() Node {
	return Lt{
		n.NodeValueString.Decode(),
	}
}

func (n LtInt) Decode() Node {
	return Lt{
		n.NodeValueInt.Decode(),
	}
}

func serialize(n Node) ([]byte, error) {
	registerGob()
	buf := bytes.NewBuffer(nil)
	zw := zlib.NewWriter(buf)
	defer zw.Close()
	nodeEncoded := n.Encode()
	err := gob.NewEncoder(zw).Encode(&nodeEncoded)
	if err != nil {
		return nil, err
	}
	zw.Flush()
	return buf.Bytes(), nil
}

func Deserialize(data []byte) (Node, error) {
	registerGob()
	var nodeEnc NodeEncoded
	buf := bytes.NewBuffer(data)
	zr, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	err = gob.NewDecoder(zr).Decode(&nodeEnc)
	if err != nil {
		return nil, err
	}

	return nodeEnc.Decode(), nil
}

func registerGob() {
	gob.Register(AndEnc{})
	gob.Register(OrEnc{})
	gob.Register(NotEnc{})
	gob.Register(EqString{})
	gob.Register(EqInt{})
	gob.Register(LtString{})
	gob.Register(LtInt{})
	gob.Register(GtString{})
	gob.Register(GtInt{})
}
