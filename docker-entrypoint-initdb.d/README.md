# docker-entrypoint-initdb.d

以下のコマンドを実行して、マイグレーションファイルを生成してください。

```bash
touch docker-entrypoint-initdb.d/"$(date +%Y%m%d%H%M%S)_{マイグレーションファイルの内容}.js"
```
