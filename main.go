package main

import (
	"flag"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"path/filepath"
	"time"
)

const (
	CONF_FILE = "./watchdog.conf"
	ENC_FILE  = "./encrypt_pwd.data"
	DB_FILE   = "./level.db"
)

var conf *Config
var db *leveldb.DB

func main() {
	timescore := time.Now()
	defer func() {
		fmt.Println("score:", time.Now().Sub(timescore).Seconds())
	}()

	fmt.Println("[[[ WatchDog ]]]")

	// コマンドライン引数
	p := flag.String("pwd", "", "Encrypt mail password")
	c := flag.String("config", "", "Config file")
	flag.Parse()

	// Mailパスワードを暗号化して保存
	if *p != "" {
		err := EncryptToFile(ENC_FILE, *p)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "- Encrypted mail password to", ENC_FILE)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// 設定ファイルの読み込み
	if *c != "" {
		conf = LoadConfig(*c)
	} else {
		conf = LoadConfig(CONF_FILE)
	}

	if conf.Mail.PASSWORD == "" {
		// パスワードの複合化
		pwd, err := DecryptFromFile(ENC_FILE)
		if err != nil {
			log.Fatal(err)
		}
		conf.Mail.PASSWORD = pwd
	}

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "- StartCalculate...")
	fmt.Print("\n")
	defer func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "- Finish.")
		fmt.Print("\n")
	}()

	// ボリューム容量を確認する
	vol, err := GetVolInfo(filepath.VolumeName(conf.Directory.ROOT_DIRECTORY))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[ " + vol.Name + " ]")
	fmt.Println("Total :", CalcByteToStr(vol.Total))
	fmt.Println("Free :", CalcByteToStr(vol.Free))
	fmt.Println("Used :", CalcByteToStr(vol.Used))
	fmt.Println("UsedPercent :", vol.UsedPercent, "%")
	fmt.Print("\n")

	// 空き容量の閾値設定
	if vol.Free > conf.Volume.FREE_BYTE_TH {
		return
	}

	// キャッシュ確認用DB
	db, err = OpenDB(DB_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 各ディレクトリの容量を確認してランキング付け
	rankedlist, err := RankingDirectories()
	if err != nil {
		log.Fatal(err)
	}
	for _, list := range rankedlist {
		fmt.Print("\n")
		fmt.Println("[ " + filepath.Dir(list[0].Path) + " ]")
		for i, d := range list {
			fmt.Println(i+1, d.Name, CalcByteToStr(d.Size), d.ModTime.Format("2006-01-02"), GetDirOwner(d.Path))
		}
	}
	fmt.Print("\n")

	// メールを送信する
	if conf.Mail.USE_SEND_MAIL {
		title, msg := BuildMailBody(vol, rankedlist)
		err := SendMail(conf.Mail.USER_NAME, conf.Mail.TO, conf.Mail.CC, conf.Mail.PASSWORD, title, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
