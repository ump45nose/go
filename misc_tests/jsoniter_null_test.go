package misc_tests

import (
	"bytes"
	"io"
	"testing"

	"github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func Test_read_null(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigDefault, `null`)
	should.True(iter.ReadNil())
	iter = jsoniter.ParseString(jsoniter.ConfigDefault, `null`)
	should.Nil(iter.Read())
	iter = jsoniter.ParseString(jsoniter.ConfigDefault, `navy`)
	iter.Read()
	should.True(iter.Error != nil && iter.Error != io.EOF)
	iter = jsoniter.ParseString(jsoniter.ConfigDefault, `navy`)
	iter.ReadNil()
	should.True(iter.Error != nil && iter.Error != io.EOF)
}

func Test_write_null(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := jsoniter.NewStream(jsoniter.ConfigDefault, buf, 4096)
	stream.WriteNil()
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("null", buf.String())
}

func Test_decode_null_object_field(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigDefault, `[null,"a"]`)
	iter.ReadArray()
	if iter.ReadObject() != "" {
		t.FailNow()
	}
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
	type TestObject struct {
		Field string
	}
	objs := []TestObject{}
	should.Nil(jsoniter.UnmarshalFromString("[null]", &objs))
	should.Len(objs, 1)
}

func Test_decode_null_array_element(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigDefault, `[null,"a"]`)
	should.True(iter.ReadArray())
	should.True(iter.ReadNil())
	should.True(iter.ReadArray())
	should.Equal("a", iter.ReadString())
}

func Test_decode_null_string(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigDefault, `[null,"a"]`)
	should.True(iter.ReadArray())
	should.Equal("", iter.ReadString())
	should.True(iter.ReadArray())
	should.Equal("a", iter.ReadString())
}

func Test_decode_null_skip(t *testing.T) {
	iter := jsoniter.ParseString(jsoniter.ConfigDefault, `[null,"a"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}

func Test_invalid_null_does_not_nil_preinitialized_map(t *testing.T) {
	should := require.New(t)
	for _, input := range []string{"none", "nan", "nil", "notjson", "n"} {
		data := map[string]interface{}{}
		err := jsoniter.Unmarshal([]byte(input), &data)
		should.Error(err, input)
		should.NotNil(data, input)
		data["a"] = 1
	}
	data := map[string]interface{}{"keep": true}
	should.Nil(jsoniter.Unmarshal([]byte("null"), &data))
	should.Nil(data)
}

func Test_invalid_null_does_not_nil_preinitialized_slice(t *testing.T) {
	should := require.New(t)
	data := []interface{}{1}
	err := jsoniter.Unmarshal([]byte("none"), &data)
	should.Error(err)
	should.NotNil(data)
	should.Equal([]interface{}{1}, data)
}
