package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mozilla-services/product-delivery-tools/post_upload/post_upload"
)

func main() {
	app := cli.NewApp()
	app.Name = "post_upload"
	app.HideVersion = true
	app.Version = Version
	app.Usage = ""
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jeremy Orem",
			Email: "oremj@mozilla.com",
		},
	}
	app.Action = doMain
	app.Flags = Flags

	app.Run(os.Args)
}

func doMain(c *cli.Context) {
	errs := make([]error, 0)
	requireArgs := func(args ...string) (hasErrors bool) {
		for _, arg := range args {
			if c.String(arg) == "" {
				hasErrors = true
				errs = append(errs, fmt.Errorf("--%s must be set", arg))
			}
		}
		return
	}

	boolRequireArgs := func(boolArg string, args ...string) bool {
		if c.Bool(boolArg) {
			return requireArgs(args...)
		}
		return false
	}

	if len(c.Args()) < 2 {
		errs = append(errs, errors.New("you must specify a directory and at least one file."))
	}

	requireArgs("product")
	boolRequireArgs("release-to-latest", "branch")
	boolRequireArgs("release-to-dated", "branch", "buildid", "nightly-dir")
	boolRequireArgs("release-to-candidates-dir", "version", "build_number")
	boolRequireArgs("release-to-mobile-candidates-dir", "version", "build-number", "builddir")
	boolRequireArgs("release-to-tinderbox-builds", "tinderbox-builds-dir")
	boolRequireArgs("release-to-dated-tinderbox-builds", "tinderbox-builds-dir", "buildid")
	boolRequireArgs("release-to-try-builds", "who", "revision", "builddir")

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println("Error:", err)
		}
		os.Exit(1)
	}

	uploadDir := c.Args()[0]
	files := c.Args()[1:]

	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			fmt.Printf("Error: %s does not exist.\n", f)
			os.Exit(1)
		}
	}

	if c.Bool("release-to-latest") {
		post_upload.ReleaseToLatest(
			c.String("branch"), c.String("tinderbox-builds-dir"), uploadDir, files)
	}
	if c.Bool("release-to-dated") {
		post_upload.ReleaseToDated(
			c.String("branch"),
			c.String("build-id"),
			c.String("product"),
			c.String("nightly-dir"),
			!c.Bool("no-shortdir"), uploadDir, files)
	}

	if c.Bool("release-to-candidates-dir") {
		post_upload.ReleaseToCandidates(
			c.String("build-dir"), c.String("build-number"),
			c.String("product"), c.String("tinderbox-builds-dir"),
			c.String("version"), c.Bool("signed"),
			uploadDir, files)
	}

	if c.Bool("release-to-mobile-candidates-dir") {
		post_upload.ReleaseToMobileCandidates(
			c.String("version"), c.String("build-number"),
			c.String("nightly-dir"), c.String("product"),
			uploadDir, files)

	}

	if c.Bool("releaset-to-tinderbox-builds") {
		post_upload.ReleaseToTinderboxBuilds(
			c.String("product"), c.String("build-id"),
			c.String("build-dir"), c.String("tinderbox-builds-dir"),
			uploadDir, files)
	}

	if c.Bool("release-to-dated-tinderbox-builds") {
		post_upload.ReleaseToDatedTinderboxBuilds(
			c.String("product"), c.String("build-id"),
			c.String("build-dir"), c.String("tinderbox-builds-dir"),
			uploadDir, files)
	}

	if c.Bool("release-to-try-builds") {
		post_upload.ReleaseToTryBuilds(
			c.String("product"), c.String("who"),
			c.String("revision"), c.String("build-dir"),
			uploadDir, files)
	}
}