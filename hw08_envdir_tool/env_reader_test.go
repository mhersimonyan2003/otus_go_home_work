package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Not found directory", func(t *testing.T) {
		dir := "testdata/not_found"
		env, err := ReadDir(dir)

		require.Nil(t, env)
		require.Error(t, err)
	})

	t.Run("Directory with files", func(t *testing.T) {
		dir := "testdata/env"
		env, err := ReadDir(dir)

		require.NoError(t, err)
		require.NotNil(t, env)

		expectedEnv := Environment{
			"BAR": EnvValue{
				Value:      "bar",
				NeedRemove: false,
			},
			"EMPTY": EnvValue{
				Value:      "",
				NeedRemove: false,
			},
			"FOO": EnvValue{
				Value:      "   foo\nwith new line",
				NeedRemove: false,
			},
			"HELLO": EnvValue{
				Value:      "\"hello\"",
				NeedRemove: false,
			},
			"UNSET": EnvValue{
				Value:      "",
				NeedRemove: true,
			},
		}

		for key, expectedValue := range expectedEnv {
			require.Contains(t, env, key)
			require.Equal(t, expectedValue, env[key])
		}
	})
}
