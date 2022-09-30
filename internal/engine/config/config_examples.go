package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Example quick start example for a common use case of this module.
func Example() {
	// A common use case is to get a private config folder for your app to
	// place its settings files into, that are specific to the local user.
	configPath := LocalConfig("my-app")
	err := makePath(configPath) // Ensure it exists.
	if err != nil {
		log.Fatal(err)
	}

	// Deal with a JSON configuration file in that folder.
	configFile := filepath.Join(configPath, "settings.json")
	type AppSettings struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var settings AppSettings

	// Does the file not exist?
	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		// Create the new config file.
		settings = AppSettings{"MyUser", "MyPassword"}
		fh, err := os.Create(configFile)
		if err != nil {
			log.Fatal(err)
		}
		defer func(fh *os.File) {
			err := fh.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(fh)

		encoder := json.NewEncoder(fh)
		err = encoder.Encode(&settings)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Load the existing file.
		fh, err := os.Open(configFile)
		if err != nil {
			log.Fatal(err)
		}
		defer func(fh *os.File) {
			err := fh.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(fh)

		decoder := json.NewDecoder(fh)
		err = decoder.Decode(&settings)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// ExampleSystemConfig example for getting a global system configuration path.
func ExampleSystemConfig() {
	// Get all the global system configuration paths.
	//
	// On Linux or BSD this might be []string{"/etc/xdg"} or the split
	// version of $XDG_CONFIG_DIRS.
	//
	// On macOS or Windows this will likely return a slice with only one entry
	// that points to the global config path; see the README.md for details.
	paths := SystemConfig()
	fmt.Printf("Global system config paths: %v\n", paths)

	// Or you can get a version of the path suffixed with a vendor folder.
	vendor := SystemConfig("acme")
	fmt.Printf("Vendor-specific config paths: %v\n", vendor)

	// Or you can use multiple path suffixes to group configs in a
	// `vendor/application` namespace. You can use as many path
	// components as you like.
	app := SystemConfig("acme", "sprockets")
	fmt.Printf("Vendor/app specific config paths: %v\n", app)
}

// ExampleLocalConfig example for getting a user-specific configuration path.
func ExampleLocalConfig() {
	// Get the default root of the local configuration path.
	//
	// On Linux or BSD this might be "$HOME/.config", or on Windows this might
	// be "C:\\Users\\$USER\\AppData\\Roaming"
	path := LocalConfig()
	fmt.Printf("Local user config path: %s\n", path)

	// Or you can get a local config path with a vendor suffix, like
	// "$HOME/.config/acme" on Linux.
	vendor := LocalConfig("acme")
	fmt.Printf("Vendor-specific local config path: %s\n", vendor)

	// Or you can use multiple path suffixes to group configs in a
	// `vendor/application` namespace. You can use as many path
	// components as you like.
	app := LocalConfig("acme", "sprockets")
	fmt.Printf("Vendor/app specific local config path: %s\n", app)
}

// ExampleLocalCache example for getting a user-specific cache folder.
func ExampleLocalCache() {
	// Get the default root of the local cache folder.
	//
	// On Linux or BSD this might be "$HOME/.cache", or on Windows this might
	// be "C:\\Users\\$USER\\AppData\\Local"
	path := LocalCache()
	fmt.Printf("Local user cache path: %s\n", path)

	// Or you can get a local cache path with a vendor suffix, like
	// "$HOME/.cache/acme" on Linux.
	vendor := LocalCache("acme")
	fmt.Printf("Vendor-specific local cache path: %s\n", vendor)

	// Or you can use multiple path suffixes to group caches in a
	// `vendor/application` namespace. You can use as many path
	// components as you like.
	app := LocalCache("acme", "sprockets")
	fmt.Printf("Vendor/app specific local cache path: %s\n", app)
}

// ExampleMakePath example for automatically creating config directories.
func ExampleMakePath() {
	// The MakePath() function can accept the output from any of the folder
	// getting functions and ensure that their path exists.

	// Create a local user configuration folder under an app prefix.
	// On Linux this may result in `$HOME/.config/my-cool-app` existing as
	// a directory, depending on the value of `$XDG_CONFIG_HOME`.
	err := makePath(LocalConfig("my-cool-app"))
	if err != nil {
		log.Fatal(err)
	}

	// Create a cache folder under a namespace.
	err = makePath(LocalCache("acme", "sprockets", "client"))
	if err != nil {
		log.Fatal(err)
	}

	// In the case of global system configuration, which may return more than
	// one path (especially on Linux/BSD that uses the XDG Base Directory Spec),
	// it will attempt to create the directories only under the *first* path.
	//
	// For example, if $XDG_CONFIG_DIRS="/etc/xdg:/opt/config" this will try
	// to create the config dir only in "/etc/xdg/acme/sprockets" and not try
	// to create any folders under "/opt/config".
	err = makePath(SystemConfig("acme", "sprockets")...)
	if err != nil {
		log.Fatal(err)
	}
}

// ExampleRefresh example for recalculating what the directories should be.
func ExampleRefresh() {
	// On your program's initialization, this module decides which paths to
	// use for global, local and cache folders, based on environment variables
	// and falling back on defaults.
	//
	// In case the environment variables change throughout the life of your
	// program, for example if you re-assigned $XDG_CONFIG_HOME, you can call
	// the Refresh() function to re-calculate the paths to reflect the new
	// environment.
	Refresh()
}
