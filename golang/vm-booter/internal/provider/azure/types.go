package azure

type AzureVirtualMachineManager struct {
	cacheFile          string
	virtualMachineInfo map[string]AzureVirtualMachineInfo
}

type AzureVirtualMachineInfo struct {
	SubscriptionId string `json:"subscription_id"`
	ResourceGroup  string `json:"resource_group"`
}
