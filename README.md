# upmyip

## Overview

If you've got a service you want people to access remotely, but you don't want to expose to the entire internet (e.g. Perforce), then a quick solution is to get everyone's public IP address and add them to the security group protecting that service.

Of course, few people's IP is truly static, and ideally a user would be able to update their IP using some secure credentials (like an AWS account).

This project aims to provide some structure and a CLI to allow your users to update their current IP as they see fit.

WARNING: If you care about security, setup a VPN. Remember that a single IP can hide an unlimited number of actual devices. Did one of your users just open up your precious service to an entire Starbucks? An entire hotel? An entire university? USE WITH CAUTION!

## Pieces

- AWS
  - Each user gets a lambda function that is hard coded to that user.
  - Each user gets an IAM User that only has permission to invoke their lambda.
  - Each lambda gets permission to alter a single security group.
- User
  - Each user gets the `upmyip` executable (which is the same for all users).
  - Each user gets a `upmyip.toml` that is specific to them, and includes their lambda's name, and the access key and secret for their user.

## Getting a new user setup

- For each new user, deploy the `per-user` CF template in the `aws` folder (see the `README` in that folder for more specifics).

- Build the upmyip.toml for this user
  ```text
  lambda = LAMBDA_FUNCTION_NAME
  access_key = ACCESS_KEY
  secret_key = SECRET
  ```

- Zip up the `upmyip.exe` and the new `upmyip.toml` and send them to the user
  - You can always send them an updated upmyip.exe that they can use with their existing toml file.

- Note that you don't have to add a security group ingress rule for them ahead of time, it will get added automatically.
