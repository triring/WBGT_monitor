package main

// 日本語の文字列は、ここにまとめておく。
// フォントデータを組込用データに変換する時に、容量を最低限にするために必要。
// 40pt表示文字列
const (
	txt_wbgt string = "暑さ指数"
	txt_thi  string = "不快指数"
	// wbgtレベル
	txt_wbgt1 string = "ほぼ安全"
	txt_wbgt2 string = "注意"
	txt_wbgt3 string = "警戒"
	txt_wbgt4 string = "厳重警戒"
	txt_wbgt5 string = "危険" // 危険の険が表示できないので、旧字を使用する。
)
