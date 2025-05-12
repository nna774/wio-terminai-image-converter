# wio-terminai-image-converter

https://wiki.seeedstudio.com/Wio-Terminal-LCD-Loading-Image/#step-2-convert-the-24-bit-bitmap-image-to-the-microcontroller-readable-8-bit-or-16-bit-bmp--format にある bmp_converter.py 同等のプログラム。
実行環境を用意するのがめんどくさかったので書いた。

## run

```bash
go run ./main.go rgb565 ./1012_240x.png out.bmp
```

main.goの中で `_ "image/png"` しかimportしていないので、png以外の画像形式は変換できないが、これを足せば他の形式にも対応できるはず。

なんかそのうち整える……。
