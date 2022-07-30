package gin

import (
	"sync"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"github.com/silverswords/devcloud/pkg/appctx"
	"github.com/silverswords/devcloud/pkg/framework"
)

func init() {
	var config Config
	if data := appctx.GlobalAppContext.Viper.Get("http.gin"); data == nil {
		return
	} else {
		mapstructure.Decode(data, &config)
	}
	copier.Copy(&exampleConfig, &config)

	if config.Enabled {
		framework.RegisterPlugin("gin", GetGinCtx())
	}
}

var exampleConfig = Config{
	Enabled:  true,
	HostPort: ":8080",
}
var exampleMap = structs.Map(exampleConfig)
var once sync.Once
var GinCtx = &ginCtx{}

type Config struct {
	Enabled  bool
	HostPort string
}
type ginCtx struct {
	Config

	*gin.Engine
}

func GetGinCtx() *ginCtx {
	once.Do(func() {
		GinCtx.Config = exampleConfig
		GinCtx.Engine = gin.New()
	})
	return GinCtx
}

func (c *ginCtx) Start() {
	c.Run(c.HostPort)
}
