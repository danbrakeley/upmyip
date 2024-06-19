package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type RequestBody struct {
	IP string `json:"ip"`
}

var (
	reIP = regexp.MustCompile(`^(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$`)
)

func hanlderWithCfg(cfg Config) func(ctx context.Context, req *RequestBody) (*string, error) {
	return func(ctx context.Context, req *RequestBody) (*string, error) {
		if !reIP.MatchString(req.IP) {
			return nil, fmt.Errorf("ip is not a valid ip address")
		}

		awsCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("loading AWS config: %w", err)
		}

		ec2Client := ec2.NewFromConfig(awsCfg)

		describeInput := &ec2.DescribeSecurityGroupRulesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("group-id"),
					Values: []string{cfg.SecurityGroup},
				},
			},
		}
		describeOutput, err := ec2Client.DescribeSecurityGroupRules(ctx, describeInput)
		if err != nil {
			return nil, fmt.Errorf("describe security group rules: %w", err)
		}

		ruleDescription := fmt.Sprintf("upmyip:%s", cfg.UserName)

		var ruleToRevoke string
		for _, rule := range describeOutput.SecurityGroupRules {
			if rule.Description != nil && *rule.Description == ruleDescription {
				if rule.SecurityGroupRuleId == nil {
					log.Printf("error looking for security group rule ID: found matching description, but rule ID is nil (should never happen)")
					return nil, fmt.Errorf("security group rule ID is nil")
				}
				ruleToRevoke = *rule.SecurityGroupRuleId
				break
			}
		}

		if len(ruleToRevoke) > 0 {
			revokeInput := &ec2.RevokeSecurityGroupIngressInput{
				GroupId: aws.String(cfg.SecurityGroup),
				SecurityGroupRuleIds: []string{
					ruleToRevoke,
				},
			}
			_, err := ec2Client.RevokeSecurityGroupIngress(ctx, revokeInput)
			if err != nil {
				return nil, fmt.Errorf("failed to revoke sg rule %s: %w", ruleToRevoke, err)
			}
		}

		authorizeInput := &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(cfg.SecurityGroup),
			IpPermissions: []types.IpPermission{
				{
					IpProtocol: aws.String("tcp"),
					FromPort:   aws.Int32(1666),
					ToPort:     aws.Int32(1666),
					IpRanges: []types.IpRange{
						{
							CidrIp:      aws.String(req.IP + "/32"),
							Description: aws.String(ruleDescription),
						},
					},
				},
			},
		}

		_, err = ec2Client.AuthorizeSecurityGroupIngress(ctx, authorizeInput)
		if err != nil {
			return nil, fmt.Errorf("authorize new security group: %w", err)
		}

		msg := fmt.Sprintf("IP updated for user %s", cfg.UserName)
		return &msg, nil
	}
}

type Config struct {
	UserName      string
	SecurityGroup string
}

func main() {
	cfg := Config{
		UserName:      os.Getenv("UPMYIP_USERNAME"),
		SecurityGroup: os.Getenv("UPMYIP_SECURITY_GROUP"),
	}
	if len(cfg.UserName) == 0 {
		panic(fmt.Errorf("UPMYIP_USERNAME environment variable not set"))
	}
	if len(cfg.SecurityGroup) == 0 {
		panic(fmt.Errorf("UPMYIP_SECURITY_GROUP environment variable not set"))
	}
	lambda.Start(hanlderWithCfg(cfg))
}
