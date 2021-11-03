# Better create file

Better create file はコマンドでのファイル作成をより簡単に行います。

`betterCreateFile.exe hoge.txt`

で`hoge.txt`を作成します

また、

`betterCreateFile.exe ./piyo/hoge.txt`

で `piyo/hoge.txt`にファイルを作成します。`./piyo/`がない場合は作成します

### -p option

`-p`オプションはファイルのプロトタイプを指定し、作成します。

プロトタイプは`-init` オプションで指定したディレクトリから参照されます。

```
betterCreateFile.exe -init ./prototypes/
//プロトタイプのディレクトリを指定。
betterCreateFile.exe -p hoge
//プロトタイプディレクトリ内の/hogeディレクトリ内がプロトタイプとしてコピーされます。
```

