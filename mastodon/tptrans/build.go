//+build ignore

// Builds and deploys the application to minipaas.
package main

import (
	"context"
	"log"
	"os"

	"within.website/x/internal"
	"within.website/x/internal/greedo"
	"within.website/x/internal/minipaas"
	"within.website/x/internal/yeet"
)

func main() {
	internal.HandleStartup()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := append(os.Environ(), []string{"CGO_ENABLED=0", "GOOS=linux"}...)
	yeet.ShouldWork(ctx, env, yeet.WD, "vgo", "build", "-o=worker")
	yeet.ShouldWork(ctx, env, yeet.WD, "appsluggr", "-worker=worker")
	fin, err := os.Open("slug.tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	fname := "sona-pi-toki-pona-" + yeet.DateTag + ".tar.gz"
	pubURL, err := greedo.CopyFile(fname, fin)
	if err != nil {
		log.Fatal(err)
	}

	err = minipaas.Exec("tar:from sona-pi-toki-pona " + pubURL)
	if err != nil {
		log.Fatal(err)
	}
}
