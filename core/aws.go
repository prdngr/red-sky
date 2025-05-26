package core

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

func InitializeAwsSession(profile string) {
	StartSpinner("Initializing AWS session")

	session, err := session.NewSessionWithOptions(session.Options{Profile: profile})

	if err != nil {
		log.Fatalf("error initializing AWS session: %s", err)
	}

	stsClient := sts.New(session)
	identity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		log.Fatalf("error retrieving AWS account details: %s", err)
	}

	var accountAlias = "n/a"
	iamClient := iam.New(session)
	aliasOutput, err := iamClient.ListAccountAliases(&iam.ListAccountAliasesInput{})

	if err != nil {
		log.Fatalf("error retrieving AWS account alias: %s", err)
	} else if len(aliasOutput.AccountAliases) > 0 {
		accountAlias = *aliasOutput.AccountAliases[0]
	}

	StopSpinner("AWS session initialized")

	PrintHeader("AWS Account Information")

	fmt.Printf("▶ Account: %s (%s)\n", *identity.Account, accountAlias)
	fmt.Printf("▶ Caller ARN: %s\n", *identity.Arn)
	fmt.Printf("▶ User ID: %s\n", *identity.UserId)

	fmt.Print("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println()
}
