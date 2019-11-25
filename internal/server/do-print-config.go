package server

import (
	"bytes"

	"github.com/go-x-pkg/log"
)

func (a *App) doPrintConfig() {
	a.ctx.cfg().DumpToFn(func(buf *bytes.Buffer) {
		log.LogfStd(log.Info,
			"using configuration:\n%s", buf.String())
	})
}
