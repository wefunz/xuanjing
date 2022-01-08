package main

import (
	// "encoding/json"
	"fmt"
	"net/http"

	// "reflect"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Info struct {
	UsedMb      uint64  `json:"usedMb"`
	TotalMb     uint64  `json:"totalMb"`
	UsedPercent float64 `json:"usedPercent"`
}

func NewQuery(usedmb uint64, totalmb uint64, usedpercent float64) *Info {
	info := &Info{UsedMb: usedmb, TotalMb: totalmb, UsedPercent: usedpercent}
	return info
}

func main() {

	v, _ := mem.VirtualMemory()
	parts, _ := disk.Partitions(true)
	d, _ := disk.Usage(parts[0].Mountpoint)
	// almost every return value is a struct
	// fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	// convert to JSON. String() is also implemented
	queryMem := NewQuery(v.Total, v.Free, v.UsedPercent)
	queryDisk := NewQuery(d.Total, d.Free, d.UsedPercent)
	fmt.Println(d.Total, d.Free, d.UsedPercent)
	r := gin.Default()
	queryAll := make(map[string]interface{})
	//mem+disk
	r.GET("/status", func(c *gin.Context) {
		queryAll["disk"] = queryDisk
		queryAll["memory"] = queryMem
		c.IndentedJSON(http.StatusOK, queryAll)
	})
	//mem
	r.GET("/status/mem", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, queryMem)
	})
	//disk
	r.GET("/status/disk", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, queryDisk)
	})

	r.Run(":9090")

}
