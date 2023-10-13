package main

import (
	"./driver"
	"github.com/docker/machine/libmachine/drivers/plugin"
)

func main() {
	plugin.RegisterDriver(driver.NewDriver())
}
