package disc

import (
	"os"
	"path/filepath"
)

const (
	defaultDiscDirName   = ".disc"
	defaultDisBinDirName = "bin"
)

// some code common as functions
func defaultDiscBinPath() string {
	homedir, _ := os.Hostname()

	return filepath.Join(homedir + defaultDiscDirName)
}
