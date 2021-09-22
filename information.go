package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)


type Information struct {
	*tview.Box
	Service string
	Info map[string]string

}

type InformationWidge struct {
	page *tview.Flex
	info *Information
	menu *Menu
	showLog bool
}

func (iw *InformationWidge) Set(serverName string, infos UnitItem) {
	iw.info.Set(serverName,infos)
}

func (i *Information) Draw(screen tcell.Screen) {
	x, y, width, height := i.GetInnerRect()
	log.Println(x,y,width,height)
	//CleanArea(screen,x,y,width,height)
	//tview.PrintSimple(screen,i.Service,x,y)
	//画一个框框，加上框框标题
	for p := 0; p < height; p++ {
		if p == 0{
			if width-LenWD(i.Service)-4 > 0{

			}
			tview.PrintSimple(screen,string(tcell.RuneULCorner)+" ",x,y)
			tview.PrintSimple(screen,i.Service,x+2,y )
			if width-LenWD(i.Service)-4 > 0{
				tview.PrintSimple(screen,strings.Repeat(string(tcell.RuneHLine),width-LenWD(i.Service)-4),x+3+LenWD(i.Service),y)
			}else{
				tview.PrintSimple(screen,strings.Repeat(string(tcell.RuneHLine),1),x+3+LenWD(i.Service),y)
			}
			tview.PrintSimple(screen,string(tcell.RuneURCorner),width-1,y)
		} else{
			if p == height-1{
				tview.PrintSimple(screen,string(tcell.RuneLLCorner),x,p)
				tview.PrintSimple(screen,strings.Repeat(string(tcell.RuneHLine),width),1,p)
				tview.PrintSimple(screen,string(tcell.RuneLRCorner),width-1,p)
			}else{
				tview.PrintSimple(screen,string(tcell.RuneVLine),0,p)
				tview.PrintSimple(screen,string(tcell.RuneVLine),width-1,p)
			}

		}


	}
	//开始展示这个服务的主要信息了，需要根据不同的类别，去展示不同的数据
	//Desc  Status [time\autostart]  Path  Memory  socket和timer的Triggers
	//还要增加几个功能按钮 Restart Stop Mask/Unmask Edit show relationships
	//这些按钮还需要用tab键切换

	//tview.PrintSimple(screen,"*" + i.Info["Description"] + "*",x+1,y+1)
	m := 2
	keys := []string{}
	for k,_:= range i.Info{
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		tview.PrintSimple(screen, k+ "\t: " + i.Info[k],x+2,y+m)
		m++
	}
}

func (i *Information)Set(serverName string, infos UnitItem){
	i.Service = serverName
	i.Info = make(map[string]string)
	i.Info["Description"] = infos.Status.Description
	i.Info["Status"] = fmt.Sprint(infos.Status.ActiveState)
	i.Info["UnitFileState"] = fmt.Sprint(infos.Property["UnitFileState"])
	i.Info["Path"] = fmt.Sprint(infos.Property["FragmentPath"])
	for _,v := range MainTabs{
		if infos.TypeID == v.id{
			if infos.TypeID == 0{
				//说明是server unit
				i.Info["ExecMainPID"] = fmt.Sprint(infos.TypeProperty["ExecMainPID"])
				i.Info["MemoryCurrent"] = fmt.Sprint(infos.TypeProperty["MemoryCurrent"])
				i.Info["ExecStart"] = fmt.Sprint(infos.TypeProperty["ExecStart"])
				i.Info["ControlGroup"] = fmt.Sprint(infos.TypeProperty["ControlGroup"])
				i.Info["WorkingDirectory"] = fmt.Sprint(infos.TypeProperty["WorkingDirectory"])
				i.Info["User"] = fmt.Sprint(infos.TypeProperty["User"])
				i.Info["Group"] = fmt.Sprint(infos.TypeProperty["Group"])
				i.Info["Environment"] = fmt.Sprint(infos.TypeProperty["Environment"])
				i.Info["Result"] = fmt.Sprint(infos.TypeProperty["Result"])
				i.Info["Restart"] = fmt.Sprint(infos.TypeProperty["Restart"])
				i.Info["KillMode"] = fmt.Sprint(infos.TypeProperty["KillMode"])
				i.Info["PIDFile"] = fmt.Sprint(infos.TypeProperty["PIDFile"])
			}
			if infos.TypeID == 2{
				i.Info["Triggers"] = fmt.Sprint(infos.Property["Triggers"])
				i.Info["Listen"] =fmt.Sprint(infos.TypeProperty["Listen"])
				i.Info["SocketGroup"] =fmt.Sprint(infos.TypeProperty["SocketGroup"])
				i.Info["SocketUser"] =fmt.Sprint(infos.TypeProperty["SocketUser"])
			}
			if infos.TypeID == 3{
				i.Info["Triggers"] = fmt.Sprint(infos.Property["Triggers"])
				i.Info["Unit"] =fmt.Sprint(infos.TypeProperty["Unit"])
				i.Info["Result"] =fmt.Sprint(infos.TypeProperty["Result"])
				if nil != infos.TypeProperty["LastTriggerUSec"]{
					a,err := strconv.Atoi(fmt.Sprintf("%v", infos.TypeProperty["LastTriggerUSec"])[:10])
					if err!= nil{
						log.Println("typeID: 3",err)
						continue
					}
					log.Println(a)
					t:=time.Unix(int64(a),0)
					i.Info["LastTriggerUSec"] =fmt.Sprint(t.String())
				}



			}
			if infos.TypeID ==4 {
				i.Info["Where"] = fmt.Sprint(infos.TypeProperty["Where"])
				i.Info["What"] = fmt.Sprint(infos.TypeProperty["What"])
				i.Info["Options"] = fmt.Sprint(infos.TypeProperty["Options"])
				i.Info["Result"] = fmt.Sprint(infos.TypeProperty["Result"])
				i.Info["Type"] = fmt.Sprint(infos.TypeProperty["Type"])

			}
		}
	}

	keys := []string{}
	for k,_:= range infos.TypeProperty{
		keys = append(keys, k)
	}
	sort.Strings(keys)
	log.Println("==============")
	for _,v := range keys{
		log.Printf("%s : %s",v, infos.TypeProperty[v])
	}
}

func NewInformationWidge() *InformationWidge {
	info := &Information{
		Box: tview.NewBox(),
	}

	info_menu := NewMenu()
	info_menu.AddItem("Help", tcell.KeyF1, nil)
	info_menu.AddItem("Refresh", tcell.KeyF2,nil)
	info_menu.AddItem("Reload",tcell.KeyF3,nil)
	info_menu.AddItem("Disable",tcell.KeyF4,nil)
	info_menu.AddItem("Create", tcell.KeyF5,nil)
	aflex := tview.NewFlex().SetDirection(tview.FlexRow)
	aflex.AddItem(info,0,1,true)
	aflex.AddItem(info_menu,1,1,false)
	widge := &InformationWidge{
		page: aflex,
		info: info,
		menu: info_menu,
	}
	return widge
}
