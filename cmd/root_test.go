package cmd

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootOptionsApplyLogFlags(t *testing.T) {
	testCases := []struct {
		name     string
		verbose  bool
		debug    bool
		expected slog.Level
	}{
		{name: "default", expected: slog.LevelWarn},
		{name: "verbose", verbose: true, expected: slog.LevelInfo},
		{name: "debug", verbose: true, debug: true, expected: slog.LevelDebug},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := newRootOptions()
			o.verbose = tc.verbose
			o.debug = tc.debug

			o.applyLogFlags()

			assert.Equal(t, tc.expected, o.logLevel.Level())
		})
	}
}
