package fs

import (
	"os"
	"os/user"
	"path/filepath"
)

// DefaultTokenFile is deprecated, and will be only used for migration.
const DefaultTokenFile = "credentials"

// DefaultConfigsFile stores cli credentials and hosts.
const DefaultConfigsFile = "configs"

// InfluxDir retrieves the influxdbv2 directory.
func InfluxDir() (string, error) {
	var dir string
	// By default, store meta and data files in current users home directory
	u, err := user.Current()
	if err == nil {
		dir = u.HomeDir
	} else if home := os.Getenv("HOME"); home != "" {
		dir = home
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		dir = wd
	}
	dir = filepath.Join(dir, ".influxdbv2")

	return dir, nil
}
