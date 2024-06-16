# upmyip

## Overview

If you've got a service you want people to access remotely, but you don't want to expose to the entire internet (e.g. Perforce), then a quick solution is to get everyone's public IP address and add them to the security group protecting that service.

Of course, few people's IP is truly static, and ideally a user would be able to update their IP using some secure credentials (like an AWS account).

This project aims to provide some structure and a CLI to allow your users to update their current IP as they see fit.

WARNING: If you care about security, setup a VPN. Remember that a single IP can hide an unlimited number of actual devices. Did one of your users just open up your precious service to an entire Starbucks? An entire hotel? An entire university? USE WITH CAUTION!

## Pieces

- IAM Policy for lambda: aws/lambda-policy.json
- IAM Role for lambda: includes AWSLambdaBasicExecutionRole and the above policy
- Lambda
- IAM Policy to execute lambda
- IAM User for each actual user: includes above policy
- CLI app to find public IP address and invoke lambda
