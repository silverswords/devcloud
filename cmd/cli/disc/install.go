package disc

import (
	"github.com/silverswords/devcloud/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	runtimeVersion string
)
var InitCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Disc",
	PreRun: func(cmd *cobra.Command, args []string) {
	},
	Example: `
# Install disc
disc install xxx
`,
	Run: func(cmd *cobra.Command, args []string) {

		dir := defaultDiscBinPath()
		_ = utils.InstallBinary("latest", "disc", "silverswords", dir)
	},
}

func init() {
	defaultRuntimeVersion := "latest"
	_ = viper.BindEnv("runtime_version_override", "DISC_RUNTIME_VERSION")
	runtimeVersionEnv := viper.GetString("runtime_version_override")
	if runtimeVersionEnv != "" {
		defaultRuntimeVersion = runtimeVersionEnv
	}
	InitCmd.Flags().StringVarP(&runtimeVersion, "runtime-version", "", defaultRuntimeVersion, "The version of the Disc runtime to install, for example: 1.0.0")
	InitCmd.Flags().BoolP("help", "h", false, "Print this help message")
	rootCmd.AddCommand(InitCmd)
}
