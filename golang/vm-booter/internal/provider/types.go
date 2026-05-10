package provider

type VirtualMachineManager interface {
	StartStopVirtualMachine(name string, start, wait bool) error
}
