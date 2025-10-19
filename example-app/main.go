package main

import "example-app/basic"
import "example-app/postgres"

func main() {

	basic.Run()
	
	postgres.Run()

}
