package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Device struct {
	OID        int    `gorm:"column:OID"`
	Name       string `gorm:"column:Name"`
	DeviceType int    `gorm:"column:DeviceType"`
}

type State struct {
	gorm.Model
	Name  string `gorm:"unique"`
	Color string
	Note  string
}
type WorkplaceSection struct {
	gorm.Model
	Name string `gorm:"unique"`
	Note string
}
type Workplace struct {
	gorm.Model
	Name                   string `gorm:"unique"`
	Code                   string
	WorkplaceSectionId     uint
	ActualStateId          uint
	ActualStateDateTime    time.Time
	ActualWorkplaceModeId  uint
	ProductionPortValue    int
	ProductionPortDateTime time.Time
	PoweroffPortDateTime   time.Time
	WorkplaceModes         []WorkplaceMode
	WorkplacePorts         []WorkplacePort
	Devices                []Device
	Note                   string
}

type WorkplacePort struct {
	gorm.Model
	Name         string
	DevicePortId uint
	WorkplaceId  uint
	LowValue     float32
	HighValue    float32
	Color        string
	StateId      uint
	Note         string
}

type WorkplaceState struct {
	Id            uint `gorm:"primary_key"`
	WorkplaceId   uint
	StateId       uint
	DateTimeStart time.Time
	DateTimeEnd   time.Time
	Interval      float32
	Note          string
}

type WorkplaceMode struct {
	gorm.Model
	Name             string `gorm:"unique"`
	DowntimeInterval int
	PoweroffInterval int
	Note             string
}

type DeviceType struct {
	gorm.Model
	Name string `gorm:"unique"`
	Note string
}

type DevicePortType struct {
	gorm.Model
	Name string `gorm:"unique"`
	Note string
}

type Setting struct {
	gorm.Model
	Key     string `gorm:"unique"`
	Value   string
	Enabled bool
	Note    string
}

type DevicePort struct {
	gorm.Model
	Name               string
	Unit               string
	PortNumber         int
	DevicePortTypeId   uint
	DeviceId           uint
	ActualDataDateTime time.Time
	ActualData         string
	PlcDataType        string
	PlcDataAddress     string
	Settings           string
	Virtual            bool
	Note               string
}

type DeviceAnalogRecord struct {
	Id           uint      `gorm:"primary_key"`
	DevicePortId uint      `gorm:"unique_index:unique_analog_data"`
	DateTime     time.Time `gorm:"unique_index:unique_analog_data"`
	Data         float32
	Interval     float32
}

type DeviceDigitalRecord struct {
	Id           uint      `gorm:"primary_key"`
	DevicePortId uint      `gorm:"unique_index:unique_digital_data"`
	DateTime     time.Time `gorm:"unique_index:unique_digital_data"`
	Data         int
	Interval     float32
}

type DeviceSerialRecord struct {
	Id           uint      `gorm:"primary_key"`
	DevicePortId uint      `gorm:"unique_index:unique_serial_data"`
	DateTime     time.Time `gorm:"unique_index:unique_serial_data"`
	Data         float32
	Interval     float32
}

func CheckDatabaseType() (string, string) {
	var connectionString string
	var dialect string
	if DatabaseType == "postgres" {
		connectionString = "host=" + DatabaseIpAddress + " sslmode=disable port=" + DatabasePort + " user=" + DatabaseLogin + " dbname=" + DatabaseName + " password=" + DatabasePassword
		dialect = "postgres"
	} else if DatabaseType == "mysql" {
		connectionString = DatabaseLogin + ":" + DatabasePassword + "@tcp(" + DatabaseIpAddress + ":" + DatabasePort + ")/" + DatabaseName + "?charset=utf8&parseTime=True&loc=Local"
		dialect = "mysql"
	}
	return connectionString, dialect
}
