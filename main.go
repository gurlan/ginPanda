package main

import (
	"pandaBook/router"
)

func main() {
	router.Init().Run(":8081") // listen and serve on 0.0.0.0:8080
}
