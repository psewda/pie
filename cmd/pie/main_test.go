package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Parallel()

	go func() {
		defer func() { q <- os.Interrupt }()
		time.Sleep(time.Second)

		url := fmt.Sprintf("http://%s:%d", "127.0.0.1", 8800)
		cmd := exec.Command("/usr/bin/curl", "-i", url)
		out, err := cmd.Output()

		assert.Nil(t, err)
		assert.Contains(t, string(out), "200 OK")
		assert.Contains(t, string(out), "Hello Pie !!")
	}()
	main()
}
