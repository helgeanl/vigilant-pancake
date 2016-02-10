package main

import (
	"./driver"
	"fmt"
)

func main() {
	fmt.Println("Started")
	driver.Init(driver.ET_simulation)
	driver.Set_bit(driver.LIGHT_COMMAND1)

	fmt.Println("Done.")

	// We wait to make sure the driver starts all its threads & connections
	select {}

}
