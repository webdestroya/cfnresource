package cfnerr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdestroya/cfnresource/cfnerr"
)

func TestErrorCreation(t *testing.T) {
	innerErr := errors.New("inner error fake")

	t.Run("As", func(t *testing.T) {
		cerr, ok := cfnerr.As(innerErr)
		require.False(t, ok)
		require.Nil(t, cerr)
	})

	t.Run("New_blank", func(t *testing.T) {
		err := cfnerr.New("", "bad req", nil)
		require.ErrorContains(t, err, "bad req")
		require.Equal(t, cfnerr.InternalFailure, err.Code())
	})

	t.Run("NewMessage", func(t *testing.T) {
		err := cfnerr.NewMessage(cfnerr.InvalidRequest, "bad req")
		require.ErrorContains(t, err, "bad req")
		require.Equal(t, cfnerr.InvalidRequest, err.Code())
		require.Nil(t, err.OrigErr())
		require.Equal(t, "bad req", err.Message())
		require.Implements(t, (*fmt.Stringer)(nil), err)
	})

	t.Run("NewMessage_blank", func(t *testing.T) {
		err := cfnerr.NewMessage(cfnerr.InvalidRequest, "")
		require.ErrorContains(t, err, "InvalidRequest")
		require.Equal(t, cfnerr.InvalidRequest, err.Code())
		require.NotErrorIs(t, err, innerErr)
	})

	t.Run("Wrap", func(t *testing.T) {
		err := cfnerr.Wrap(cfnerr.InvalidRequest, innerErr)
		require.Error(t, err)
		require.ErrorIs(t, err, innerErr)
		require.ErrorContains(t, err, "InvalidRequest")
		require.ErrorContains(t, err, "inner error fake")
		require.Equal(t, cfnerr.InvalidRequest, err.Code())
		require.True(t, cfnerr.Is(err))
		inErr := errors.Unwrap(err)
		require.ErrorIs(t, inErr, innerErr)

		cerr, ok := cfnerr.As(err)
		require.True(t, ok)
		require.NotNil(t, cerr)

	})
}
