package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"alkalax/vm-booter/internal/provider/azure"
)

func main() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startCmdWait := startCmd.Bool("wait", false, "Wait for the start command to complete")
	startCmdNoCache := startCmd.Bool("no-cache", false, "Search for machines without consulting cache file")

	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	stopCmdWait := stopCmd.Bool("wait", false, "Wait for the stop command to complete")
	stopCmdNoCache := stopCmd.Bool("no-cache", false, "Search for machines without consulting cache file")

	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)

	logsCmd := flag.NewFlagSet("logs", flag.ExitOnError)

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

		vmNames := strings.Split(startCmd.Arg(0), ",")
		for _, vmName := range vmNames {
			fmt.Printf("Starting VM: %s\n", vmName)

			err = manager.StartStopVirtualMachine(vmName, true, *startCmdWait, *startCmdNoCache)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("Done.")
		}
	case "stop":
		stopCmd.Parse(os.Args[2:])
		if stopCmd.NArg() < 1 {
			fmt.Println("error: virtual machine name is required")
			os.Exit(1)
		}

		vmNames := strings.Split(stopCmd.Arg(0), ",")
		for _, vmName := range vmNames {
			fmt.Printf("Stopping VM: %s\n", vmName)

			err = manager.StartStopVirtualMachine(vmName, false, *stopCmdWait, *stopCmdNoCache)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("Done.")
		}
	case "status":
		statusCmd.Parse(os.Args[2:])
		if statusCmd.NArg() < 1 {
			fmt.Println("error: virtual machine name is required")
			os.Exit(1)
		}

		vmNames := strings.Split(statusCmd.Arg(0), ",")
		for _, vmName := range vmNames {
			state, err := manager.GetVirtualMachineState(vmName)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			fmt.Printf("%s: %s\n", vmName, state)
		}
	case "logs":
		logsCmd.Parse(os.Args[2:])
		if logsCmd.NArg() < 1 {
			fmt.Println("error: virtual machine name is required")
			os.Exit(1)
		}

		err := manager.GetActivityLogs(logsCmd.Arg(0))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	default:
		fmt.Println("error: unknown command")
		os.Exit(1)
	}
}
