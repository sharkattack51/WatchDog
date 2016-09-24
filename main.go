package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"
)

const (
	CONF_FILE = "./watchdog.conf"
	ENC_FILE  = "./encrypt_pwd.data"
)

func main() {
	fmt.Println("\n")
	fmt.Println("[[[ WatchDog ]]]")

	// コマンドライン引数
	p := flag.String("pwd", "", "Encrypt mail password")
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
	conf,err := LoadConfig(CONF_FILE)
	if err != nil {
		log.Fatal(err)
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
	fmt.Println("\n")
	defer func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "- Finish.")
	}()

	// ボリューム容量を確認する
	vol, err := GetVolInfo(filepath.VolumeName(conf.Directory.ROOT_DIRECTORY))
	if err != nil {
		log.Fatal(err)
	}

	// 空き容量の閾値設定
	if vol.Free > conf.Volume.FREE_BYTE_TH {
		return
	}

	// 各ディレクトリの容量を確認してランキング付け
	rankedlist, err := RankingDirectories()
	if err != nil {
		log.Fatal(err)
	}

	// メールを送信する
	if conf.Mail.USE_SEND_MAIL {
		title, msg := BuildMailBody(vol, rankedlist)
		err := SendMail(conf.Mail.USER_NAME, conf.Mail.TO, conf.Mail.CC, conf.Mail.PASSWORD, title, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
