package cli

import (
	"strings"
	"testing"

	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &testhelpers.StubPlayerStore{}
	cli := &CLI{playerStore, in}
	cli.PlayPoker()
	assert.True(t, len(playerStore.WinCalls) != 0)

	got := playerStore.WinCalls[0]
	want := "Chris"
	assert.Equal(t, want, got)
}
