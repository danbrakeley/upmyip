# upmyip

## Overview

This is a solution to maintaining an IP allow list on a service using Security Groups in AWS. It gives a cli to end users that they can run to first look up their public-facing IP address, and to then send it along to a lambda that will revoke any old IP address for that user, and authorize the new one.

The lambda and CLI are both written in Go, and deployment and adding users happens via the included CloudFormation scripts (see aws/README.md). The CLI reads its credentials from its own config file (and not from `~/.aws`, this was done to keep it simple to roll out to users).

## WARNING

If you want to keep something secure, then put it behind a VPN. This solution is just meant to reduce the attack surface, but doesn't offer any real protection.

I TAKE NO RESPONSIBILITY FOR THIS WORKING OR NOT WORKING. ASSUME IT DOESN'T WORK, AND THEN PROVE TO YOURSELF THAT IT DOES BEFORE USING.

If you find something wrong/broken, please let me know and/or open a PR to help fix it!

## Dev Setup

- Install a recent version of [Go](https://go.dev/dl/)
- Install [Mage](https://magefile.org/)
  - There's a helper bash script for installing and upgrading Mage here: `./scripts/reinstall-mage.sh`.
- run `mage` to see the build targets, e.g.
  ```text
  $ mage
  Targets:
    all            tests, builds, and packages all targets
    ci             runs all CI tasks
    lambda         tests, builds, and packages the lambda
    lambdaBuild    builds the lamda (output goes to "local" folder)
    lambdaZip      zips the lambda
    test           tests all packages
    upMyIP         tests and builds the upmyip cli app
    upMyIPBuild    builds the upmyip cli app
    upMyIPRun      runs the upmyip cli app in the "local" folder
  ```

Building and packaging happen in the `local` folder, which is ignored by git.

### Building the Lambda

- run `mage lambda`
  - output is `local/lambda.zip`

Note that deploying code changes to all running lambdas can be automated via some bash script in [aws/README.md](aws/README.md).

### Building the CLI

- run `mage build`
  - output is `local/upmyip[.exe]`
- it will require a `upmyip.toml` config file in the current folder, in the form:
  ```toml
  lambda = "LAMBDA_FUNCTION_NAME"
  access_key = "ACCESS_KEY"
  secret_key = "SECRET"
  ```

## AWS Setup

See [aws/README.md](aws/README.md).

### Adding users

- For each new user, deploy the `per-user.yaml` CF template in the `aws` folder (see the `README` in that folder for more specifics).
- Create a `upmyip.toml` for this user by hand (see `Building the CLI` above for an example).
- Securely send the user the config file.
- Send the user the latest cli executable.
