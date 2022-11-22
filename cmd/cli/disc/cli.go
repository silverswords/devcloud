package disc

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "disc",
	Short: "Devcloud CLI",
	Long:  `your dev infrastructure self cloud`,
}

func Execute() {
	cobra.OnInitialize(func() {
		viper.SetEnvPrefix("disc")
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.AutomaticEnv()
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
