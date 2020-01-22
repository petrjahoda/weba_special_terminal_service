package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Device struct {
	OID        int    `gorm:"column:OID"`
	Name       string `gorm:"column:Name"`
	DeviceType int    `gorm:"column:DeviceType"`
}

func (Device) TableName() string {
	return "device"
}

type Workplace struct {
	OID          int    `gorm:"column:OID"`
	Name         string `gorm:"column:Name"`
	IdleFromTime int    `gorm:"column:IdleFromTime"`
	DeviceID     string `gorm:"column:DeviceID"`
}

func (Workplace) TableName() string {
	return "workplace"
}

type WorkplacePort struct {
	OID          int `gorm:"column:OID"`
	DevicePortID int `gorm:"column:DevicePortID"`
	WorkplaceID  int `gorm:"column:WorkplaceID"`
}

func (WorkplacePort) TableName() string {
	return "workplace_port"
}

type DeviceInputDigital struct {
	OID          int       `gorm:"column:OID;primary_key"`
	DevicePortID int       `gorm:"column:DevicePortID"`
	DT           time.Time `gorm:"column:DT"`
	Data         int       `gorm:"column:Data"`
}

func (DeviceInputDigital) TableName() string {
	return "device_input_digital"
}

type TerminalInputIdle struct {
	OID      int       `gorm:"column:OID;primary_key"`
	DTS      time.Time `gorm:"column:DTS"`
	DTE      time.Time `gorm:"column:DTE;default:'null'"`
	IdleID   int       `gorm:"column:IdleID"`
	UserID   int       `gorm:"column:UserID;default:'null'"`
	Interval float32   `gorm:"column:Interval"`
	DeviceID int       `gorm:"column:DeviceID"`
	Note     string    `gorm:"column:Note"`
}

func (TerminalInputIdle) TableName() string {
	return "terminal_input_idle"
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
