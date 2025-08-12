package main

// go run main.go
/*
go get go.bug.st/serial
go mod tidy
go run main.go
*/

import (
	"bufio"
	"flag" // コマンドライン引数処理のためのパッケージ
	"fmt"
	"log"
	"os"
	"time"

	"go.bug.st/serial"
)

func main() {
	// --- コマンドライン引数の定義 ---
	// -port オプションでシリアルポート名を受け取る (デフォルト値: /dev/ttyUSB0)
	portName := flag.String("port", "/dev/ttyUSB0", "監視するシリアルポート名 (例: COM1, /dev/ttyUSB0)")
	// -file オプションで出力ファイル名を受け取る (デフォルト値: serial_log.txt)
	outputFilename := flag.String("file", "serial_log.txt", "受信データを保存するファイル名")
	// -baud オプションでボーレートを受け取る (デフォルト値: 9600)
	baudRate := flag.Int("baud", 115200, "シリアル通信のボーレート")

	flag.Parse() // コマンドライン引数をパース

	// 引数で取得した値をローカル変数に格納
	pName := *portName
	oFilename := *outputFilename
	bRate := *baudRate
	// --- コマンドライン引数の定義ここまで ---

	// シリアルポートを開く
	mode := &serial.Mode{
		BaudRate: bRate,
	}
	port, err := serial.Open(pName, mode)
	if err != nil {
		log.Fatalf("シリアルポート %s を開けませんでした: %v", pName, err)
	}
	defer port.Close() // アプリケーション終了時にポートを閉じる

	fmt.Printf("シリアルポート %s (ボーレート: %d) を開きました。\n", pName, bRate)
	fmt.Printf("受信データを %s に追記します。\n", oFilename)
	fmt.Println("Ctrl+C で終了します。")

	// 出力ファイルを開く (追記モード)
	outputFile, err := os.OpenFile(oFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("出力ファイル %s を開けませんでした: %v", oFilename, err)
	}
	defer outputFile.Close() // アプリケーション終了時にファイルを閉じる

	// 受信バッファ
	scanner := bufio.NewScanner(port)

	// 受信ループ
	for scanner.Scan() {
		receivedData := scanner.Text()
		// timestamp := time.Now().Format("2006-01-02,15:04:05") // タイムスタンプ
		// timestamp := time.Now().Format(time.RFC1123Z) // タイムスタンプ
		timestamp := time.Now() // タイムスタンプ
		logEntry := fmt.Sprintf("%s,%s,%s\r\n",
			timestamp.Format(time.DateOnly),
			timestamp.Format(time.TimeOnly),
			receivedData)

		fmt.Print(logEntry) // コンソールにも表示

		// ファイルに追記
		if _, err := outputFile.WriteString(logEntry); err != nil {
			log.Printf("ファイルへの書き込みに失敗しました: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("シリアルポートからの読み込みエラー: %v", err)
	}
}
