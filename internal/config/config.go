// Package config provides an interface to the .trpc.conf config file
package config

import (
	"fmt"
	"os/user"
	"path"

	"github.com/pelletier/go-toml"
)

// Config contains all the external configuration data from .trpc.conf.
type Config struct {
	Trackernames map[string]string
	Settings     *toml.Tree
}

// ReadConfig attempts to read ~/.trpc.conf as a toml file and returns a config tree.
func ReadConfig() *Config {
	var TomlConfig *toml.Tree

	c := &Config{
		Trackernames: make(map[string]string),
	}

	usr, err := user.Current()
	if err != nil {
		return nil
	}

	TomlConfig, err = toml.LoadFile(path.Join(usr.HomeDir, ".trpc.conf"))
	if err != nil {
		return c
	}

	trackernames := TomlConfig.Get("trackernames")

	tnames := &toml.Tree{}

	if trackernames != nil {
		tnames = trackernames.(*toml.Tree)
	}

	for _, shortname := range tnames.Keys() {
		trackers := tnames.Get(shortname)
		switch v := trackers.(type) {
		case string:
			c.Trackernames[trackers.(string)] = shortname
		case []interface{}:
			for _, tracker := range trackers.([]interface{}) {
				c.Trackernames[tracker.(string)] = shortname
			}
		default:
			fmt.Printf("Unknown %T\n", v)
		}
	}

	settings := TomlConfig.Get("settings")
	if settings != nil {
		c.Settings = settings.(*toml.Tree)
	} else {
		x := toml.Tree{}
		c.Settings = &x
	}
	return c
}
