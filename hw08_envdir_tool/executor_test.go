package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Empty command and environment", func(t *testing.T) {
		returnCode := RunCmd([]string{}, Environment{})
		require.Equal(t, -1, returnCode)
	})

	t.Run("Valid command and empty environment", func(t *testing.T) {
		returnCode := RunCmd([]string{"echo", "Hello, World!"}, Environment{})
		require.Equal(t, 0, returnCode)
	})

	t.Run("Valid command and non-empty environment", func(t *testing.T) {
		env := Environment{
			"PATH": {
				Value:      "/usr/local/bin:/usr/bin:/bin",
				NeedRemove: false,
			},
			"HOME": {
				Value:      "/home/user",
				NeedRemove: false,
			},
		}
		returnCode := RunCmd([]string{"ls", "-l"}, env)
		require.Equal(t, 0, returnCode)
	})

	// Test case 4: Test with command that fails
	t.Run("Command that fails", func(t *testing.T) {
		returnCode := RunCmd([]string{"invalid-command"}, Environment{})
		require.Equal(t, -1, returnCode)
	})
}
