package storage

import (
	"alkalax/skip-shutdown-html/internal/models"
	"database/sql"

	_ "modernc.org/sqlite"
)

type SQLiteVMRepository struct {
	dbFile string
}

func NewSQLiteVMRepository(dbFile string) *SQLiteVMRepository {
	return &SQLiteVMRepository{dbFile: dbFile}
}

func (m *SQLiteVMRepository) GetVMs() ([]models.VirtualMachineInfo, error) {
	vms := []models.VirtualMachineInfo{}

	db, err := sql.Open("sqlite", m.dbFile)
	if err != nil {
		return vms, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT Name FROM VirtualMachines")
	if err != nil {
		return vms, err
	}

	for rows.Next() {
		var vm models.VirtualMachineInfo
		if err = rows.Scan(&vm.Name); err != nil {
			return vms, err
		}

		vms = append(vms, vm)
	}

	return vms, nil
}

func (m *SQLiteVMRepository) FindVM(name string) (*models.VirtualMachineInfo, error) {
	db, err := sql.Open("sqlite", m.dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var vm models.VirtualMachineInfo
	err = db.QueryRow(
		"SELECT Name, ShutdownTime, SkipToday FROM VirtualMachines WHERE Name = ?",
		name,
	).Scan(&vm.Name, &vm.ShutdownTime, &vm.SkipToday)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &vm, nil
}
