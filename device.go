package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

func (device Device) CreateTerminalIdle() {
	LogInfo(device.Name, "Creating idle")
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError(device.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	var terminalInputIdle TerminalInputIdle
	terminalInputIdle.DTS = time.Now()
	terminalInputIdle.IdleID = 1
	terminalInputIdle.Interval = 0
	terminalInputIdle.DeviceID = device.OID
	db.Save(&terminalInputIdle)
	LogInfo(device.Name, "Idle created")

}

func (device Device) CloseTerminalIdle() {
	LogInfo(device.Name, "Closing idle")
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError(device.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return
	}
	defer db.Close()
	var terminalInputIdle TerminalInputIdle
	db.Where("DTE is null").Where("DeviceID = ?", device.OID).First(&terminalInputIdle)
	db.Model(&terminalInputIdle).Where("OID =?", terminalInputIdle.OID).UpdateColumns(TerminalInputIdle{DTE: time.Now()})
	LogInfo(device.Name, "Idle closed")

}

func (device Device) GetActualIdleDuration(devicePortId int) int {
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError(device.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	var deviceInputDigital DeviceInputDigital
	db.Where("DevicePortID=?", devicePortId).Last(&deviceInputDigital)
	return int(time.Now().Sub(deviceInputDigital.DT).Seconds())
}

func (device Device) GetDefaultIdleDuration(workplaceId int) int {
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError(device.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	var workplace Workplace
	db.Where("OID=?", workplaceId).Find(&workplace)
	return workplace.IdleFromTime
}

func (device Device) CheckForOpenTerminalIdle() bool {
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError(device.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return false
	}
	defer db.Close()
	var terminalInputIdle TerminalInputIdle
	db.Where("DTE is null").Where("DeviceID = ?", device.OID).Find(&terminalInputIdle)
	return terminalInputIdle.OID > 0
}

func (device Device) GetActualDevicePortValue(devicePortId int) int {
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError(device.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	var deviceInputDigital DeviceInputDigital
	db.Where("DevicePortID=?", devicePortId).Last(&deviceInputDigital)
	return deviceInputDigital.Data
}

func (device Device) GetDevicePortOid(workplaceID int) int {
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError(device.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	var workplacePort WorkplacePort
	db.Where("WorkplaceID=?", workplaceID).Find(&workplacePort)
	return workplacePort.DevicePortID
}

func (device Device) GetWorkplaceOid() int {
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError(device.Name, "Problem opening "+DatabaseName+" database: "+err.Error())
		return 0
	}
	defer db.Close()
	var workplace Workplace
	db.Where("DeviceID=?", device.OID).Find(&workplace)
	return workplace.OID
}

func (device Device) Sleep(start time.Time) {
	if time.Since(start) < (downloadInSeconds * time.Second) {
		sleepTime := downloadInSeconds*time.Second - time.Since(start)
		LogInfo(device.Name, "Sleeping for "+sleepTime.String())
		time.Sleep(sleepTime)
	}
}
