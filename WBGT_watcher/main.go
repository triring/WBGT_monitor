// tinygo flash  -target=m5stack -size=short -monitor .

// go mod download tinygo.org/x/tinyfont
// go get tinygo.org/x/drivers/ili9341/initdisplay
// 以下のI2C接続のリアルタイムクリック RCT を追加しようとした。
// 	"tinygo.org/x/drivers/pcf8563"
// しかし、温湿度気圧センサーのI2C通信を競合するのか、機能しなくなった。
// それぞれ、単品で動かすと、問題ない。
// 2つのセンサーを同じI2Cに接続すると動かない。

package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/bmp280"
	"tinygo.org/x/drivers/examples/ili9341/initdisplay"
	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/drivers/sht4x"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
)

// カラーユニバーサルデザイン(CUD) カラーセット
var (
	// Accent Colors アクセントカラー
	red      = color.RGBA{R: 0xFF, G: 0x4B, B: 0x0, A: 0xFF}  //  Red : 赤
	yellow   = color.RGBA{R: 0xFF, G: 0xF1, B: 0x0, A: 0xFF}  //  Yellow : 黄色
	green    = color.RGBA{R: 0x3, G: 0xAF, B: 0x7A, A: 0xFF}  //  Green : 緑
	blue     = color.RGBA{R: 0x0, G: 0x5A, B: 0xFF, A: 0xFF}  //  Blue : 青
	sky_blue = color.RGBA{R: 0x4D, G: 0xC4, B: 0xFF, A: 0xFF} //  Sky blue : 空色
	pink     = color.RGBA{R: 0xFF, G: 0x80, B: 0x82, A: 0xFF} //  Pink : ピンク
	orange   = color.RGBA{R: 0xF6, G: 0xAA, B: 0x0, A: 0xFF}  //  Orange : オレンジ
	purple   = color.RGBA{R: 0x99, G: 0x0, B: 0x99, A: 0xFF}  //  Purple : 紫
	brown    = color.RGBA{R: 0x80, G: 0x40, B: 0x0, A: 0xFF}  //  Brown : 茶色

	// Base Colors  ベースカラー
	light_pink         = color.RGBA{R: 0xFF, G: 0xCA, B: 0xBF, A: 0xFF} //  Light pink : 明るいピンク
	cream              = color.RGBA{R: 0xFF, G: 0xFF, B: 0x80, A: 0xFF} //  Cream : クリーム
	light_yellow_green = color.RGBA{R: 0xD8, G: 0xF2, B: 0x55, A: 0xFF} //  Light yellow-green : 明るい黄緑
	light_sky_blue     = color.RGBA{R: 0xBF, G: 0xE4, B: 0xFF, A: 0xFF} //  Light sky blue : 明るい空色
	beige              = color.RGBA{R: 0xFF, G: 0xCA, B: 0x80, A: 0xFF} //  Beige : ベージュ
	light_green        = color.RGBA{R: 0x77, G: 0xD9, B: 0xA8, A: 0xFF} //  Light green : 明るい緑
	light_purple       = color.RGBA{R: 0xC9, G: 0xAC, B: 0xE6, A: 0xFF} //  Light purple : 明るい紫

	// Achromatic Colors 無彩色
	white      = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} //  White  白
	light_gray = color.RGBA{R: 0xC8, G: 0xC8, B: 0xCB, A: 0xFF} //  Light gray  明るいグレー
	gray       = color.RGBA{R: 0x84, G: 0x91, B: 0x9E, A: 0xFF} //  Gray  グレー
	black      = color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xFF}    //  Black  黒
)

var (
	display *ili9341.Device
)

// Unix時間への相互変換（ミリ秒）

