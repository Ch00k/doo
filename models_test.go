package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStruct struct {
	CompletedAt CompletedAt
}

func TestMarshalCompletedAtNull(t *testing.T) {
	s := TestStruct{}
	v, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "{\"CompletedAt\":null}", string(v))
}

func TestMarshalCompletedAtNotNull(t *testing.T) {
	s := TestStruct{CompletedAt: CompletedAt{Int64: 123456789, Valid: true}}
	v, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "{\"CompletedAt\":123456789}", string(v))
}

func TestUnmarshalCompletedAtNull(t *testing.T) {
	var s TestStruct
	json.Unmarshal([]byte("{\"CompletedAt\":null}"), &s)
	assert.Equal(t, TestStruct{}, s)
}

func TestUnmarshalCompletedAtnotNull(t *testing.T) {
	var s TestStruct
	json.Unmarshal([]byte("{\"CompletedAt\":123456789}"), &s)
	assert.Equal(t, TestStruct{CompletedAt: CompletedAt{Int64: 123456789, Valid: true}}, s)
}
