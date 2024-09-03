package xhttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetURLHost(t *testing.T) {
	host1 := "http://minio-service:9000"
	host1Str, err := GetURLHost(host1)
	if err != nil {
		t.Errorf("get host port error: %v ", err)
	}
	assert.Equal(t, "minio-service:9000", host1Str)

	host2 := "https://shimo.im"
	host2Str, err := GetURLHost(host2)
	if err != nil {
		t.Errorf("get host port error: %v ", err)
	}
	assert.Equal(t, "shimo.im", host2Str)
}
