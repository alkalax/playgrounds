package storage

import (
	"alkalax/skip-shutdown-html/internal/models"
	"time"
)

type MockVMRepository struct {
	vms []models.VirtualMachineInfo
}

func NewMockVMRepository() *MockVMRepository {
	return &MockVMRepository{
		vms: []models.VirtualMachineInfo{
			{Name: "VM1", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: false},
			{Name: "VM2", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: true},
			{Name: "VM3", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: false},
		},
	}
}

func (m *MockVMRepository) GetVMs() []models.VirtualMachineInfo {
	return m.vms
}

func (m *MockVMRepository) FindVM(name string) *models.VirtualMachineInfo {
	for i := range m.vms {
		if m.vms[i].Name == name {
			return &m.vms[i]
		}
	}

	return nil
}
