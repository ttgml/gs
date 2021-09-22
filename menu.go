package main

import (
	"fmt"
	tcell "github.com/gdamore/tcell/v2"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type menuItem struct {
	Text     string
	Shortcut tcell.Key
	Selected func()
}

type Menu struct {
	*tview.Box
	Items []menuItem
	Tip string
}

func NewMenu() *Menu {
	m := &Menu{
		Box: tview.NewBox(),
	}
	return m
}

func (m *Menu) Draw(screen tcell.Screen) {
	x, y, width, _ := m.GetInnerRect()
	x++

	for _, i := range m.Items {
		tview.Print(screen, tcell.KeyNames[i.Shortcut], x, y, width-2, tview.AlignLeft, tcell.Color(tcell.GetColor("red")))
		x += 3

		tview.Print(screen, i.Text, x, y, width, tview.AlignLeft, tcell.Color(tcell.GetColor("white")))
		x += runewidth.StringWidth(i.Text) + 2
	}
	tview.PrintSimple(screen, m.Tip, width-LenWD(m.Tip)-2,y)
}

func (m *Menu) AddItem(text string, shortcut tcell.Key, selected func()) *Menu {
	m.Items = append(m.Items, menuItem{
		Text:     text,
		Shortcut: shortcut,
		Selected: selected,
	})

	return m
}

func (m *Menu) Clear() *Menu {
	m.Items = nil

	return m
}

func (m *Menu) SetTip(tip string){
	m.Tip = tip
}

func (m *Menu) SetRateTip(current,all int){
	m.Tip = fmt.Sprintf("%v/%v",current,all)
}