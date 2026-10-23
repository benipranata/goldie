package goldie_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/benipranata/goldie/pkg/goldie"
)

//go:generate go test -v -update -clean .

func TestEqual_Integration(t *testing.T) {
	type user struct {
		Name  string   `json:"name"`
		Roles []string `json:"roles"`
	}
	userJSON, err := json.MarshalIndent(user{Name: "Beni", Roles: []string{"admin", "dev"}}, "", "  ")
	require.NoError(t, err)

	tests := []struct {
		name   string
		actual []byte
	}{
		{
			name: "text",
			actual: []byte(`Hello, Gopher!
This is a golden file.
`),
		},
		{name: "json", actual: userJSON},
		{name: "unicode", actual: []byte("héllo, 世界 🌏\n")},
		{name: "binary", actual: []byte{0x00, 0x01, 0xfe, 0xff}},
		{name: "empty", actual: []byte{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			goldie.New(t).Equal(t.Name(), tc.actual)
		})
	}
}
