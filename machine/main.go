package main

import(
	"github.com/docker/machine/libmachine/drivers/plugin"
	"./driver"
)


func main() {
	plugin.RegisterDriver(driver.NewDriver("",""))
}
