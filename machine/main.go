package main

import (
	"github.com/hyunkim13/kt-driver/machine/driver"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(driver.NewDriver())
}
