package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func InitializeAwsSession(profile string) {
	StartSpinner("Initializing AWS session")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("error loading AWS config: %s", err)
	}

	stsClient := sts.NewFromConfig(cfg)
	identity, err := stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})

	if err != nil {
		log.Fatalf("error retrieving AWS account details: %s", err)
	}

	var accountAlias = "n/a"
	iamClient := iam.NewFromConfig(cfg)
	paginator := iam.NewListAccountAliasesPaginator(iamClient, &iam.ListAccountAliasesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			log.Fatalf("error retrieving AWS account alias: %s", err)
		}

		if len(page.AccountAliases) > 0 {
			accountAlias = page.AccountAliases[0]
			break
		}
	}

	StopSpinner("AWS session initialized")

	PrintHeader("AWS Account Information")

	fmt.Printf("▶ Account: %s (%s)\n", aws.ToString(identity.Account), accountAlias)
	fmt.Printf("▶ Caller ARN: %s\n", aws.ToString(identity.Arn))
	fmt.Printf("▶ User ID: %s\n", aws.ToString(identity.UserId))

	fmt.Print("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println()
}
