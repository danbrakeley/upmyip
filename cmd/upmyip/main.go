package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func main() {
	cfg, err := LoadConfig("upmyip.toml")
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()

	fmt.Printf("Finding public IP address...\n")
	info, err := RequestPublicInfo(ctx)
	if err != nil {
		fmt.Printf("Error requesting public IP: %v\n", err)
		return
	}

	credProvider := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
	)

	awscfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credProvider),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	fmt.Printf("Invoking lambda %s...\n", cfg.LambdaName)
	err = InvokeLambda(ctx, awscfg, cfg.LambdaName, info.IP)
	if err != nil {
		fmt.Printf("Error invoking lambda: %v\n", err)
		return
	}

	fmt.Printf("Done\n")
}

func InvokeLambda(ctx context.Context, awscfg aws.Config, lambdaName, ip string) error {
	svc := lambda.NewFromConfig(awscfg)

	payload := fmt.Sprintf(`{"ip": "%s"}`, ip)
	result, err := svc.Invoke(ctx, &lambda.InvokeInput{
		FunctionName: aws.String(lambdaName),
		Payload:      []byte(payload),
	})
	if err != nil {
		return err
	}
	if result.FunctionError != nil {
		return fmt.Errorf("lambda function error: %s", *result.FunctionError)
	}

	return nil
}
