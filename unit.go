package main

import (
	"github.com/coreos/go-systemd/v22/dbus"
)

type UnitItem struct {
	TypeID int
	Status   dbus.UnitStatus
	Property map[string]interface{}
	TypeProperty map[string]interface{}
}

type UnitTree struct {
	Type   string
	TypeID int
	Items  UnitItems
}

type UnitItems []*UnitItem

func (u UnitTree) AddItem(item *UnitItem) {
	u.Items = append(u.Items, item)
}

func (u UnitItems) Len() int{
	return len(u)
}

func (u UnitItems) Less(i,j int) bool {
	return (u)[i].Status.Name[0] < (u)[j].Status.Name[1]
}

func (u UnitItems) Swap(i, j int){
	(u)[i],(u)[j] = (u)[j],(u)[i]
}