package bet

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

const PANIC_UNMATCHED_TYPE = "unmatched_type"

type Node interface {
	Eval(p map[string]interface{}) bool
	Serialize() ([]byte, error)
}

type BinaryOperator string

const (
	OpAND BinaryOperator = "&&"
	OpOR  BinaryOperator = "||"
	OpNOT BinaryOperator = "!"
)

type ComparisonOperator string

const (
	OpLt    ComparisonOperator = "<"
	OpEq    ComparisonOperator = "=="
	OpGt    ComparisonOperator = ">"
	OpNotEq ComparisonOperator = "!="
)

type BinaryOperation struct {
	Op          BinaryOperator
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

func (n ComparisonOperation) EvalNotEq(p map[string]interface{}) bool {
	cmpInt := func(a int, b int) bool { return a != b }
	cmpString := func(a string, b string) bool { return a != b }
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
	case OpNotEq:
		return n.EvalNotEq(p)
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

func convertExprToBET(e ast.Expr) Node {
	switch val_e := e.(type) {
	case *ast.BinaryExpr:
		return convertBinaryExprToBET(val_e)
	case *ast.ParenExpr:
		return convertExprToBET(val_e.X)
	case *ast.UnaryExpr:
		if val_e.Op.String() != string(OpNOT) {
			panic("Unexpected Op " + val_e.Op.String())
		}
		return BinaryOperation{
			Op:   OpNOT,
			Left: convertExprToBET(val_e.X),
		}
	}
	return nil
}

func convertBinaryExprToBET(a *ast.BinaryExpr) Node {
	switch a.Op.String() {
	case string(OpLt), string(OpGt), string(OpEq), string(OpNotEq):
		node := &ComparisonOperation{
			Op: ComparisonOperator(a.Op.String()),
		}
		switch val_x := (a.X).(type) {
		case *ast.Ident:
			node.Key = val_x.Name
		default:
			panic("Unexpected type " + reflect.TypeOf(a.X).String())
		}
		switch val_y := (a.Y).(type) {
		case *ast.Ident:
			node.Value = val_y.Name
		case *ast.BasicLit:
			if val_y.Kind == token.INT {
				var err error
				node.Value, err = strconv.Atoi(val_y.Value)
				if err != nil {
					panic(err.Error())
				}
			} else if val_y.Kind == token.STRING {
				node.Value = strings.Replace(val_y.Value, `"`, "", -1)
			} else {
				panic("Unexpected BasicLit Kind " + val_y.Kind.String())
			}
		default:
			panic("Unexpected type " + reflect.TypeOf(a.X).String())
		}

		return node
	case string(OpAND), string(OpOR):
		node := BinaryOperation{
			Op: BinaryOperator(a.Op.String()),
		}
		node.Left = convertExprToBET(a.X)
		node.Right = convertExprToBET(a.Y)
		return node
	default:
		panic("Unexpected Op " + a.Op.String())
	}
}

func ParseExpr(exprStr string) (Node, error) {
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		return nil, err
	}
	//ast.Print(nil, expr)
	return convertExprToBET(expr), nil
}
