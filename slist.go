package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type listItem struct {
	Id            int    //
	MainText      string //name
	SecondaryText string //desct
	RunStatus     string //not running
	PropStatus    string //disable
	TipInfo       string //more info
}

type SList struct {
	*tview.Box
	items             []*listItem
	currentItem       map[int]int                                     //当前第几个 key是当前第几个tab，value是具体值
	currentTab        int                                             //当前第几个tab
	mainTextColor     tcell.Color                                     //列表的颜色
	selectedTextColor tcell.Color                                     //被选中的颜色
	RunStatusColor    tcell.Color                                     //运行状态的颜色
	PropStatusColor   tcell.Color                                     //运行状态2的颜色
	itemOffset        map[int]int                                     //列表第一行数据是在所有列表里面第几
	displayOffset     map[int]int                                     //当前应该高亮第几行，默认是1，如果切换，记录,key是当前第几个tab，value是具体值
	changed           func(index int, mainText, secondaryText string) //被改变是运行的函数
	selected          func(index int, mainText, secondaryText string) //同上
	done              func()
	height            int //列表的高度，方便翻页
}

func NewSList() *SList {
	s := &SList{
		Box:               tview.NewBox(),
		mainTextColor:     tview.Styles.PrimaryTextColor,
		selectedTextColor: tview.Styles.TertiaryTextColor,
		currentItem:       make(map[int]int),
		itemOffset:        make(map[int]int),
		displayOffset:     make(map[int]int),
	}
	s.SetBorder(true)
	return s
}

func (s *SList) GetItemCount() int {
	return len(s.items)
}

func (s *SList) AddItem(id int, mainText, secondaryText, runStatus, propStatus, tipInfo string) *SList {
	switch runStatus {
	case "active":
		runStatus = "Running"
	case "failed":
		runStatus = "Failed"
	default:
		runStatus = "Not running"
	}

	item := &listItem{
		Id:            id,
		MainText:      mainText,
		SecondaryText: secondaryText,
		RunStatus:    runStatus,
		PropStatus:     propStatus,
		TipInfo:       tipInfo,
	}
	s.items = append(s.items, item)
	return s
}

func (s *SList) Clear() *SList {
	s.items = []*listItem{}
	return s
}

