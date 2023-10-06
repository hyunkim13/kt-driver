package main

import(
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/hyunkim13/kt-driver/driver"
)


func main() {
	plugin.RegisterDriver(driver.NewDriver("",""))
}
