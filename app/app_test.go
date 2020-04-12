package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppRun(t *testing.T) {
	port := GetRandPort()
	app := NewApp()

	// test invalid port #
	if err := app.Run(1); true {
		assert.Error(t, err)
	}

	// test app run - happy path
	if err := app.Run(port); true {
		assert.Nil(t, err)

		// duplicate app runs
		if err := app.Run(port); true {
			assert.Nil(t, err)
		}
	}
	app.Dispose()
}
