# AWS CloudFormation Templates

## Outline

- Core: S3 Bucket to store lambda function zip
- Per User
  - Lambda function using shared zip, with an env var set to this user's name
    - tag with a common value as all such user lambdas can be updated easily
  - Policy allowing user to invoke lambda
  - Policy allowing lambda to edit a security group

## Setup

1. Deploy `core.yaml`
   ```text
   aws cloudformation deploy --template ./aws/core.yaml --stack-name upmyip-core
   ```
2. Get the new bucket's name:
   ```text
   aws cloudformation describe-stacks --stack-name upmyip-core --query "Stacks[0].Outputs[?OutputKey=='BucketName'].OutputValue" --output text
   ```
3. Build and upload lambda function
   ```text
   $ mage lambda
   Running unit tests...
   Building lambda...
   Zipping executable to local/lambda.zip...

   $ aws s3 cp ./local/lambda.zip s3://BUCKET_NAME
   upload: local\lambda.zip to s3://BUCKET_NAME/lambda.zip
   ```
4. For each user, deploy `per-user.yaml`
   ```text
   aws cloudformation deploy --template ./aws/per-user.yaml --stack-name upmyip-user-USERNAME \
     --capabilities CAPABILITY_IAM \
     --parameter UserName=USERNAME BucketForLambdaZips=BUCKET_NAME SecurityGroupId=SECURITY_GROUP_ID
   ```

## Generate access key for the newly created user

```text
aws iam create-access-key --user-name INVOCATION_USER
```

## Invoke lambda

Run upmyip, or:

```text
aws lambda invoke --function-name FUNCITON_NAME_OR_ARN --payload '{"ip":"10.0.0.1"}' --cli-binary-format raw-in-base64-out output.json
```

## Update lambda when code changes

1. Build and upload lambda function
   ```text
   $ mage lambda
   Running unit tests...
   Building lambda...
   Zipping executable to local/lambda.zip...

   $ aws s3 cp ./local/lambda.zip s3://BUCKET_NAME
   upload: local\lambda.zip to s3://BUCKET_NAME/lambda.zip
   ```
2. Update all lambda functions with "Updatable"="upmyip"
   ```bash
   function_arns=$(aws lambda list-functions --query 'Functions[*].FunctionArn' --output text)

   for arn in $function_arns; do
     tags=$(aws lambda list-tags --resource $arn --query 'Tags' --output json)

     if echo $tags | jq -e --arg key "UpdatableBy" --arg value "upmyip" '.[$key] == $value' > /dev/null; then
       aws lambda update-function-code --function-name $arn --s3-bucket $S3BUCKETNAME --s3-key lambda.zip
     fi
   done
   ```
