package bet

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnd(t *testing.T) {
	tree := BinaryOperation{
		Op: OpAND,
		Left: ComparisonOperation{
			Op:    OpEq,
			Key:   "val1",
			Value: "1",
		},
		Right: ComparisonOperation{
			Op:    OpEq,
			Key:   "val2",
			Value: "2",
		},
	}

	serialized, err := tree.Serialize()
	assert.Equal(t, nil, err)
	tree2, err := Deserialize(serialized)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{
		"val1": "1",
		"val2": "2",
	}
	assert.Equal(t, true, tree2.Eval(p))

	p["val2"] = "3"
	assert.Equal(t, false, tree2.Eval(p))
}

func TestOr(t *testing.T) {
	tree := BinaryOperation{
		Op: OpOR,
		Left: ComparisonOperation{
			Op:    OpEq,
			Key:   "val1",
			Value: "1",
		},
		Right: ComparisonOperation{
			Op:    OpEq,
			Key:   "val2",
			Value: "2",
		},
	}

	serialized, err := tree.Serialize()
	assert.Equal(t, nil, err)
	tree2, err := Deserialize(serialized)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{
		"val1": "2",
		"val2": "2",
	}
	assert.Equal(t, true, tree2.Eval(p))

	p["val2"] = "3"
	assert.Equal(t, false, tree2.Eval(p))
}

func TestNot(t *testing.T) {
	tree := BinaryOperation{
		Op: OpNOT,
		Left: ComparisonOperation{
			Op:    OpEq,
			Key:   "val1",
			Value: "1",
		},
	}

	serialized, err := tree.Serialize()
	assert.Equal(t, nil, err)
	tree2, err := Deserialize(serialized)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{
		"val1": "2",
	}
	assert.Equal(t, true, tree2.Eval(p))

	p["val1"] = "1"
	assert.Equal(t, false, tree2.Eval(p))
}

func TestEq(t *testing.T) {
	tree := ComparisonOperation{
		Op:    OpEq,
		Key:   "val1",
		Value: "1",
	}

	serialized, err := tree.Serialize()
	assert.Equal(t, nil, err)
	tree2, err := Deserialize(serialized)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{}
	assert.Equal(t, false, tree2.Eval(p))

	p["val1"] = "1"
	assert.Equal(t, true, tree2.Eval(p))

	p["val1"] = "2"
	assert.Equal(t, false, tree2.Eval(p))
}

func TestEq2(t *testing.T) {
	tree := ComparisonOperation{
		Op:    OpEq,
		Key:   "val1",
		Value: 1,
	}

	serialized, err := tree.Serialize()
	assert.Equal(t, nil, err)
	tree2, err := Deserialize(serialized)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{}
	assert.Equal(t, false, tree2.Eval(p))

	p["val1"] = 1
	assert.Equal(t, true, tree2.Eval(p))

	defer func() {
		err := recover()
		if err != PANIC_UNMATCHED_TYPE {
			t.Errorf("panic %v", err)
		}
	}()
	p["val1"] = "1"
	assert.Equal(t, false, tree2.Eval(p))
}

func TestLt(t *testing.T) {
	tree := ComparisonOperation{
		Op:    OpLt,
		Key:   "val1",
		Value: "1",
	}

	serialized, err := tree.Serialize()
	assert.Equal(t, nil, err)
	tree2, err := Deserialize(serialized)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{}
	assert.Equal(t, false, tree2.Eval(p))

	p["val1"] = "0"
	assert.Equal(t, true, tree2.Eval(p))

	p["val1"] = "2"
	assert.Equal(t, false, tree2.Eval(p))
}

func TestLt2(t *testing.T) {
	tree := ComparisonOperation{
		Op:    OpLt,
		Key:   "val1",
		Value: 1,
	}

	serialized, err := tree.Serialize()
	assert.Equal(t, nil, err)
	tree2, err := Deserialize(serialized)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{}
	assert.Equal(t, false, tree2.Eval(p))

	p["val1"] = 0
	assert.Equal(t, true, tree2.Eval(p))

	defer func() {
		err := recover()
		if err != PANIC_UNMATCHED_TYPE {
			t.Errorf("panic %s", err)
		}
	}()
	p["val1"] = "2"
	assert.Equal(t, false, tree.Eval(p))
}

func TestGt(t *testing.T) {
	tree := ComparisonOperation{
		Op:    OpGt,
		Key:   "val1",
		Value: "1",
	}

	serialized, err := tree.Serialize()
	assert.Equal(t, nil, err)
	tree2, err := Deserialize(serialized)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{}
	assert.Equal(t, false, tree2.Eval(p))

	p["val1"] = "2"
	assert.Equal(t, true, tree2.Eval(p))

	p["val1"] = "0"
	assert.Equal(t, false, tree2.Eval(p))
}

func TestGt2(t *testing.T) {
	tree := ComparisonOperation{
		Op:    OpGt,
		Key:   "val1",
		Value: 1,
	}

	serialized, err := tree.Serialize()
	assert.Equal(t, nil, err)
	tree2, err := Deserialize(serialized)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{}
	assert.Equal(t, false, tree2.Eval(p))

	p["val1"] = 2
	assert.Equal(t, true, tree2.Eval(p))

	defer func() {
		err := recover()
		if err != PANIC_UNMATCHED_TYPE {
			t.Errorf("panic %s", err)
		}
	}()
	p["val1"] = "0"
	assert.Equal(t, false, tree2.Eval(p))
}

func TestSerializeDeserialize(t *testing.T) {
	tree := BinaryOperation{
		Op: OpAND,
		Left: ComparisonOperation{
			Op:    OpEq,
			Key:   "val1",
			Value: "1",
		},
		Right: ComparisonOperation{
			Op:    OpEq,
			Key:   "val2",
			Value: "2",
		},
	}

	serialized, err := tree.Serialize()
	assert.Equal(t, nil, err)
	//t.Logf("%v", base64.StdEncoding.EncodeToString(serialized))
	tree2, err := Deserialize(serialized)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{
		"val1": "1",
		"val2": "2",
	}
	assert.Equal(t, true, tree2.Eval(p))

	p["val2"] = "3"
	assert.Equal(t, false, tree2.Eval(p))
}

func TestDeserialize(t *testing.T) {
	encodedTree := "eJzKFGDQT88syShN0kvOz9XPS8zPzrS0NDTUd3IN0U3Pz0nMS9dzysxLLKr0L0gtSizJzM/738jMyMiPJsj4v4mBkZmRyb+AkYeBkcUnNa2EUYCBkTUoMz0DxGJg+P/kf1MJI7OjnwujMSEbnfNzCxKLMovz8xC2NjMzMgpjkWD834JkM7N3aiWIZg1LzClNBVuc879FgpHJtZCRpSwxx5CRrbikKDMvnYeZgdGQgTynIJlnhGyeEQMDAAAA//8="
	dec, _ := base64.StdEncoding.DecodeString(encodedTree)
	tree2, err := Deserialize(dec)
	assert.Equal(t, nil, err)

	p := map[string]interface{}{
		"val1": "1",
		"val2": "2",
	}
	assert.Equal(t, true, tree2.Eval(p))

	p["val2"] = "3"
	assert.Equal(t, false, tree2.Eval(p))
}
