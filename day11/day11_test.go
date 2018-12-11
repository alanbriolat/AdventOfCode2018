package day11

import "testing"

func TestFuelGrid_PowerLevel(t *testing.T) {
	tables := []struct{
		serialNo, x, y, power int
	}{
		{8, 3, 5, 4},
		{57, 122, 79, -5},
		{39, 217, 196, 0},
		{71, 101, 153, 4},
	}

	for _, table := range tables {
		fg := FuelGrid{SerialNo: table.serialNo}
		result := fg.CellPower(table.x, table.y)
		if result != table.power {
			t.Errorf("%+v != %v", table, result)
		}
	}
}

func TestFuelGrid_GroupPower(t *testing.T) {
	tables := []struct{
		serialNo, x, y, power int
	}{
		{18, 33, 45, 29},
		{42, 21, 61, 30},
	}

	for _, table := range tables {
		fg := FuelGrid{SerialNo: table.serialNo}
		result := fg.GroupPower(table.x, table.y, 3)
		if result != table.power {
			t.Errorf("%+v != %v", table, result)
		}
	}
}

func TestFuelGrid_FindBestGroup(t *testing.T) {
	tables := []struct{
		serialNo, bestX, bestY, bestPower int
	}{
		{18, 33, 45, 29},
		{42, 21, 61, 30},
	}

	for _, table := range tables {
		fg := FuelGrid{SerialNo: table.serialNo}
		bestX, bestY, bestPower := fg.FindBestGroup(3)
		if bestX != table.bestX || bestY != table.bestY || bestPower != table.bestPower {
			t.Errorf("%+v != %v,%v=%v", table, bestX, bestY, bestPower)
		}
	}
}

func TestFuelGrid_FindBestGroupAnySize(t *testing.T) {
	tables := []struct{
		serialNo, bestX, bestY, bestSize, bestPower int
	}{
		{18, 90, 269, 16, 113},
		{42, 232, 251, 12, 119},
	}

	for _, table := range tables {
		fg := FuelGrid{SerialNo: table.serialNo}
		bestX, bestY, bestSize, bestPower := fg.FindBestGroupAnySize()
		if bestX != table.bestX || bestY != table.bestY ||
			bestPower != table.bestPower || bestSize != table.bestSize {
			t.Errorf("%+v != %v,%v,%v=%v", table, bestX, bestY, bestSize, bestPower)
		}
	}
}
