package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

func checkError(error error, message string) {
	if error != nil {
		fmt.Printf("%s: %v\n", message, error)
		os.Exit(1)
	}
}

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	checkError(err, "Failed to create Azure credential")

	subClient, err := armsubscriptions.NewClient(cred, nil)
	checkError(err, "Failed to create Azure subscription client")

	subPager := subClient.NewListPager(nil)
	ctx := context.Background()
	for subPager.More() {
		subResp, err := subPager.NextPage(ctx)
		checkError(err, "Failed to list subscriptions")

		for _, sub := range subResp.Value {
			fmt.Println(*sub.DisplayName, *sub.SubscriptionID)
		}
	}
}
