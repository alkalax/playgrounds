package main

import (
	"flag"
	"fmt"
	"os"

	"alkalax/vm-booter/internal/provider/azure"
)

func main() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startCmdWait := startCmd.Bool("wait", false, "Wait for the start command to complete")
	startCmdNoCache := startCmd.Bool("no-cache", false, "Search for machines without consulting cache file")

	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	stopCmdWait := stopCmd.Bool("wait", false, "Wait for the stop command to complete")
	stopCmdNoCache := stopCmd.Bool("no-cache", false, "Search for machines without consulting cache file")

	if len(os.Args) < 2 {
		fmt.Println("error: expected subcommands")
		os.Exit(1)
	}

	manager, err := azure.NewVirtualMachineManager("vm_info.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		if startCmd.NArg() < 1 {
			fmt.Println("error: virtual machine name is required")
			os.Exit(1)
		}

		vmName := startCmd.Arg(0)
		fmt.Printf("Starting VM: %s\n", vmName)

		err = manager.StartStopVirtualMachine(vmName, true, *startCmdWait, *startCmdNoCache)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Done.")
	case "stop":
		stopCmd.Parse(os.Args[2:])
		if stopCmd.NArg() < 1 {
			fmt.Println("error: virtual machine name is required")
			os.Exit(1)
		}

		vmName := stopCmd.Arg(0)
		fmt.Printf("Stopping VM: %s\n", vmName)

		err = manager.StartStopVirtualMachine(vmName, false, *stopCmdWait, *stopCmdNoCache)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Done.")

	default:
		fmt.Println("error: unknown command")
		os.Exit(1)
	}
}
