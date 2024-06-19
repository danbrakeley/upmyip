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

	s := NewSpinner()
	fmt.Print("Finding public IP address... ")
	s.Start()
	info, err := RequestPublicInfo(ctx)
	if err != nil {
		s.Stop()
		fmt.Printf("%sError: %s%v%s\n", SGR(FgRed), SGR(FgYellow), err, SGR(FgReset))
		return
	}
	s.Stop()
	fmt.Printf("%s%s %s(%s)%s\n", SGR(FgWhite), info.IP, SGR(FgCyan), info.ISP, SGR(FgReset))

	credProvider := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
	)

	awscfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credProvider),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		fmt.Printf("%sError loading AWS config: %s%v%s\n", SGR(FgRed), SGR(FgYellow), err, SGR(FgReset))
		return
	}

	fmt.Print("Updating IP... ")
	s.Start()
	err = InvokeLambda(ctx, awscfg, cfg.LambdaName, info.IP)
	if err != nil {
		s.Stop()
		fmt.Printf("%sError invoking lambda: %s%v%s\n", SGR(FgRed), SGR(FgYellow), err, SGR(FgReset))
		return
	}
	s.Stop()
	fmt.Printf("%sDone%s\n", SGR(FgWhite), SGR(FgReset))

	fmt.Printf("Success!\n\n")
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
