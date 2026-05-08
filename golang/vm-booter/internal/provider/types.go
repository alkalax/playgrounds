package provider

type VirtualMachineManager interface {
	StartStopVirtualMachine(name string, start bool) error
}
