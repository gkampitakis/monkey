package object_test

import (
	"testing"

	"github.com/gkampitakis/monkey/object"
	"github.com/stretchr/testify/require"
)

func TestStringHashingKey(t *testing.T) {
	hello1 := &object.String{Value: "Hello World"}
	hello2 := &object.String{Value: "Hello World"}

	diff1 := &object.String{Value: "My name is johnny"}
	diff2 := &object.String{Value: "My name is johnny"}

	require.Equal(t, hello1.HashKey(), hello2.HashKey())
	require.Equal(t, diff1.HashKey(), diff2.HashKey())
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &object.Boolean{Value: true}
	true2 := &object.Boolean{Value: true}
	false1 := &object.Boolean{Value: false}
	false2 := &object.Boolean{Value: false}

	require.Equal(t, true1.HashKey(), true2.HashKey())
	require.Equal(t, false1.HashKey(), false2.HashKey())
	require.NotEqual(t, true1.HashKey(), false1.HashKey())
}

func TestIntegerHashKey(t *testing.T) {
	one1 := &object.Integer{Value: 1}
	one2 := &object.Integer{Value: 1}
	two1 := &object.Integer{Value: 2}
	two2 := &object.Integer{Value: 2}

	require.Equal(t, one1.HashKey(), one2.HashKey())
	require.Equal(t, two1.HashKey(), two2.HashKey())
	require.NotEqual(t, one1.HashKey(), two1.HashKey())
}
