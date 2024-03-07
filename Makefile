.PHONY: run

run:
	go run main.go

export_log_tsv:
	go run command/log/main.go > log.tsv