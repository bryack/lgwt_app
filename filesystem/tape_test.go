package filesystem

import (
	"io"
	"testing"

	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestTape_Write(t *testing.T) {
	file, clean := testhelpers.CreateTempFile(t, "12345")
	defer clean()

	tape := &tape{file: file}

	tape.Write([]byte("abc"))

	file.Seek(0, io.SeekStart)
	newFileContent, _ := io.ReadAll(file)

	got := string(newFileContent)
	want := "abc"
	assert.Equal(t, want, got)
}
