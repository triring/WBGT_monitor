# serial_logger

シリアルからの入力を監視し、受け取った文字列の先頭に、受信した日付と時刻を加えて、コンソールへの表示とログライルへの追記を行う。


## コンパイル

``` bash
$ go build -o serial_logger main.go
```

## 実行

* -baud オプションでボーレートを受け取る (デフォルト値: 115200)
* -file オプションで出力ファイル名を受け取る (デフォルト値: serial_log.txt)
* -port シリアルポート名を受け取る (デフォルト値: /dev/ttyUSB0)

## 使用例

``` bash
$ ./serial_logger -baud 115200 -port /dev/ttyUSB0 -file WBGT`date "+%Y%m%d"`.csv
```

上記は、Raspberry Pi上で使用することを想定した起動例である。

* シリアルポート/dev/ttyUSB0を監視する。
* 通信速度は、115200 bps に設定する。
* シリアルから送られてくる計測データは、その日の日付を埋め込んだファイルに追記していく。
