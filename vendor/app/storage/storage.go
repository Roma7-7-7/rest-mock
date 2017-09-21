package storage

import "app/api"
import "app/shared"

//DataStorage :
type DataStorage interface {
	GetAll() ([]api.Endpoint, error)
}

//Storage :
var Storage DataStorage

func init() {
	Storage = BuildFileStorage(shared.ResponseFolder)
}
