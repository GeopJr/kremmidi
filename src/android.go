// Responsible for getting Android builds' version
// and generating Android builds' links.
package main

import (
	"net/http"

	"gopkg.in/ini.v1"
)

type Android_Build struct {
	binary string
	sig    string
}

// Gets the Android builds' version from
// the versions.ini. If it fails it will
// return the fallback.
func get_android_version(fallback string) string {
	android_version := fallback

	resp, err := http.Get(VERSIONS_INI)
	check(err)

	if resp.StatusCode == 200 {
		cfg, err := ini.Load(resp.Body)
		check(err)

		defer resp.Body.Close()

		android_version = cfg.Section("torbrowser-android-stable").Key("version").String()
		if android_version == "" {
			android_version = fallback
		}
	}

	return android_version
}

// Returns a map of the Android builds
// based on the android_version.
func get_android_builds(android_version string) map[string]Android_Build {
	res := make(map[string]Android_Build)

	archs := []string{"aarch64", "armv7", "x86_64", "x86"}
	for _, arch := range archs {
		res[arch] = Android_Build{
			binary: "https://dist.torproject.org/torbrowser/" + android_version + "/tor-browser-" + android_version + "-android-" + arch + "-multi.apk",
			sig:    "https://dist.torproject.org/torbrowser/" + android_version + "/tor-browser-" + android_version + "-android-" + arch + "-multi.apk.asc",
		}
	}

	return res
}
