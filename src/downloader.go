// Responsible for downloading large binaries.
package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Downloads a binary of url to dest using mirror.
// If overwrite is false and the binary exist, it
// will return its (abs) path without downloading it.
// It first downloads it into filename.part and then
// moves it to filename. This is to avoid overwriting
// already existing binaries with half-downloaded ones
// or treating half-downloaded ones as full for future
// runs.
func download_binary(url string, dest string, mirror Mirror, overwrite bool) string {
	name := filepath.Base(url)
	path := filepath.Join(dest, name)
	part := path + ".part"
	safe_url := use_mirror(mirror, url)

	if !overwrite {
		if _, err := os.Stat(path); err == nil {
			log.Print("\"" + path + "\" already exists. Skipping.")

			abs_path, err := filepath.Abs(path)
			check(err)

			return abs_path
		}
	}

	log.Print("Started downloading " + safe_url)

	out, err := os.Create(part)
	check(err)

	defer out.Close()

	resp, err := http.Get(url)
	check(err)

	if resp.StatusCode > 400 {
		log.Fatalf("Downloading "+safe_url+" failed with status code: %d\n", resp.StatusCode)
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	check(err)

	err = os.Rename(part, path)
	check(err)

	log.Print("Finished downloading " + safe_url)

	abs_path, err := filepath.Abs(path)
	check(err)

	return abs_path
}
