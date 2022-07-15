.ONESHELL:

define HELP
Usage:
	make [command]

Commands:
	help	ヘルプを表示します.
	init	開発に使用する周辺ツールをセットアップします
endef
export HELP

define GO_TOOLS
github.com/volatiletech/sqlboiler/v4@v4.11.0
github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-sqlite3@v4.11.0
endef
export GO_TOOLS

.PHONY: help
help:
	@echo "$$HELP"

.PHONY: init
init:
	@for v in $$GO_TOOLS; do go install $$v; done

.PHONY: boil
boil:
	sqlboiler sqlite3
	go mod tidy