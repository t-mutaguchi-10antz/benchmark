# SQLite パフォーマンス検証

## 実行環境

```
MacBook Pro
2.7 GHz クアッドコア Intel Core i7
メモリ 16 GB 2133 MHz LPDDR3
```

## 結果

1 レコードの読み込みを 100,000 回実行、1 回辺りの平均時間

#### 1. ヒープメモリ

```go
_ = Heap[id]
```

```
0.000000052 sec/load [heap memory]
```

#### 2. Go 標準パッケージ

```go
_ = db.QueryRow(query, id).Scan(&sample.ID, &sample.Field1, &sample.Field2, &sample.Field3)
```

```
0.000035458 sec/load [database/sql]
```

#### 3. 3rd Party の ORM パッケージ

```go
_, err := boiler.FindSample(ctx, db, null.NewString(id, true))
```

```
0.000039202 sec/load [volatiletech/sqlboiler ( ORM )]
```

## 考察

当然、ヒープメモリに生成済みのインスタンスへアクセスするのは圧倒的に高速。
ORM パッケージは色々あるが、標準パッケージとほぼ変わらない速度が出ているものもある。

- SQL を使わずにヒープメモリ ( マップ ) を参照した方が 1,000 倍近く高速
- SQL ( ORM ) を使った場合でも 30 〜 40 マイクロ秒なので高速
- SQL が使えると開発者体験は高く、場合に依って開発効率は上がる

など、それぞれの方法にはメリット・デメリットがあり、どちらが正解という話でも無いため、プロジェクト側で適宜判断すれば良いのではないか？