func (s *SList) Draw(screen tcell.Screen) {

	x, y, width, height := s.GetInnerRect()
	//log.Println(height)
	s.height = height
	if width < 70 {
		fmt.Println("The window is too small.")
		os.Exit(1)
		return
	}
	//获取行高限制
	bottonLimit := y + height
	_, totalHeight := screen.Size()
	if bottonLimit > totalHeight {
		bottonLimit = totalHeight
	}

	//重置一下第一行
	CleanArea(screen, x, y-1, width, y)

	//根据tab位置打印第一行上边框
	switch s.currentTab {
	case 0:
		tview.PrintSimple(screen, string(tcell.RuneVLine), x-1, y-1)
		tview.PrintSimple(screen, string(tcell.RuneLLCorner), x+LenWD(MainTabs[0].Text)+2, y-1)
		tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine), width-LenWD(MainTabs[0].Text)-4), x+LenWD(MainTabs[0].Text)+4, y-1)
		tview.PrintSimple(screen, string(tcell.RuneURCorner), width, y-1)
	case 1:
		tview.PrintSimple(screen, string(tcell.RuneULCorner), x-1, y-1)
		tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine), x+LenWD(MainTabs[0].Text)+1), x, y-1)
		tview.PrintSimple(screen, string(tcell.RuneLRCorner), x+LenWD(MainTabs[1].Text)+3, y-1)
		tview.PrintSimple(screen, string(tcell.RuneLLCorner), x+LenWD(MainTabs[0].Text+MainTabs[1].Text)+5, y-1)
		tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine), width-LenWD(MainTabs[0].Text+MainTabs[1].Text)-6), x+LenWD(MainTabs[0].Text+MainTabs[1].Text)+6, y-1)
		tview.PrintSimple(screen, string(tcell.RuneURCorner), width, y-1)
	case 2:
		tview.PrintSimple(screen, string(tcell.RuneULCorner), x-1, y-1)
		tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine), x+LenWD(MainTabs[0].Text+MainTabs[1].Text)+4), x, y-1)
		tview.PrintSimple(screen, string(tcell.RuneLRCorner), x+LenWD(MainTabs[0].Text+MainTabs[1].Text)+5, y-1)
		tview.PrintSimple(screen, string(tcell.RuneLLCorner), x+LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text)+8, y-1)
		tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine), width-LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text)-10), x+LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text)+9, y-1)
		tview.PrintSimple(screen, string(tcell.RuneURCorner), width, y-1)
	case 3:
		tview.PrintSimple(screen, string(tcell.RuneULCorner), x-1, y-1)
		tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine), x+LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text)+7), x, y-1)
		tview.PrintSimple(screen, string(tcell.RuneLRCorner), x+LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text)+8, y-1)
		tview.PrintSimple(screen, string(tcell.RuneLLCorner), x+LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text+MainTabs[3].Text)+11, y-1)
		tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine), width-LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text+MainTabs[3].Text)-13), x+LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text+MainTabs[3].Text)+12, y-1)
		tview.PrintSimple(screen, string(tcell.RuneURCorner), width, y-1)
	case 4:
		tview.PrintSimple(screen, string(tcell.RuneULCorner), x-1, y-1)
		tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine), x+LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text+MainTabs[3].Text)+10), x, y-1)
		tview.PrintSimple(screen, string(tcell.RuneLRCorner), x+LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text+MainTabs[2].Text)+10, y-1)
		tview.PrintSimple(screen, string(tcell.RuneLLCorner), x+LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text+MainTabs[3].Text+MainTabs[4].Text)+14, y-1)
		tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine), width-LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text+MainTabs[3].Text+MainTabs[4].Text)-16), x+LenWD(MainTabs[0].Text+MainTabs[1].Text+MainTabs[2].Text+MainTabs[3].Text+MainTabs[4].Text)+15, y-1)
		tview.PrintSimple(screen, string(tcell.RuneURCorner), width, y-1)
	}

	//画一个边框
	for i := y; i < height+4; i++ {
		if i == height+3{
			tview.PrintSimple(screen, string(tcell.RuneLLCorner), x-1, i)
			tview.PrintSimple(screen, strings.Repeat(string(tcell.RuneHLine),width), x, i)
			tview.PrintSimple(screen, string(tcell.RuneLRCorner), width, i)
		}else{
			tview.PrintSimple(screen, string(tcell.RuneVLine), x-1, i)
			tview.PrintSimple(screen, string(tcell.RuneVLine), width, i)
		}

	}

	tview.PrintSimple(screen, string(tcell.RuneURCorner), width, y-1)

	x += 1     //x轴增加一点空白
	width -= 2 //行末留一点空白

	//写入之前重置一下
	//CleanArea(screen, x-2, y, width+3, height+4)

	//计算diplayoffset 和 itemoffset 位置
	if s.currentItem[s.currentTab] < height {
		s.displayOffset[s.currentTab] = s.currentItem[s.currentTab]
		s.itemOffset[s.currentTab] = 0
	} else {
		s.displayOffset[s.currentTab] = s.currentItem[s.currentTab] % height
		s.itemOffset[s.currentTab] = s.currentItem[s.currentTab] - s.displayOffset[s.currentTab]

	}

	for i := 0; i < s.GetItemCount(); i++ {
		if y >= bottonLimit {
			//log.Println("a")
			continue
		}
		if i < s.itemOffset[s.currentTab] {
			continue
		}
		//log.Println(i,y,bottonLimit)

		if y-3 == s.displayOffset[s.currentTab] {
			//高亮这行
			tview.Print(screen, string(tcell.RuneRArrow), x-1, y, width, tview.AlignLeft, tcell.GetColor("green"))
			tview.Print(screen, s.items[i].MainText, x, y, width, tview.AlignLeft, tcell.GetColor("green"))
			if s.items[i].RunStatus == "Failed"{ //如果runstatus是失败，则标红
				tview.Print(screen, s.items[i].RunStatus, width-22, y, width, tview.AlignLeft, tcell.GetColor("red"))
			}else{
				tview.Print(screen, s.items[i].RunStatus, width-22, y, width, tview.AlignLeft, tcell.GetColor("white"))
			}
			tview.Print(screen, s.items[i].PropStatus, width-8, y, width, tview.AlignLeft, tcell.GetColor("white"))
			tview.Print(screen, string(tcell.RuneBlock), width+1, y, width, tview.AlignLeft, tcell.GetColor("green"))
			y++
		} else {
			tview.Print(screen, s.items[i].MainText, x, y, width-25, tview.AlignLeft, tcell.GetColor("white"))
			if s.items[i].RunStatus == "Failed"{ //如果runstatus是失败，则标红
				tview.Print(screen, s.items[i].RunStatus, width-22, y, width, tview.AlignLeft, tcell.GetColor("red"))
			}else{
				tview.Print(screen, s.items[i].RunStatus, width-22, y, width, tview.AlignLeft, tcell.GetColor("white"))
			}
			tview.Print(screen, s.items[i].PropStatus, width-8, y, width, tview.AlignLeft, tcell.GetColor("white"))
			y++
		}

	}

}

func (s *SList) DownKey() {
	s.currentItem[s.currentTab]++
	if s.currentItem[s.currentTab] == s.GetItemCount() {
		s.currentItem[s.currentTab] = 0
	}
}
func (s *SList) UpKey() {
	s.currentItem[s.currentTab]--
	if s.currentItem[s.currentTab] == -1 {
		s.currentItem[s.currentTab] = s.GetItemCount() - 1
	}
}

func (s *SList) SetCurrentTab(tabID int) {
	s.currentTab = tabID
}

func (s *SList) Refresh(){
	s.Clear()
	for index, item := range unitDBS[s.currentTab].Items {
		FreezerState := fmt.Sprint(item.Status.ActiveState)
		UnitFileState := fmt.Sprint(item.Property["UnitFileState"])
		s.AddItem(index, item.Status.Name, item.Status.Description, FreezerState, UnitFileState,"")
	}
}

func (s *SList) NextPage(){
	if s.currentItem[s.currentTab] + s.height < s.GetItemCount(){
		s.currentItem[s.currentTab] = s.currentItem[s.currentTab] + s.height
	}else{
		s.currentItem[s.currentTab] = s.GetItemCount()-1
	}
}
func (s *SList) PreviousPage(){
	if s.currentItem[s.currentTab] - s.height > 0{
		s.currentItem[s.currentTab] = s.currentItem[s.currentTab]-s.height
	}else {
		s.currentItem[s.currentTab] = 0
	}
}
func (s *SList) HomePage(){
	s.currentItem[s.currentTab] = 0
}
func (s *SList) EndPage(){
	s.currentItem[s.currentTab] = s.GetItemCount()-1
}
func (s *SList) NextTab() {
	if s.currentTab == len(MainTabs)-1{
		s.currentTab = 0
	}else{
		s.currentTab ++
	}
}
func (s *SList) PreviousTab() {
	if s.currentTab == 0{
		s.currentTab = len(MainTabs)-1
	}else{
		s.currentTab --
	}
}

func LenWD(string2 string) int {
	return tview.TaggedStringWidth(string2)
}
