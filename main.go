package main

import(
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/hyunkim13/kt-driver/driver/kt"
)


func main() {
	plugin.RegisterDriver(kt.NewDriver())
}
