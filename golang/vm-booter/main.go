package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/v3"
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
			rgClient, err := armresources.NewResourceGroupsClient(*sub.SubscriptionID, cred, nil)
			checkError(err, "Failed to create resource group client")

			rgPager := rgClient.NewListPager(nil)
			for rgPager.More() {
				rgResp, err := rgPager.NextPage(ctx)
				checkError(err, "Failed to list resource groups")

				for _, rg := range rgResp.Value {
					fmt.Printf("%s | %s (%s)\n", *sub.DisplayName, *rg.Name, *rg.Location)
				}
			}
		}
	}
}
