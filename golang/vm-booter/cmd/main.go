package main

import (
	"fmt"
	"os"

	"alkalax/vm-booter/internal/provider/azure"
)

func main() {
	manager, err := azure.NewVirtualMachineManager("vm_info.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = manager.StartStopVirtualMachine("ubuntu-0", true)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
