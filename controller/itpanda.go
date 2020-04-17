package controller

import (
	"pandaBook/model"
)

type ItPandaController struct {

}

func (I *ItPandaController) List() {
	new(model.ItPandaModel).List()
}



