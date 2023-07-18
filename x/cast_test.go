package x

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS2I64(t *testing.T) {
	cases := []struct {
		in  string
		exp int64
	}{
		{"0", 0},
		{"999", 999},
		{"-999", -999},
	}
	for _, c := range cases {
		res := S2I64(c.in)
		log.Printf("%+v\n"+"res--------------->", res)
		assert.Equal(t, c.exp, res)
	}
}

func TestI2S(t *testing.T) {
	type cs[T Int] struct {
		in  T
		exp string
	}
	c1 := cs[int]{0, "0"}
	res := I2S(c1.in)
	log.Printf("res--------------->"+"%+v\n", res)
	assert.Equal(t, c1.exp, res)

	c2 := cs[int8]{9, "9"}
	res = I2S(c2.in)
	log.Printf("res--------------->"+"%+v\n", res)
	assert.Equal(t, c2.exp, res)

	c3 := cs[int32]{999, "999"}
	res = I2S(c3.in)
	log.Printf("res--------------->"+"%+v\n", res)
	assert.Equal(t, c3.exp, res)

	type MyI32 int32
	c4 := cs[MyI32]{999, "999"}
	res = I2S(c4.in)
	log.Printf("res--------------->"+"%+v\n", res)
	assert.Equal(t, c4.exp, res)

	c5 := cs[int16]{9999, "9999"}
	res = I2S(c5.in)
	log.Printf("res--------------->"+"%+v\n", res)
	assert.Equal(t, c5.exp, res)
}

func TestUI2S(t *testing.T) {
	type cs[T Uint] struct {
		in  T
		exp string
	}
	c1 := cs[uint8]{0, "0"}
	res := UI2S(c1.in)
	log.Printf("res--------------->"+"%+v\n", res)
	assert.Equal(t, c1.exp, res)

	c2 := cs[uint32]{9, "9"}
	res = UI2S(c2.in)
	log.Printf("res--------------->"+"%+v\n", res)
	assert.Equal(t, c2.exp, res)

	c3 := cs[uint64]{999, "999"}
	res = UI2S(c3.in)
	log.Printf("res--------------->"+"%+v\n", res)
	assert.Equal(t, c3.exp, res)
}
