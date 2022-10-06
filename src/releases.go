// Responsible for downloading & parsing downloads.json.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Link struct {
	Binary string
	Sig    string
}

type Releases struct {
	Tag       string
	Version   string
	Downloads map[string]map[string]Link
}

// Decodes io.Reader into Releases.
func parse_releases(releases io.Reader) Releases {
	var res Releases

	err := json.NewDecoder(releases).Decode(&res)
	check(err)

	return res
}

// Renews downloads.json/releases.json.
func renew_releases(dest string) Releases {
	resp, err := http.Get(API)
	check(err)

	if resp.StatusCode != 200 {
		log.Fatal("\"downloads.json\" responded with: " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	file, err := os.Create(filepath.Join(dest, "releases.json"))
	check(err)

	defer file.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	io.Copy(file, tee)

	log.Print("Renewed releases.json")

	return parse_releases(&buf)
}

// [Unused] Checks if a platform+locale exist in Releases.
// func available_release(releases Releases, platform string, locale string) bool {
// 	return releases.Downloads[platform][locale].Sig != ""
// }

// Whether it should renew.
// True only if releases.json doesn't exit
// or a day has passed since it got last
// renewd.
func should_renew(data_path string) bool {
	if file_info, err := os.Stat(filepath.Join(data_path, "releases.json")); err == nil {
		last_mod_plus_1_day := file_info.ModTime().AddDate(0, 0, 1)
		if last_mod_plus_1_day.UnixMilli() < time.Now().UnixMilli() {
			return true
		} else {
			return false
		}
	} else if errors.Is(err, os.ErrNotExist) {
		return true
	} else {
		check(err)
	}

	return true
}

// Helper function that checks
// if it should renew and returns
// Releases.
func get_releases(dest string) Releases {
	if should_renew(dest) {
		return renew_releases(dest)
	}

	file, err := os.Open(filepath.Join(dest, "releases.json"))
	check(err)

	return parse_releases(file)
}
