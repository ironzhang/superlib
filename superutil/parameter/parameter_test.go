package parameter

import (
	"testing"

	"github.com/ironzhang/superlib/fileutil"
)

func TestParameter(t *testing.T) {
	t.Logf("%+v", Param)
}

func TestWriteParameter(t *testing.T) {
	fileutil.WriteTOML("supername.conf", Param)
}
