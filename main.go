package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Prog は各番組の情報を格納する
type Prog struct {
	XMLName xml.Name `xml:"prog"`
	Title   string   `xml:"title"`
	Pfm     string   `xml:"pfm"`
}

// Radiko はRadiko週間番組表（局ごと）を格納する
type Radiko struct {
	XMLName  xml.Name `xml:"radiko"`
	Stations *struct {
		Station *struct {
			Name string `xml:"name"`
			Days []*struct {
				Date     string  `xml:"date"`
				Programs []*Prog `xml:"prog"`
			} `xml:"progs"`
		} `xml:"station"`
	} `xml:"stations"`
}

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	exitCode = 0

	client := http.Client{}
	resp, err := client.Get("https://radiko.jp/v3/program/station/weekly/TBS.xml")
	if err != nil {
		log.Printf("alert: Radiko番組表を読み込めませんでした：%s", err)
		exitCode = 1
		return
	}
	defer resp.Body.Close()

	var table Radiko
	if err := xml.NewDecoder(resp.Body).Decode(&table); err != nil {
		log.Printf("alert: xmlデータを構造体に読み込めませんでした：%v", err)
		exitCode = 1
		return
	}

	// 今日の該当番組のゲスト名抽出
	ts := (time.Now()).Format("20060102")
	guest := ""
	for _, day := range table.Stations.Station.Days {
		if strings.Contains(day.Date, ts) {
			for _, prog := range day.Programs {
				if strings.Contains(prog.Title, "伊集院光とらじおと") {
					guest = getGuest(prog)
					if guest != "" {
						break
					}
				}
			}
			break
		}
	}

	fmt.Printf("%s", guest)

	return
}

func getGuest(prog *Prog) (gst string) {
	pfm := prog.Pfm
	gsl := strings.Split(pfm, "ゲスト：")
	if len(gsl) == 1 {
		return
	}
	gst = gsl[len(gsl)-1]
	gst = strings.ReplaceAll(gst, "/", "／")

	return
}
