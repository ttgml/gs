package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

var conn *dbus.Conn
var ctx = context.TODO()
var unitDBS []UnitTree

func main_tui(c *cli.Context) error {
	tuiApp := tview.NewApplication()
	pages := tview.NewPages()
	menu := NewMenu()
	menu.AddItem("Help", tcell.KeyF1, nil)
	menu.AddItem("Refresh", tcell.KeyF2,nil)
	menu.AddItem("Reload",tcell.KeyF3,nil)
	menu.AddItem("Edit",tcell.KeyF4,nil)
	menu.AddItem("Create", tcell.KeyF5,nil)
	menu.AddItem("Filter", tcell.KeyF6,nil)
	menu.AddItem("Quit", tcell.KeyF7,nil)
	tabBox := NewTab(tuiApp)
	tabBox.SetCurrentTab(0)
	slist := NewSList()
	slist.SetCurrentTab(0)
	slist.Refresh()
	main_flex := tview.NewFlex().SetDirection(tview.FlexRow)
	main_flex.AddItem(tabBox, 2, 1, false)
	main_flex.AddItem(slist, 0, 1, true)
	main_flex.AddItem(menu, 1, 1, false)

	info_flex := NewInformationWidge()
	create_flex := NewCreateWidge()
	pages.AddPage("info", info_flex.page, true, false)
	pages.AddPage("main", main_flex, true, true)
	pages.AddPage("create", create_flex.page, true, false)

	log.Println(main_flex.HasFocus())
	tuiApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		log.Println("KEY DOWN: ",event.Key())
		log.Println(main_flex.HasFocus())
		log.Println(slist.HasFocus())
		log.Println(tabBox.HasFocus())
		log.Println(info_flex.page.HasFocus())
		if main_flex.HasFocus(){
			for _, i := range tabBox.Items {
				if event.Key() == i.Shortcut {
					tabBox.SetCurrentTab(i.id)
					slist.SetCurrentTab(i.id)
					slist.Refresh()
					break
				}
			}
			if event.Key() == tcell.KeyDown || event.Rune() == 'j' {
				slist.DownKey()
			}
			if event.Key() == tcell.KeyUp || event.Rune() == 'k' {
				slist.UpKey()
			}
			if event.Key() == tcell.KeyPgDn{
				slist.NextPage()
			}
			if event.Key() == tcell.KeyPgUp{
				slist.PreviousPage()
			}
			if event.Key() == tcell.KeyEnd{
				slist.EndPage()
			}
			if event.Key() == tcell.KeyHome{
				slist.HomePage()
			}

			if event.Rune() == 'H'{
				slist.PreviousTab()
				tabBox.PreviousTab()
				slist.Refresh()
			}
			if event.Rune() == 'L'{
				slist.NextTab()
				tabBox.NextTab()
				slist.Refresh()
			}
			if event.Key() == tcell.KeyCtrlN || event.Key() == tcell.KeyF5{
				pages.SwitchToPage("create")
				pages.ShowPage("create")
				tuiApp.SetFocus(create_flex.page)
				log.Println("Create New Unit")
			}
			//在列表界面，按回车(256)，查看当前光标所在行的服务详细信息
			if event.Key() == tcell.KeyEnter {
				pages.SwitchToPage("info")
				pages.ShowPage("info")
				tuiApp.SetFocus(info_flex.page)
				name:=unitDBS[slist.currentTab].Items[slist.currentItem[slist.currentTab]].Status.Name
				info:=*unitDBS[slist.currentTab].Items[slist.currentItem[slist.currentTab]]
				info_flex.Set(name,info)
			}
		}
		if info_flex.page.HasFocus(){
			if event.Key() == tcell.KeyEsc{
				pages.SwitchToPage("main")
				tuiApp.SetFocus(main_flex)
				log.Println("focus slist")
			}
			if event.Key() == tcell.KeyCtrlN || event.Key() == tcell.KeyF5{
				pages.SwitchToPage("create")
				pages.ShowPage("create")
				tuiApp.SetFocus(create_flex.page)
				log.Println("Create New Unit")
			}
		}
		if create_flex.page.HasFocus(){
			if event.Key() == tcell.KeyEsc{
				pages.SwitchToPage("main")
				tuiApp.SetFocus(main_flex)
				log.Println("focus slist")
			}
		}
		menu.SetRateTip(slist.currentItem[slist.currentTab]+1,slist.GetItemCount())
		return event
	})

	if err := tuiApp.SetRoot(pages, true).SetFocus(main_flex).Run(); err != nil {
		panic(err)
	}

	return nil
}

func init() {
	var err error
	conn, err = dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Load...")
	initUnitDB()

}


//初始化db，其实就是一个字典
func initUnitDB() {
	units, err := conn.ListUnitsContext(ctx)
	if err != nil {
		log.Println(err)
	}
	for _, tab := range MainTabs {
		unitDBS = append(unitDBS, UnitTree{Type: tab.UnitType, TypeID: tab.id})
	}

	for _, item := range units {
		propertys, err := conn.GetUnitPropertiesContext(ctx, item.Name)
		if err != nil {
			log.Println(err)
		}
		unititem := new(UnitItem)
		unititem.Property = propertys
		unititem.Status = item
		if strings.Contains(item.Name, MainTabs[0].Ext) {
			pro,_:=getUnitTypeProp(MainTabs[0].id,item.Name)
			unititem.TypeProperty = pro
			unititem.TypeID = MainTabs[0].id
			unitDBS[0].Items = append(unitDBS[0].Items, unititem)
		}
		if strings.Contains(item.Name, MainTabs[1].Ext) {
			pro,_:=getUnitTypeProp(MainTabs[1].id,item.Name)
			unititem.TypeProperty = pro
			unititem.TypeID = MainTabs[1].id
			unitDBS[1].Items = append(unitDBS[1].Items, unititem)
		}
		if strings.Contains(item.Name, MainTabs[2].Ext) {
			pro,_:=getUnitTypeProp(MainTabs[2].id,item.Name)
			unititem.TypeProperty = pro
			unititem.TypeID = MainTabs[2].id
			unitDBS[2].Items = append(unitDBS[2].Items, unititem)
		}
		if strings.Contains(item.Name, MainTabs[3].Ext) {
			pro,_:=getUnitTypeProp(MainTabs[3].id,item.Name)
			unititem.TypeProperty = pro
			unititem.TypeID = MainTabs[3].id
			unitDBS[3].Items = append(unitDBS[3].Items, unititem)
		}
		if strings.Contains(item.Name, MainTabs[4].Ext) {
			pro,_:=getUnitTypeProp(MainTabs[4].id,item.Name)
			unititem.TypeProperty = pro
			unititem.TypeID = MainTabs[4].id
			unitDBS[4].Items = append(unitDBS[4].Items, unititem)
		}

	}
	//排序
	for i := 0; i < len(unitDBS); i++ {
		sort.Sort(unitDBS[i].Items)
	}
}

func getUnitTypeProp(typeID int,unitName string) (map[string]interface{},error) {
	properties,err:=conn.GetUnitTypePropertiesContext(ctx,unitName,MainTabs[typeID].UnitType)
	if err!=nil{
		return nil, err
	}
	return properties,nil
}