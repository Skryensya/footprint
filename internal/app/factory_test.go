package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	// StyleEnabled should be true by default
	require.True(t, opts.StyleEnabled)
}

func TestNewForTesting(t *testing.T) {
	app := NewForTesting()

	require.NotNil(t, app.Git)
	require.NotNil(t, app.Repo)
	require.NotNil(t, app.Config)
	require.NotNil(t, app.Logger)
	require.NotNil(t, app.Output)
	require.NotNil(t, app.Styler)
	require.NotNil(t, app.Hooks)
}

func TestClose_NilComponents(t *testing.T) {
	app := NewForTesting()
	// Store is intentionally nil in test app
	app.Store = nil

	// Should not panic
	err := Close(app)
	require.NoError(t, err)
}
