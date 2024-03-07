.PHONY: run

run:
	go run main.go

benchimark:
	go run command/benchmark/main.go

# benchmark の測定結果を tsv で出力
export_log_tsv:
	go run command/log/main.go > log.tsv