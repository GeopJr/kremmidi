package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const VERSION = "0.2.1"
const API = "https://aus1.torproject.org/torbrowser/update_3/release/downloads.json"
const VERSIONS_INI = "https://gitlab.torproject.org/tpo/web/tpo/-/raw/main/databags/versions.ini"

// Calls download_binary & split and returns
// an array of Binary and an array of Part.
func index(links []string, platform string, dest string, version string, locale string, mirror Mirror, overwrite bool, should_split bool, limit int, binary_col *mongo.Collection) ([]Binary, []Part) {
	var binaries []Binary
	var parts []Part

	for _, link := range links {
		var file_parts []string
		id := primitive.NewObjectID()

		end_path := download_binary(link, dest, mirror, overwrite)
		is_sig := filepath.Ext(strings.ToLower(end_path)) == ".asc"

		if is_sig {
			continue
		}

		binary_name := filepath.Base(end_path)
		parts_no := 0

		// Check if it already exists and skip (including splitting).
		if db_doc_exists(binary_col, binary_name, platform, version, locale) {
			continue
		}

		if should_split {
			file_parts = split(end_path, filepath.Join(filepath.Dir(end_path), "parts"), int64(limit))

			parts_no = len(file_parts)

			for i, file_part := range file_parts {
				abs_path, err := filepath.Abs(file_part)
				check(err)

				doc := Part{
					Name:       binary_name,
					Version:    version,
					Path:       abs_path,
					Part_no:    i,
					Belongs_to: id,
					Platform:   platform,
				}

				parts = append(parts, doc)
			}

		}

		abs_path, err := filepath.Abs(end_path)
		check(err)

		doc := Binary{
			ID:       id,
			Name:     binary_name,
			Version:  version,
			Path:     abs_path,
			Sig:      abs_path + ".asc",
			Locale:   locale,
			Parts:    parts_no,
			Platform: platform,
		}

		binaries = append(binaries, doc)
	}

	return binaries, parts
}

func main() {
	android := true
	desktop := true
	should_split := true
	mirror := TOR_PROJECT

	// from CLI
	var db_url string
	var dest string
	var overwrite bool
	var drop_db bool
	var split_limit int
	var locales string

	// Whether it passed
	// the CLI. (aka not help)
	passed := false
	no_android := false
	no_desktop := false
	no_should_split := false

	// mirror as string from CLI
	var mirror_tmp string

	app := &cli.App{
		Name:      "kremmidi",
		Usage:     "Experimental GetTor backend with multiple frontends.",
		Version:   "v" + VERSION,
		UsageText: "kremmidi [global options] [command]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "database",
				Aliases:     []string{"d", "db"},
				Usage:       "Database url.",
				Value:       os.Getenv("KREMMIDI_DB"),
				Destination: &db_url,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o", "out", "dest"},
				Usage:       "Binary destination.",
				DefaultText: "./data",
				Value:       "data",
				Destination: &dest,
			},
			&cli.StringFlag{
				Name:        "mirror",
				Aliases:     []string{"m"},
				Usage:       "Use a mirror. Available: " + get_mirrors(),
				DefaultText: "TOR_PROJECT",
				Destination: &mirror_tmp,
			},
			&cli.IntFlag{
				Name:        "limit",
				Aliases:     []string{"l"},
				Usage:       "Amount of *MB* per part.",
				DefaultText: "5",
				Value:       5,
				Destination: &split_limit,
			},
			&cli.StringFlag{
				Name:        "locales",
				Usage:       "Locales to download seperated by a comma or 'ALL'.",
				DefaultText: "en-us",
				Value:       "en-us",
				Destination: &locales,
			},
			&cli.BoolFlag{
				Name:        "overwrite",
				Usage:       "Whether to overwrite builds if they already exist.",
				DefaultText: "false",
				Value:       false,
				Destination: &overwrite,
			},
			&cli.BoolFlag{
				Name:        "drop-db",
				Usage:       "Whether to drop all collections on start.",
				DefaultText: "false",
				Value:       false,
				Destination: &drop_db,
			},
			&cli.BoolFlag{
				Name:        "no-android",
				Usage:       "Disable grabbing Android builds.",
				Destination: &no_android,
			},
			&cli.BoolFlag{
				Name:        "no-desktop",
				Usage:       "Disable grabbing Desktop builds.",
				Destination: &no_desktop,
			},
			&cli.BoolFlag{
				Name:        "no-split",
				Usage:       "Disable splitting binaries into parts.",
				Destination: &no_should_split,
			},
		},
		Action: func(*cli.Context) error {
			// string -> Mirror
			mirror = get_mirror(mirror_tmp)
			android = !no_android
			desktop = !no_desktop
			should_split = !no_should_split

			locales = strings.ToLower(locales)

			passed = true
			return nil
		},
	}

	err := app.Run(os.Args)
	check(err)

	if !passed {
		os.Exit(0)
	}

	// Default path for binaries.
	binaries_dest := filepath.Join(dest, "browser")
	os.MkdirAll(binaries_dest, os.ModePerm)

	// Connect to db
	mongodb_context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodb_client, err := mongo.Connect(mongodb_context, options.Client().ApplyURI(db_url))
	check(err)

	defer func() {
		err := mongodb_client.Disconnect(mongodb_context)
		check(err)
	}()

	mongodb_context, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = mongodb_client.Ping(mongodb_context, readpref.Primary())
	check(err)

	// Default db name and collections.
	db := mongodb_client.Database("kremmidi")
	coll_binaries := db.Collection("binaries")
	coll_parts := db.Collection("parts")

	if drop_db {
		mongodb_context, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = db.Drop(mongodb_context)
		check(err)
	}

	releases := get_releases(dest)

	if android {
		android_version := get_android_version(releases.Version)
		android_path := filepath.Join(binaries_dest, android_version, "android")
		os.MkdirAll(android_path, os.ModePerm)

		android_builds := get_android_builds(android_version)

		// Go through all Android Builds, index
		// and save to the db.
		for platform, arch := range android_builds {
			binary_docs, part_docs := index([]string{arch.binary, arch.sig}, "android-"+platform, android_path, android_version, "multi", mirror, overwrite, should_split, split_limit, coll_binaries)

			db_insert(coll_binaries, coll_parts, binary_docs, part_docs)
		}

	}

	var locales_arr []string

	// if multiple locales,
	// create an array of string.
	if strings.Contains(locales, ",") {
		locales_arr = strings.Split(locales, ",")
	}

	if desktop {
		desktop_path := filepath.Join(binaries_dest, releases.Version, "desktop")
		os.MkdirAll(desktop_path, os.ModePerm)

		for platform := range releases.Downloads {
			for locale, link := range releases.Downloads[platform] {
				lower_locale := strings.ToLower(locale)

				// If the array of locales is not empty AND it includes lower_locale,
				// OR if lower_locale is equal to the cli locales
				// OR if the cli locales is "all"
				// index & save to the db.
				if (len(locales_arr) > 0 && includes(locales_arr, lower_locale)) || (locales == lower_locale) || (locales == "all") {
					binary_docs, part_docs := index([]string{link.Binary, link.Sig}, platform, desktop_path, releases.Version, lower_locale, mirror, overwrite, should_split, split_limit, coll_binaries)

					db_insert(coll_binaries, coll_parts, binary_docs, part_docs)
				}
			}
		}
	}
}
