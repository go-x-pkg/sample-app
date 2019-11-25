package server

import (
	"github.com/go-x-pkg/log"
	"github.com/go-x-pkg/sample-app/appversion"
)

const (
	company string = "©ACME"
	logo           = `
 ______  ______  ______   __   __  ______  ______
/_____/\/_____/\/_____/\ /_/\ /_/\/_____/\/_____/\
\::::_\/\::::_\/\:::_ \ \\:\ \\ \ \::::_\/\:::_ \ \
 \:\/___/\:\/___/\:(_) ) )\:\ \\ \ \:\/___/\:(_) ) )_
  \_::._\:\::___\/\: __  \ \:\_/.:\ \::___\/\: __  \ \
    /____\:\:\____/\ \  \ \ \ ..::/ /\:\____/\ \  \ \ \
    \_____\/\_____\/\_\/ \_\/\___/_(  \_____\/\_\/ \_\/
`
)

func (a *App) doVersion() {
	log.LogStd(log.Info, logo)
	log.LogStd(log.Info, company)
	log.LogfStd(log.Info, "v%s, build-at %s",
		appversion.Version, appversion.BuildDate)
}
