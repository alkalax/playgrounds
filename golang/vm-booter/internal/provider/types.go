package provider

type VirtualMachineManager interface {
	StartStopVirtualMachine(name string, start, wait, noCache bool) error
}
