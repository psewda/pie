package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/psewda/pie/app"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	port := app.GetRandPort()
	os.Setenv("PIE_PORT", strconv.Itoa(int(port)))

	go func() {
		defer func() { q <- os.Interrupt }()
		time.Sleep(time.Second)

		url := fmt.Sprintf("http://%s:%d/api/v1/pie-store/version", "127.0.0.1", port)
		cmd := exec.Command("/usr/bin/curl", "-i", url)
		out, err := cmd.Output()

		assert.Nil(t, err)
		assert.Contains(t, string(out), "200 OK")
		assert.Contains(t, string(out), "version")
	}()
	main()
}

func TestParsePort(t *testing.T) {
	// test happy path
	port, ok := parsePort("8800")
	assert.True(t, ok)
	assert.Equal(t, uint16(8800), port)

	// test invalid inputs
	data := []string{"", "invalid", "1"}
	for _, i := range data {
		if port, ok := parsePort(i); true {
			assert.False(t, ok)
			assert.Equal(t, uint16(0), port)
		}
	}
}
