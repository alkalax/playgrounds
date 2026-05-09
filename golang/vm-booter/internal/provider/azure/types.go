package azure

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

type AzureVirtualMachineManager struct {
	credential         *azidentity.DefaultAzureCredential
	cacheFile          string
	virtualMachineInfo map[string]AzureVirtualMachineInfo
}

type AzureVirtualMachineInfo struct {
	SubscriptionId string `json:"subscription_id"`
	ResourceGroup  string `json:"resource_group"`
}
