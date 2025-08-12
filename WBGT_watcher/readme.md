# tinygoで液晶に日本語を表示する手順

現在は、tinygoは、直接、液晶画面に多バイト文字のフォントを表示することはできない。
また、組み込みの場合、巨大なフォントデータを埋め込んでおいても、そのごく一部しか使用されることがないので、不合理である。

https://github.com/tinygo-org/tinyfont/tree/release/examples/unicode_font3_const2bit

そこで、上記のtinygo本家に書かれていた多バイト文字の表示手順を参考にして、
M5Stackの液晶画面に日本語を表示してみた。

1. 作業用ディレクトリの作成とモジュール管理ファイルの生成

```bash
mkdir UnicodeTest
cd UnicodeTest
go mod init UnicodeTest
```

2. 事前にUnicode_fontをダウンロードしておく。

    https://fonts.google.com/noto/specimen/Noto+Sans+JP
    https://fonts.google.com/noto/specimen/Noto+Sans+KR
        NotoSansJP-Regular.ttf
        NotoSansKR-Regular.ttf

3. テストディレクトリを作成し、以下のURLより多バイトの文字列を含むソースコードをダウンロードしてくる。

https://github.com/tinygo-org/tinyfont/blob/release/examples/unicode_font3_const2bit/main.go

```bash
mkdir test
cd test
```

4. フォントデータの変換ツールを導入する。

```bash
> go install tinygo.org/x/tinyfont/cmd/tinyfontgen@latest
go: downloading github.com/hajimehoshi/go-jisx0208 v1.0.0
go: downloading github.com/sago35/go-bdf v0.0.0-20200313142241-6c17821c91c4
go: downloading golang.org/x/image v0.0.0-20220617043117-41969df76e82

> go install tinygo.org/x/tinyfont/cmd/tinyfontgen-ttf@latest
go: downloading golang.org/x/text v0.3.7 
```

5. 多バイトの文字列を含むソースコードから、使用されているフォントだけのイメージデータをフォントファイルから抽出し、
その一覧をデータ化したgoのソースコードを作成する。

C:\Users\089241\go\bin\tinyfontgen-ttf.exe

```bash
> dir jp*.go
Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----        2025/08/06      9:01            332 jp-24pt_string.go
-a----        2025/08/06      9:00            484 jp-40pt_string.go
```

```bash
> tinyfontgen-ttf.exe --size 24 --verbose --output ./jp_font24.go --string-file ./jp-24pt_string.go --package main --fontname Notosans24pt ./NotoSansJP-Regular.ttf ./NotoSansKR-Regular.ttf
> tinyfontgen-ttf.exe --size 40 --verbose --output ./jp_font40.go --string-file ./jp-40pt_string.go --package main --fontname Notosans40pt ./NotoSansJP-Regular.ttf ./NotoSansKR-Regular.ttf

> dir jp*.go
Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----        2025/08/06      9:01            332 jp-24pt_string.go
-a----        2025/08/06      9:00            484 jp-40pt_string.go
-a----        2025/08/06      9:16         169288 jp_font24.go
-a----        2025/08/06      9:17         424931 jp_font40.go
```

6. main.goのコードをコンパイルを行う。
最初に、モジュールが足りないという警告がでたので、指摘されたモジュールを追加インストールした。

``` bash
go get tinygo.org/x/tinyfont/const2bit
```

``` bash
> tinygo build -o UnicodeTest.uf2 -target=wioterminal -size short ./test
test\font.go:6:2: no required module provides package tinygo.org/x/tinyfont/const2bit; to add it:
        go get tinygo.org/x/tinyfont/const2bit
test\main.go:7:2: no required module provides package tinygo.org/x/tinyfont; to add it:
        go get tinygo.org/x/tinyfont
test\main.go:8:2: no required module provides package tinygo.org/x/tinyfont/examples/initdisplay; to add it:
        go get tinygo.org/x/tinyfont/examples/initdisplay

> go get tinygo.org/x/tinyfont
go: added github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
go: added tinygo.org/x/drivers v0.29.0
go: added tinygo.org/x/tinyfont v0.5.0

> go get tinygo.org/x/tinyfont/examples/initdisplay
```

7. 再度、コンパイルを行った。無事に、バイナリファイルを生成できた。

```bash
> tinygo build -o UnicodeTest.uf2 -target=wioterminal -size short ./test
   code    data     bss |   flash     ram
  25352     296    6688 |   25648    6984

> dir UnicodeTest.uf2

Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----        2025/01/08     22:03          51712 UnicodeTest.uf2

```

7. UnicodeTest.uf2ファイルを Wio Terminal に転送し、多バイトコードが表示されることを確認した。
