//go:build mage

package main

import (
	"bytes"
	"strings"
	"time"

	"github.com/danbrakeley/bsh"
	"github.com/magefile/mage/mg"
)

var sh = &bsh.Bsh{}

// Test tests all packages
func Test() {
	sh.Echo("Running unit tests...")
	sh.Cmd("go test ./...").Run()
}

// LambdaBuild builds the lamda (output goes to "local" folder)
func LambdaBuild() {
	sh.MkdirAll("local/")

	sh.Echof("Building lambda...")
	sh.Cmd(
		"go build -o local/bootstrap -tags lambda.norpc ./cmd/lambda/",
	).Env("GOOS=linux", "GOARCH=arm64", "CGO_ENABLED=0").Run()
}

// LambdaZip zips the lambda
func LambdaZip() {
	sh.Echof("Zipping executable to local/lambda.zip...")
	sh.ZipFile("local/bootstrap", "local/lambda.zip")
}

// Lambda tests, builds, and packages the lambda
func Lambda() {
	mg.SerialDeps(Test, LambdaBuild, LambdaZip)
}

// UpMyIPBuild builds the upmyip cli app
func UpMyIPBuild() {
	exeName := sh.ExeName("upmyip")
	sh.Echof("Building local/%s...", exeName)
	sh.MkdirAll("local/")

	// grab git commit hash to use as version for local builds
	commit := "(dev)"
	var b bytes.Buffer
	n := sh.Cmd(`git log --pretty=format:'%h' -n 1`).Out(&b).Err(&b).RunExitStatus()
	if n == 0 {
		commit = strings.TrimSpace(b.String())
	}

	sh.Cmdf(
		`go build -ldflags '`+
			`-X "github.com/danbrakeley/upmyip/internal/buildvar.Version=%s" `+
			`-X "github.com/danbrakeley/upmyip/internal/buildvar.BuildTime=%s" `+
			`-X "github.com/danbrakeley/upmyip/internal/buildvar.ReleaseURL=https://github.com/danbrakeley/upmyip"`+
			`' -o local/%s ./cmd/upmyip/`, commit, time.Now().Format(time.RFC3339), exeName,
	).Run()
}

// UpMyIPRun runs the upmyip cli app in the "local" folder
func UpMyIPRun() {
	exeName := sh.ExeName("upmyip")
	sh.Echo("Running...")
	sh.Cmdf("./%s", exeName).Dir("local").Run()
}

// UpMyIP tests and builds the upmyip cli app
func UpMyIP() {
	mg.SerialDeps(Test, UpMyIPBuild)
}

// All tests, builds, and packages all targets
func All() {
	mg.SerialDeps(Test, Lambda, UpMyIPBuild)
}

// CI runs all CI tasks
func CI() {
	mg.SerialDeps(Test, LambdaBuild, UpMyIPBuild)
}
