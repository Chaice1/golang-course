package main

import (
	collectorapp "github.com/Chaice1/golang-course/task2/collector/app"
)

func main() {
	config := collectorapp.Config{
		GRPCport: ":8082",
	}
	collectorapp.RunApp(config)
}
