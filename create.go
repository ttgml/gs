package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)



type CreateWidge struct {
	page *tview.Flex
	//createFrom *CreateF
	menu *Menu
	form *tview.Form
}

//func (i *CreateF) Draw(screen tcell.Screen) {
//	x, y, width, height := i.GetInnerRect()
//	log.Println(x,y,width,height)
//	//CleanArea(screen,x,y,width,height)
//	//tview.PrintSimple(screen,i.Service,x,y)
//}

func NewCreateWidge() *CreateWidge {

	//最下面一行的快捷键菜单
	formP_menu := NewMenu()
	formP_menu.AddItem("Help", tcell.KeyF1, nil)
	formP_menu.AddItem("Reset", tcell.KeyF2,nil)
	formP_menu.AddItem("Preview",tcell.KeyF3,nil)
	formP_menu.AddItem("Disable",tcell.KeyF4,nil)
	formP_menu.AddItem("Preview", tcell.KeyF5,nil)

	formP := tview.NewForm()

	formP.AddInputField("Unit Name","",100,nil,nil)
	formP.AddInputField("Description","",100,nil,nil)
	formP.AddDropDown("Type",[]string{"simple","forking","oneshot","dbus","notify","idle"},0,nil)
	formP.AddInputField("ExecStart","",100,nil,nil)
	formP.AddInputField("ExecStop","",100,nil,nil)
	formP.AddDropDown("Restart",[]string{"no","always","on-success","on-abnormal","on-abort","on-watchdog",},0,nil)
	formP.SetTitle("Create Service Unit").SetTitleAlign(tview.AlignLeft).SetBorder(true)

	//布局
	pflex := tview.NewFlex().SetDirection(tview.FlexRow)
	pflex.AddItem(formP,0,1,true)
	pflex.AddItem(formP_menu,1,1,false)
	widge := &CreateWidge{
		page: pflex,
		menu: formP_menu,
		form: formP,
	}
	return widge
}
