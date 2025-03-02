# Go Batch MongoDB Aggregate

記事（未公開）：<br>
https://ap-ep.com/mongodb-aggregation-pipeline-vs-go/

## 背景・目的
集計処理を行うとき、GoのロジックとMongoDBのアグリゲーションパイプラインのどちらが高速かを調べることを目的とします。

## 実行方法
```zsh
$ make benchmark
```

## 測定内容
### 測定対象の概要
以下の処理を行います。
1. ユーザーとユーザーのポイントを管理するコレクションを用意する
2. アグリゲーションパイプライン or Goのロジックを使ってユーザーのポイントを集計し、集計結果をコレクションに格納する

アグリゲーションパイプラインとGoのロジックではどちらがよりパフォーマンスが優れているかを計測します。

## 計測結果
[benchmark_results.txt](https://github.com/taako-502/go-batch-mongodb-aggregate/blob/main/benchmark_results.txt) を参照してください。

### グラフ
グラフは以下です。
![スコアの集計処理の速度測定](https://github.com/taako-502/go-batch-mongodb-aggregate/assets/36348377/3fcf50d8-5c0f-4b98-9d95-c7e464579035)
