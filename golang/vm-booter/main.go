package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v8"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

type VirtualMachineInfo struct {
	SubscriptionId string `json:"subscription_id"`
	ResourceGroup  string `json:"resource_group"`
}

var vmInfo map[string]VirtualMachineInfo

const vmInfoFile = "vm_info.json"

func checkError(error error, message string) {
	if error != nil {
		fmt.Printf("%s: %v\n", message, error)
		os.Exit(1)
	}
}

func loadVirtualMachineInfo() error {
	_, err := os.Stat(vmInfoFile)
	if errors.Is(err, os.ErrNotExist) {
		err = generateVirtualMachineInfo()
		if err != nil {
			return err
		}

		err = saveVirtualMachineInfo()
		if err != nil {
			return err
		}
	} else {
		data, err := os.ReadFile(vmInfoFile)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(data, &vmInfo); err != nil {
			return err
		}
	}

	return nil
}

func saveVirtualMachineInfo() error {
	data, err := json.MarshalIndent(vmInfo, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(vmInfoFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func generateVirtualMachineInfo() error {
	vmInfo = map[string]VirtualMachineInfo{}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	subClient, err := armsubscriptions.NewClient(cred, nil)
	if err != nil {
		return err
	}

	subPager := subClient.NewListPager(nil)
	ctx := context.Background()
	for subPager.More() {
		subResp, err := subPager.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, sub := range subResp.Value {
			vmClient, err := armcompute.NewVirtualMachinesClient(*sub.SubscriptionID, cred, nil)
			if err != nil {
				return err
			}

			vmPager := vmClient.NewListAllPager(nil)
			for vmPager.More() {
				vmResp, err := vmPager.NextPage(ctx)
				if err != nil {
					return err
				}

				for _, vm := range vmResp.Value {
					parsedId, err := arm.ParseResourceID(*vm.ID)
					if err != nil {
						return err
					}

					vmInfo[*vm.Name] = VirtualMachineInfo{
						SubscriptionId: *sub.SubscriptionID,
						ResourceGroup:  parsedId.ResourceGroupName,
					}
				}
			}
		}
	}

	return nil
}

func main() {
	err := loadVirtualMachineInfo()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(vmInfo)
}
