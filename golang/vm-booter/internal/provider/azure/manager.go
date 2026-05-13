package azure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v8"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"

	"alkalax/vm-booter/internal/provider"
)

func NewVirtualMachineManager(cacheFile string) (provider.VirtualMachineManager, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	return &AzureVirtualMachineManager{
		credential:         credential,
		cacheFile:          cacheFile,
		virtualMachineInfo: map[string]AzureVirtualMachineInfo{},
	}, nil
}

func (manager *AzureVirtualMachineManager) loadVirtualMachineInfo(noCache bool) error {
	if noCache {
		return manager.generateVirtualMachineInfo()
	}

	_, err := os.Stat(manager.cacheFile)
	if errors.Is(err, os.ErrNotExist) {
		err = manager.generateVirtualMachineInfo()
		if err != nil {
			return err
		}

		err = manager.saveVirtualMachineInfo()
		if err != nil {
			return err
		}
	} else {
		data, err := os.ReadFile(manager.cacheFile)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(data, &manager.virtualMachineInfo); err != nil {
			return err
		}
	}

	return nil
}

func (manager *AzureVirtualMachineManager) saveVirtualMachineInfo() error {
	data, err := json.MarshalIndent(manager.virtualMachineInfo, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(manager.cacheFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (manager *AzureVirtualMachineManager) generateVirtualMachineInfo() error {
	manager.virtualMachineInfo = map[string]AzureVirtualMachineInfo{}

	subClient, err := armsubscriptions.NewClient(manager.credential, nil)
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
			vmClient, err := armcompute.NewVirtualMachinesClient(*sub.SubscriptionID, manager.credential, nil)
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

					manager.virtualMachineInfo[*vm.Name] = AzureVirtualMachineInfo{
						SubscriptionId: *sub.SubscriptionID,
						ResourceGroup:  parsedId.ResourceGroupName,
					}
				}
			}
		}
	}

	return nil
}

func (manager *AzureVirtualMachineManager) GetVirtualMachineState(name string) (string, error) {
	vm, found := manager.virtualMachineInfo[name]
	if !found {
		err := manager.loadVirtualMachineInfo(false)
		if err != nil {
			return "", err
		}

		vm, found = manager.virtualMachineInfo[name]
		if !found {
			return "", fmt.Errorf("virtual machine '%s' not found", name)
		}
	}

	clientFactory, err := armcompute.NewClientFactory(vm.SubscriptionId, manager.credential, nil)
	if err != nil {
		return "", err
	}

	client := clientFactory.NewVirtualMachinesClient()
	ctx := context.Background()
	resp, err := client.InstanceView(ctx, vm.ResourceGroup, name, nil)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve instance view for '%s'", vm)
	}

	for _, status := range resp.Statuses {
		if strings.HasPrefix(*status.Code, "PowerState") {
			return strings.Split(*status.Code, "/")[1], nil
		}
	}

	return "", fmt.Errorf("failed to get status for '%s'", vm)
}

func (manager *AzureVirtualMachineManager) StartStopVirtualMachine(name string, start, wait, noCache bool) error {
	vm, found := manager.virtualMachineInfo[name]
	if !found {
		err := manager.loadVirtualMachineInfo(noCache)
		if err != nil {
			return err
		}

		vm, found = manager.virtualMachineInfo[name]
		if !found {
			return fmt.Errorf("virtual machine '%s' not found", name)
		}
	}

	clientFactory, err := armcompute.NewClientFactory(vm.SubscriptionId, manager.credential, nil)
	if err != nil {
		return err
	}

	client := clientFactory.NewVirtualMachinesClient()

	ctx := context.Background()
	if start {
		poller, err := client.BeginStart(ctx, vm.ResourceGroup, name, nil)
		if err != nil {
			return err
		}

		if wait {
			_, err = poller.PollUntilDone(ctx, nil)
			if err != nil {
				return err
			}
		}
	} else {
		poller, err := client.BeginDeallocate(ctx, vm.ResourceGroup, name, nil)
		if err != nil {
			return err
		}

		if wait {
			_, err = poller.PollUntilDone(ctx, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
