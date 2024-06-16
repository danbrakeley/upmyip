package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

	awscfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		return
	}

	fmt.Printf("Invoking lambda %s...\n", cfg.LambdaName)
	err = InvokeLambda(ctx, awscfg, cfg.LambdaName, cfg.User, info.IP)
	if err != nil {
		fmt.Printf("Error invoking lambda: %v\n", err)
		return
	}

	fmt.Printf("Done\n")
}

type LambdaRequest struct {
	User string `json:"user"`
	IP   string `json:"ip"`
}

func InvokeLambda(ctx context.Context, awscfg aws.Config, lambdaName, user, ip string) error {
	svc := lambda.NewFromConfig(awscfg)

	// Create the request payload
	payload, err := json.Marshal(LambdaRequest{
		User: user,
		IP:   ip,
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	// Invoke the Lambda function
	result, err := svc.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName: aws.String(lambdaName),
		Payload:      payload,
	})
	if err != nil {
		return err
	}

	// Check for function error
	if result.FunctionError != nil {
		return fmt.Errorf("lambda function error: %s", *result.FunctionError)
	}

	return nil
}
