# PTCBomTreeView

## 使い方
1. CreoのBOMフォーマット (コンフィギュレーションの *bom_format*)に *bom_format.fmt* を指定する
2. Creoで部品表を出力する
3. 出力したbomファイルをCSVに変換する
```sh
node main.mjs hogehoge.bom.1 > out.csv
```

## 実行オプション
-t : CSVではなくTSV形式で出力する

-e : Assemblyの末端で空行を挿入する
