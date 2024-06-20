package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func main() {
	_, noColor := os.LookupEnv("NO_COLOR")
	prn := Printer{
		NoColor: noColor,
	}

	prn.Header("UpMyIP v0.1.2")

	cfg, err := LoadConfig("upmyip.toml")
	if err != nil {
		prn.Error("upmyip.toml", err)
		return
	}

	ctx := context.Background()

	s := NewSpinner(noColor)
	prn.Print("Finding public IP address... ")
	s.Start()
	info, err := RequestPublicInfo(ctx)
	if err != nil {
		s.Stop()
		prn.Error("Error", err)
		return
	}
	s.Stop()
	prn.BrightIPln(info)

	credProvider := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
	)

	awscfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credProvider),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		prn.Error("Error loading AWS config", err)
		return
	}

	prn.Print("Updating IP... ")
	s.Start()
	err = InvokeLambda(ctx, awscfg, cfg.LambdaName, info.IP)
	if err != nil {
		s.Stop()
		prn.Error("AWS Error", err)
		return
	}
	s.Stop()
	prn.BrightPrintln("Done")

	prn.Print("Success!\n\n")
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
		return fmt.Errorf("remote function error: %s", *result.FunctionError)
	}

	return nil
}
