package models

type WareHouse struct { //TODO:Refactor this for optimal memory usage
	ID          uint
	Name        string
	IsAvailable bool
	Products    []Product
}
