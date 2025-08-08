package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	v2iam "github.com/aws/aws-sdk-go-v2/service/iam"
	v2sts "github.com/aws/aws-sdk-go-v2/service/sts"
)

func InitializeAwsSession(profile string) {
	StartSpinner("Initializing AWS session")

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile))
	if err != nil {
		log.Fatalf("error loading AWS config: %s", err)
	}

	stsClient := v2sts.NewFromConfig(cfg)
	identity, err := stsClient.GetCallerIdentity(context.Background(), &v2sts.GetCallerIdentityInput{})

	if err != nil {
		log.Fatalf("error retrieving AWS account details: %s", err)
	}

	var accountAlias = "n/a"
	iamClient := v2iam.NewFromConfig(cfg)
	paginator := v2iam.NewListAccountAliasesPaginator(iamClient, &v2iam.ListAccountAliasesInput{})

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
