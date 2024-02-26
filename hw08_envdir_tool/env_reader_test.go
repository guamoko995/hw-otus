package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("complex", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.Nil(t, err)
		require.Equal(t, Environment{
			"BAR":   {Value: "bar"},
			"EMPTY": {Value: ""},
			"FOO":   {Value: "   foo\nwith new line"},
			"HELLO": {Value: "\"hello\""},
			"UNSET": {NeedRemove: true},
		}, env)
	})
}

func TestEnvRead(t *testing.T) {
	t.Run("without changes", func(t *testing.T) {
		actual := listToMap(EnvRead(Environment{}))
		expected := listToMap(os.Environ())
		equal(t, expected, actual)
	})

	t.Run("Add", func(t *testing.T) {
		actual := listToMap(
			EnvRead(
				Environment{
					"HELLO": EnvValue{
						Value: "hello",
					},
				},
			),
		)

		expected := listToMap(os.Environ())
		expected["HELLO"] = "hello"

		equal(t, expected, actual)
	})

	t.Run("change", func(t *testing.T) {
		expected := listToMap(os.Environ())
		var rundomRealEnv string
		for rundomRealEnv = range expected {
			break
		}
		expected[rundomRealEnv] = "\"hello\""

		actual := listToMap(
			EnvRead(
				Environment{
					rundomRealEnv: EnvValue{
						Value: "\"hello\"",
					},
				},
			),
		)

		equal(t, expected, actual)
	})

	t.Run("remove", func(t *testing.T) {
		expected := listToMap(os.Environ())
		var rundomRealEnv string
		for rundomRealEnv = range expected {
			break
		}

		actual := EnvRead(Environment{
			rundomRealEnv: EnvValue{
				NeedRemove: true,
			},
		})

		for _, str := range actual {
			key := strings.SplitN(str, "=", 2)[0]
			require.NotEqual(t, key, rundomRealEnv)
		}
	})
}

// equal проверяет эквивалентность словарей.
func equal(t *testing.T, expected, actual map[string]string) {
	t.Helper()
	// Каждая актуальная переменная найдется в ожидаемом словаре
	for key, val := range actual {
		require.Equal(t, expected[key], val)
	}

	// Каждая ожидаемая переменная найдется в актуальном словаре
	for key, val := range expected {
		require.Equal(t, actual[key], val)
	}
}

// listToMap преобразует массив строк вида "key=val" в
// словарь значений val, с ключами key.
func listToMap(list []string) map[string]string {
	expEnv := make(map[string]string, len(list))
	for _, str := range list {
		keyVal := strings.SplitN(str, "=", 2)
		key, val := keyVal[0], keyVal[1]
		expEnv[key] = val
	}
	return expEnv
}
