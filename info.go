package main

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type VolumeInfo struct {
	Total       int
	Free        int
	Used        int
	UsedPercent float64
}

type DirectoryInfo struct {
	Name    string
	Path    string
	Size    int
	ModTime time.Time
	Owner   string
}

// ボリューム情報を取得
func GetVolInfo(volume string) (*VolumeInfo, error) {
	u, err := disk.Usage(filepath.VolumeName(volume))
	if err != nil {
		return nil, err
	}

	total := int(u.Total)
	free := int(u.Free)
	used := int(u.Used)
	usedp := math.Floor(u.UsedPercent*10.0) / 10.0

	fmt.Println("[ " + volume + " ]")
	fmt.Println("Total :", CalcByteToStr(total))
	fmt.Println("Free :", CalcByteToStr(free))
	fmt.Println("Used :", CalcByteToStr(used))
	fmt.Println("UsedPercent :", usedp, "%")
	fmt.Println("\n")

	return &VolumeInfo{total, free, used, usedp}, nil
}

// ランキングソート
func RankingDirectories() ([][]*DirectoryInfo, error) {
	var rankedlist [][]*DirectoryInfo
	for _, t := range strings.Split(conf.Directory.TARGET_DIRECTORIES, ",") {
		dirs := MakeDirInfoList(filepath.Join(conf.Directory.ROOT_DIRECTORY, t))
		list := SortList(<-dirs)
		sort.Sort(list)
		rankedlist = append(rankedlist, []*DirectoryInfo(list))

		fmt.Println("[ " + filepath.Join(conf.Directory.ROOT_DIRECTORY, t) + " ]")
		for i, d := range list {
			fmt.Println(i+1, d.Name, CalcByteToStr(d.Size), d.ModTime.Format("2006-01-02"))
		}
		fmt.Println("\n")
	}

	return rankedlist, nil
}

// ディレクトリ情報を作成
func MakeDirInfoList(root string) <-chan []*DirectoryInfo {
	ch := make(chan []*DirectoryInfo)
	var list []*DirectoryInfo

	fmt.Print(">")
	go func() {
		var wg sync.WaitGroup

		dirs, _ := GetProjDirList(root)
		for _, d := range dirs {
			wg.Add(1)
			go func(f os.FileInfo) {
				path := filepath.Join(root, f.Name())
				size, _ := GetDirSize(path)
				owner := GetDirOwner(path)

				info := &DirectoryInfo{f.Name(), path, size, f.ModTime(), owner}
				list = append(list, info)

				wg.Done()
				fmt.Print("*")
			}(d)
		}

		wg.Wait()
		fmt.Print("\n\n")

		ch <- list
	}()

	return ch
}

// プロジェクトディレクトリ一覧を取得
func GetProjDirList(root string) ([]os.FileInfo, error) {
	var dirs []os.FileInfo

	fileInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, f := range fileInfos {
		if f.IsDir() {
			ck := true
			for _, n := range strings.Split(conf.Directory.IGNORE_DIRECTORIES, ",") {
				if f.Name() == n {
					ck = false
					break
				}
			}

			if ck {
				dirs = append(dirs, f)
			}
		}
	}

	return dirs, nil
}

// ディレクトリ容量を取得
func GetDirSize(path string) (int, error) {
	size := 0

	err := filepath.Walk(path,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			info, err = os.Stat(p)
			if err != nil {
				return nil
			}

			size += int(info.Size())
			return nil
		})

	if err != nil {
		return size, err
	}

	return size, nil
}

// ディレクトリオーナーを取得
func GetDirOwner(path string) string {
	//out, _ := exec.Command("dir", "/q", path).Output()
	return ""
}

// バイト表示に変換
func CalcByteToStr(size int) string {
	u := 1.0
	s := "B"
	if size > 1024*1024*1024*1024 {
		u = 1024 * 1024 * 1024 * 1024
		s = "TB"
	} else if size > 1024*1024*1024 {
		u = 1024 * 1024 * 1024
		s = "GB"
	} else if size > 1024*1024 {
		u = 1024 * 1024
		s = "MB"
	} else {
		u = 1024
		s = "KB"
	}
	c := math.Floor((float64(size)/u*10.0)+0.5) / 10.0

	return strconv.FormatFloat(c, 'f', 1, 64) + " " + s
}
