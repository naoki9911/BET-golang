package bet

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnd(t *testing.T) {
	tree := And{
		Eq{
			NodeValue{
				Key:   "val1",
				Value: "1",
			},
		},
		Eq{
			NodeValue{
				Key:   "val2",
				Value: "2",
			},
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
	tree := Or{
		Eq{
			NodeValue{
				Key:   "val1",
				Value: "1",
			},
		},
		Eq{
			NodeValue{
				Key:   "val2",
				Value: "2",
			},
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
	tree := Not{
		Eq{
			NodeValue{
				Key:   "val1",
				Value: "1",
			},
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
	tree := Eq{
		NodeValue{
			Key:   "val1",
			Value: "1",
		},
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
	tree := Eq{
		NodeValue{
			Key:   "val1",
			Value: 1,
		},
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
	tree := Lt{
		NodeValue{
			Key:   "val1",
			Value: "1",
		},
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
	tree := Lt{
		NodeValue{
			Key:   "val1",
			Value: 1,
		},
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
	tree := Gt{
		NodeValue{
			Key:   "val1",
			Value: "1",
		},
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
	tree := Gt{
		NodeValue{
			Key:   "val1",
			Value: 1,
		},
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
	tree := And{
		Eq{
			NodeValue{
				Key:   "val1",
				Value: "1",
			},
		},
		Eq{
			NodeValue{
				Key:   "val2",
				Value: "2",
			},
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
	encodedTree := "eJyUjs9Kh0AUhc/5jUYLIZ8i2qSMO5cFLqKI6N9+0mkcshkKDVpm1gP2Qjcq3AXR7rv34xzOWY5d58d+uinaeF8GE+98XWtdHjaX+y4OJrjiIHRNaOVFkVs/TJnBDZMTezsyB9Nz7/ovAuRD5ivu/dXZPFyMjz44eVXk9npRFpDcOY2dvTbDZNf3O4BS3tTvjhuqY/vMDEy/FTMAR7JkZPJkBk1q4B+j1mBFVgA+AQAA//8="
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
