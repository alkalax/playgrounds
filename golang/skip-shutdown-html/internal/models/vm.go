package models

import "time"

type VirtualMachineInfo struct {
	Name         string
	ShutdownTime time.Time
	SkipToday    bool
}

type PageData struct {
	VMs          []VirtualMachineInfo
	SelectedName string
	Selected     *VirtualMachineInfo
}
