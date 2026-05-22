package mock

import (
	"alkalax/skip-shutdown-html/internal/models"
	"time"
)

var Vms = []models.VirtualMachineInfo{
	{Name: "VM1", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: false},
	{Name: "VM2", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: true},
	{Name: "VM3", ShutdownTime: time.Date(2024, 6, 30, 22, 0, 0, 0, time.UTC), SkipToday: false},
}
