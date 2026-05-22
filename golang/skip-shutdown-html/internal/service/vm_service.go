package service

import (
	"alkalax/skip-shutdown-html/internal/models"
	mock "alkalax/skip-shutdown-html/internal/storage"
)

func FindVM(name string) *models.VirtualMachineInfo {
	for i := range mock.Vms {
		if mock.Vms[i].Name == name {
			return &mock.Vms[i]
		}
	}
	return nil
}

func GetVMs() []models.VirtualMachineInfo {
	return mock.Vms
}
