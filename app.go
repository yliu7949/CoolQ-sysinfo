package main

import (
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"strings"
)

//go:generate cqcfg -c .
// cqp: 名称: CoolQ-sysinfo
// cqp: 版本: 1.0.0:1
// cqp: 作者: Underworld
// cqp: 简介: 在QQ群里监控Linux服务器~
func main() { /*此处应当留空*/ }

func init() {
	cqp.AppID = "me.cqp.underworld.sysinfo" // TODO: 修改为这个插件的ID
	cqp.GroupMsg = onGroupMsg
	cqp.PrivateMsg = onPrivateMsg
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	defer handlePanic()
	if strings.HasPrefix(msg, "[CQ:at,qq=3********]") && msg[22:30] == "#sysinfo" {	//群机器人的QQ号
		reply := handleCmd()
		cqp.SendGroupMsg(fromGroup, reply)
	}
	return 0
}

func onPrivateMsg(subType, msgID int32, fromQQ int64, msg string, font int32) int32 {
	defer handlePanic()
	if msg == "#sysinfo" {
		reply := handleCmd()
		cqp.SendPrivateMsg(fromQQ, reply)
	}
	return 0
}

func handleCmd() string {
	text := "服务器的实时数据如下：\n"
	var cpuText string
	c, err := cpu.Times(false)		//CPU
	if err!= nil {
		return "出错啦！"
	}
	for _, elem := range c {
		cpuInfo := strconv.FormatFloat(100.0-elem.Idle/elem.Total()*100.0,'f',1,64)
		cpuText = "CPU的使用量为" + cpuInfo + "%\n"
	}
	m, err := mem.VirtualMemory()		//VirtualMemory
	if err!= nil {
		return "出错啦！"
	}
	memInfo := strconv.FormatFloat(m.UsedPercent,'f',1,64)
	memText := "物理内存的使用量为" + memInfo + "%\n"
	s, err := mem.SwapMemory()		//SwapMemory
	if err!= nil {
		return "出错啦！"
	}
	swapInfo := strconv.FormatFloat(s.UsedPercent,'f',1,64)
	swapText := "虚拟内存的使用量为" + swapInfo + "%\n"
	d,err := disk.Usage("/")		//Disk
	if err!= nil {
		return "出错啦！"
	}
	diskInfo := strconv.FormatFloat(d.UsedPercent,'f',1,64)
	diskText := "硬盘的使用量为" + diskInfo + "%"

	text =  text + cpuText + memText + swapText + diskText
	return text
}

func handlePanic() {
	if r := recover(); r != nil {
		cqp.AddLog(cqp.Error, "未知错误", fmt.Sprint(r))
		return
	}
}
