package main

import(
	"github.com/docker/machine/libmachine/drivers/plugin"
	"./driver/kt"
)


func main() {
	plugin.RegisterDriver(kt.NewDriver("",""))
}
