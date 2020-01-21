package main

import (
	"github.com/jinzhu/gorm"
	"strconv"
	"sync"
	"time"
)

const version = "2019.4.3.31"
const deleteLogsAfter = 240 * time.Hour
const downloadInSeconds = 10

var (
	activeDevices  []Device
	runningDevices []Device
	deviceSync     sync.Mutex
)

func main() {
	LogDirectoryFileCheck("MAIN")
	LogInfo("MAIN", "Program version "+version+" started")
	CreateConfigIfNotExists()
	LoadSettingsFromConfigFile()
	LogDebug("MAIN", "Using ["+DatabaseType+"] on "+DatabaseIpAddress+":"+DatabasePort+" with database "+DatabaseName)
	for {
		start := time.Now()
		LogInfo("MAIN", "Program running")
		UpdateActiveDevices("MAIN")
		DeleteOldLogFiles()
		LogInfo("MAIN", "Active devices: "+strconv.Itoa(len(activeDevices))+", running devices: "+strconv.Itoa(len(runningDevices)))
		for _, activeDevice := range activeDevices {
			activeDeviceIsRunning := CheckDevice(activeDevice)
			if !activeDeviceIsRunning {
				go RunDevice(activeDevice)
			}
		}
		if time.Since(start) < (downloadInSeconds * time.Second) {
			sleepTime := downloadInSeconds*time.Second - time.Since(start)
			LogInfo("MAIN", "Sleeping for "+sleepTime.String())
			time.Sleep(sleepTime)
		}
	}
}

func CheckDevice(device Device) bool {
	for _, runningDevice := range runningDevices {
		if runningDevice.Name == device.Name {
			return true
		}
	}
	return false
}

func RunDevice(device Device) {
	LogInfo(device.Name, "Device started running")
	deviceSync.Lock()
	runningDevices = append(runningDevices, device)
	deviceSync.Unlock()
	deviceIsActive := true
	var workplaceOid int = device.GetWorkplaceOid()
	var devicePortOid int = device.GetDevicePortOid(workplaceOid)
	for deviceIsActive {
		start := time.Now()
		LogInfo(device.Name, "Device loop started")
		var actualDevicePortStatus int = device.GetActualDevicePortStatus(devicePortOid)
		var terminalInputIdleIsOpened bool = device.CheckForOpenTerminalInputIdle()
		if actualDevicePortStatus == 1 && terminalInputIdleIsOpened {
			device.CloseOpenTerminalInputIdle()
		}
		if actualDevicePortStatus == 0 && !terminalInputIdleIsOpened {
			var workplaceIdleCycleTime int = device.GetWorkplaceIdleTime(workplaceOid)
			var actualIdleCycleTime int = device.GetActualIdleCycleTime(devicePortOid)
			if actualIdleCycleTime > workplaceIdleCycleTime {
				device.CreateTerminalInputIdle()
			}
		}
		device.Sleep(start)
		deviceIsActive = CheckActive(device)
	}
	RemoveDeviceFromRunningDevices(device)
	LogInfo(device.Name, "Device not active, stopped running")

}

func CheckActive(device Device) bool {
	for _, activeDevice := range activeDevices {
		if activeDevice.Name == device.Name {
			LogInfo(device.Name, "Device still active")
			return true
		}
	}
	LogInfo(device.Name, "Device not active")
	return false
}

func RemoveDeviceFromRunningDevices(device Device) {
	for idx, runningDevice := range runningDevices {
		if device.Name == runningDevice.Name {
			deviceSync.Lock()
			runningDevices = append(runningDevices[0:idx], runningDevices[idx+1:]...)
			deviceSync.Unlock()
		}
	}
}

func UpdateActiveDevices(reference string) {
	connectionString, dialect := CheckDatabaseType()
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError(reference, "Problem opening "+DatabaseName+" database: "+err.Error())
		activeDevices = nil
		return
	}
	defer db.Close()
	db.Where("DeviceType=?", 9).Find(&activeDevices)
}
