package config

import (
	"github.com/spf13/viper"
)

var GlobalViper *viper.Viper

// todo(abserari): make it happen once.
func init() {
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
	vp.ReadInConfig()

	vp.SetEnvPrefix("dc") // will be uppercased automatically
	vp.AutomaticEnv()

	vp.WatchConfig()

	GlobalViper = vp
}
