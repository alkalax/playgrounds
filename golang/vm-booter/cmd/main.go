package main

import (
	"fmt"
	"os"

	"alkalax/vm-booter/internal/provider/azure"
)

func main() {
	manager := azure.NewVirtualMachineManager("vm_info.json")
	err := manager.StartStopVirtualMachine("ubuntu-0", true)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
