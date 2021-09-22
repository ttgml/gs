package main

import (
	"strings"

	tcell "github.com/gdamore/tcell/v2"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type tabItem struct {
	id       int
	Text     string
	Shortcut tcell.Key
	Selected func(id int)
	Ext      string
	UnitType string
}

type Tab struct {
	*tview.Box
	Items      []tabItem
	CurrentTab int
	NewUnit    tabItem  //创建按钮
}

var MainTabs = []tabItem{
	tabItem{0, "Services(s)", tcell.KeyCtrlS, nil, ".service","Service"},
	tabItem{1, "Targets(t)", tcell.KeyCtrlT, nil, ".target","Target"},
	tabItem{2, "Sockets(o)", tcell.KeyCtrlO, nil, ".socket","Socket"},
	tabItem{3, "Timers(i)", tcell.KeyCtrlI, nil, ".timer","Timer"},
	tabItem{4, "Mount(p)", tcell.KeyCtrlP, nil, ".mount","Mount"},
}

func NewTab(app *tview.Application) *Tab {
	t := &Tab{
		Box: tview.NewBox(),
	}
	t.Items = MainTabs
	t.CurrentTab = 0
	t.NewUnit = tabItem{100, "Create", tcell.KeyCtrlN, nil, "Create",""}
	return t
}

func (t *Tab) Draw(screen tcell.Screen) {
	x, y, width, height := t.GetInnerRect()
	CleanArea(screen, x, y, width, height)
	for _, i := range t.Items {
		wordLen := runewidth.StringWidth(i.Text)
		if i.id == t.CurrentTab {
			tview.PrintSimple(screen, string(tcell.RuneULCorner), x, y)
			tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine), wordLen+2), x+1, y)
			tview.PrintSimple(screen, string(tcell.RuneURCorner), x+wordLen+2+1, y)
			tview.PrintSimple(screen, string(tcell.RuneVLine), x, y+1)
			//tview.PrintSimple(screen," " + i.Text + " ",x+1,y+1)
			tview.Print(screen, " "+i.Text+" ", x+1, y+1, width-2, tview.AlignLeft, tcell.GetColor("green"))
			tview.PrintSimple(screen, string(tcell.RuneVLine), x+wordLen+2+1, y+1)
		} else {
			tview.PrintSimple(screen, " "+i.Text+" ", x+1, y+1)
		}
		x = x + wordLen + 3
	}
	tview.PrintSimple(screen, t.NewUnit.Text, width-runewidth.StringWidth(t.NewUnit.Text)-2, y+1)

}

func (t *Tab) SetCurrentTab(tabID int) {
	t.CurrentTab = tabID
}
func (s *Tab) NextTab() {
	if s.CurrentTab == len(MainTabs)-1{
		s.CurrentTab = 0
	}else{
		s.CurrentTab ++
	}
}
func (s *Tab) PreviousTab() {
	if s.CurrentTab == 0{
		s.CurrentTab = len(MainTabs)-1
	}else{
		s.CurrentTab --
	}
}

func (t *Tab) Clear() *Tab {
	t.Items = nil
	return t
}

func CleanArea(screen tcell.Screen, x, y, width, height int) {
	for yy := y; yy < height; yy++ {
		for xx := x; xx < width; xx++ {
			tview.PrintSimple(screen, " ", xx, yy)
		}
	}
}
