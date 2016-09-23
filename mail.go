package main

import (
	"bytes"
	"fmt"
	"math"
	"net/smtp"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

func SendMail(from string, to string, cc string, pwd string, title string, msg string) error {
	header := ""
	header += "From:" + from + "\n"
	header += "To:" + to + "\n"
	header += "Cc:" + cc + "\n"
	header += "Subject:" + title + "\n"
	header += "\n"

	msg = header + msg

	auth := smtp.PlainAuth("", from, pwd, "smtp.gmail.com")

	var recievers []string
	if cc != "" {
		recievers = append(strings.Split(to, ","), strings.Split(cc, ",")...)
	} else {
		recievers = strings.Split(to, ",")
	}

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, recievers, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func BuildMailBody(vol *VolumeInfo, list [][]*DirectoryInfo) (string, string) {
	title := "【容量確認のお知らせ】"

	msg := ""
	msg += "このメールは自動実行で送信しています。\n"
	msg += "残り容量が少なくなってきています。早めのバックアップをお願いいたします。\n"
	msg += "\n"
	msg += fmt.Sprintf("[ %v ]の現在状況\n", filepath.VolumeName(conf.Directory.ROOT_DIRECTORY))
	msg += fmt.Sprintf("総容量: %v\n", CalcByteToStr(vol.Total))
	msg += fmt.Sprintf("空き容量: %v\n", CalcByteToStr(vol.Free))
	msg += fmt.Sprintf("使用容量: %v\n", CalcByteToStr(vol.Used))
	msg += fmt.Sprintf("使用率: %v", vol.UsedPercent) + "%\n"
	msg += "\n"

	for i, dir := range strings.Split(conf.Directory.TARGET_DIRECTORIES, ",") {
		total := 0
		for _, d := range list[i] {
			total += d.Size
		}

		msg += "[ " + filepath.Join(conf.Directory.ROOT_DIRECTORY, dir) + " : " + CalcByteToStr(total) + " ] 使用状況順(容量x保存経過時間)\n"

		buf := new(bytes.Buffer)
		w := new(tabwriter.Writer)
		w.Init(buf, 0, 8, 0, '\t', 0)

		fmt.Fprintln(w, "順\tディレクトリ名\t容量\t更新日時")
		list[i] = list[i][:int(math.Min(float64(conf.Rank.MAX), float64(len(list[i]))))]
		for j, d := range list[i] {
			n := strconv.Itoa(j + 1)
			s := CalcByteToStr(d.Size)
			t := d.ModTime.Format("2006-01-02")

			fmt.Fprintln(w, n+"\t"+d.Name+"\t"+s+"\t"+t)
		}

		fmt.Fprintln(w)
		w.Flush()

		msg += buf.String()
		msg += "\n"
	}

	return title, msg
}