// システム時間からUnix時間への変換
func timeToUnixMilli(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

// Unix時間からシステム時間への変換
func unixMilliToTime(millis int64) time.Time {
	return time.Unix(0, millis*1000000)
}

func main() {
	display = initdisplay.InitDisplay()

	width, height := display.Size()
	if width < 320 || height < 240 {
		display.SetRotation(ili9341.Rotation270)
	}

	display.FillScreen(light_gray)
	time.Sleep(5 * time.Second)

	machine.I2C0.Configure(machine.I2CConfig{})
	sensorPressure := bmp280.New(machine.I2C0)
	sensorPressure.Address = 0x76
	sensorPressure.Configure(bmp280.STANDBY_125MS, bmp280.FILTER_4X, bmp280.SAMPLING_16X, bmp280.SAMPLING_16X, bmp280.MODE_FORCED)

	connected := sensorPressure.Connected()
	if !connected {
		println("\nBMP280 Sensor not detected\n")
		return
	}
	println("\nBMP280 Sensor detected\n")

	println("Calibration:")
	sensorPressure.PrintCali()

	machine.I2C0.Configure(machine.I2CConfig{})
	sensorTempHum := sht4x.New(machine.I2C0)
	var prev int = 0
	for {
		temp, humidity, _ := sensorTempHum.ReadTemperatureHumidity()
		h := float32(humidity) / 1000.0
		t := float32(temp) / 1000.0
		// fmt.Printf("気温(℃): %2.2f\n", t) //	Temperature
		// fmt.Printf("湿度(％): %2.2f\n", h) //	humidity
		pres, _ := sensorPressure.ReadPressure()
		p := float32(pres) / 100000.0
		// Pressure in hectoPascal
		// fmt.Printf("気圧(hPa): %4.2f \n", p) // Pressure

		WBGT := (t*0.003289+0.01844)*h + (0.6868*t - 2.022) //暑さ指数の計算
		THI := (t*0.81 + h*0.01*(t*0.99-14.3) + 46.3)       //不快指数の計算

		// 不快指数は、温湿指数（略称THI）とも言われる。
		// 人間が生活するうえで不快を感じるような体感を、気温と湿度で表す指数
		// 測定結果を温度、湿度、気圧、暑さ指数、不快指数 の順でシリアルに出力する。
		fmt.Printf("%05.2f,%05.2f,%05.2f,%05.2f,%05.2f\n", t, h, p, WBGT, THI)

		// 以下は、液晶表示の更新
		str_temp := fmt.Sprintf("%10s:  %6.2f", txt_temp, t)  //  "気温(℃)" jp-24pt_string.go 内で定義
		str_hum := fmt.Sprintf("%10s:  %6.2f", txt_hum, h)    //  "湿度(％)" jp-24pt_string.go 内で定義
		str_pres := fmt.Sprintf("%10s:%6.2f", txt_pres, p)    //  "気圧(hPa)" jp-24pt_string.go 内で定義
		str_wbgt := fmt.Sprintf("%s : %4.1f", txt_wbgt, WBGT) // "暑さ指数"  jp-40pt_string.go 内で定義
		str_thi := fmt.Sprintf("%s : %4.1f", txt_thi, THI)    // "不快指数"  jp-40pt_string.go 内で定義

		// 利用可能なフォント
		// Bold9pt7b,Bold12pt7b,Bold18pt7b,Bold24pt7b
		// BoldOblique9pt7b,BoldOblique12pt7b,BoldOblique18pt7b,BoldOblique24pt7b
		// Oblique9pt7b,Oblique12pt7b,Oblique18pt7b,Oblique24pt7b
		// Regular9pt7b,Regular12pt7b,Regular18pt7b,Regular24pt7b

		// 暑さ指数によるWBGTレベルの判定
		level := 1 // ほぼ安全	21未満
		if WBGT >= 21 {
			level = 2
		} // 注意	21以上25未満
		if WBGT >= 25 {
			level = 3
		} // 警戒	25以上28未満
		if WBGT >= 28 {
			level = 4
		} // 厳重警戒	28以上31未満
		if WBGT >= 31 {
			level = 5
		} // 危険	31以上

		switch level {
		case 1:
			tinydraw.FilledRectangle(display, 0, 0, 320, 48, blue)              // Draw yellow rectangle
			tinyfont.WriteLine(display, &Notosans40pt, 0, 42, txt_wbgt1, white) // "ほぼ安全"    jp-40pt_string.go 内で定義
		case 2:
			tinydraw.FilledRectangle(display, 0, 0, 320, 48, sky_blue)          // Draw yellow rectangle
			tinyfont.WriteLine(display, &Notosans40pt, 0, 42, txt_wbgt2, white) // "注意"        jp-40pt_string.go 内で定義
		case 3:
			tinydraw.FilledRectangle(display, 0, 0, 320, 48, yellow)            // Draw yellow rectangle
			tinyfont.WriteLine(display, &Notosans40pt, 0, 42, txt_wbgt3, black) // "警戒"        jp-40pt_string.go 内で定義
		case 4:
			tinydraw.FilledRectangle(display, 0, 0, 320, 48, orange)            // Draw yellow rectangle
			tinyfont.WriteLine(display, &Notosans40pt, 0, 42, txt_wbgt4, black) // "厳重警戒"    jp-40pt_string.go 内で定義
		case 5:
			tinydraw.FilledRectangle(display, 0, 0, 320, 48, red)               // Draw yellow rectangle
			tinyfont.WriteLine(display, &Notosans40pt, 0, 42, txt_wbgt5, black) // "危險"  jp-40pt_string.go 内で定義
			// 危険の険が表示できないので、旧字を使用する。
		}

		tinydraw.FilledRectangle(display, 0, 48, 320, 46, light_gray) // 暑さ指数の表示
		tinyfont.WriteLine(display, &Notosans40pt, 0, 86, str_wbgt, black)
		tinydraw.FilledRectangle(display, 0, 91, 320, 46, light_gray) // 不快指数の表示
		tinyfont.WriteLine(display, &Notosans40pt, 0, 132, str_thi, black)

		tinydraw.FilledRectangle(display, 0, 136, 320, 34, light_gray) // 気温の表示
		tinyfont.WriteLine(display, &Notosans24pt, 72, 168, str_temp, black)
		tinydraw.FilledRectangle(display, 0, 170, 320, 34, light_gray) // 湿度の表示
		tinyfont.WriteLine(display, &Notosans24pt, 72, 202, str_hum, black)
		tinydraw.FilledRectangle(display, 0, 204, 320, 36, light_gray) // 気圧の表示
		tinyfont.WriteLine(display, &Notosans24pt, 72, 235, str_pres, black)

		// 1分間のインターバル
		for {
			log_time := time.Now()
			if prev != log_time.Minute() {
				prev = log_time.Minute()
				break
			}
			time.Sleep(time.Second)
		}
	}
}
