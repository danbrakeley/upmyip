//go:build mage

package main

import (
	"github.com/danbrakeley/bsh"
	"github.com/magefile/mage/mg"
)

var sh = &bsh.Bsh{}

func Test() {
	sh.Echo("Running unit tests...")
	sh.Cmd("go test ./...").Run()
}

// BuildLambda tests and builds the lamda (output goes to "local" folder)
func BuildLambda() {
	mg.SerialDeps(Test)

	sh.MkdirAll("local/")

	sh.Echof("Building lambda...")
	sh.Cmd(
		"go build -o local/bootstrap -tags lambda.norpc ./cmd/lambda/",
	).Env("GOOS=linux", "GOARCH=arm64", "CGO_ENABLED=0").Run()

	sh.Echof("Zipping executable to local/lambda.zip...")
	sh.ZipFile("local/bootstrap", "local/lambda.zip")
}

func Build() {
	mg.SerialDeps(Test)

	sh.MkdirAll("local/")

	exeName := sh.ExeName("upmyip")
	sh.Echof("Building local/%s...", exeName)
	sh.Cmdf("go build -o local/%s ./cmd/upmyip/", exeName).Run()
}

// Run runs unit tests, builds, and runs the app
func UpMyIP() {
	mg.SerialDeps(Test, Build)

	exeName := sh.ExeName("upmyip")

	sh.Echo("Running...")
	sh.Cmdf("./%s", exeName).Dir("local").Run()
}
