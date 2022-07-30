package appctx

import (
	"embed"

	"github.com/silverswords/devcloud/pkg/config"
	"github.com/spf13/viper"
)

// all the config init here.
var (
	GlobalAppContext *AppContext = &AppContext{

		Viper: config.GlobalViper,
	}
)

type AppContext struct {
	// support the embedFS
	embedFS map[string]map[string]*embed.FS

	app appInfo

	*viper.Viper
}

type appInfo struct {
	Name string
}
