package main

import (
	"pandaBook/model"
	"pandaBook/router"
)

func main() {

	defer new(model.BaseModel).GetDbInstance().Close()
	router.Init().Run(":8081") // listen and serve on 0.0.0.0:8080
}
