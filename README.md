# Create

createはコマンドでのファイル作成をより簡単に行います。

`create.exe hoge.txt`

で`hoge.txt`を作成します

また、

`create.exe ./piyo/hoge.txt`

で `piyo/hoge.txt`にファイルを作成します。`./piyo/`がない場合は作成します

### -p option

`-p`オプションはファイルのプロトタイプを指定し、作成します。

プロトタイプは`-init` オプションで指定したディレクトリから参照されます。

```
create.exe -init ./prototypes/
//プロトタイプのディレクトリを指定。
create.exe -p hoge
//プロトタイプディレクトリ内の/hogeディレクトリ内がプロトタイプとしてコピーされます。
```

