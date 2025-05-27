export CGO_ENABLED=0

# ==========================
# プロジェクト全般
# ==========================

.DEFAULT_GOAL := help
.PHONY: help
help: ## helpを表示
	@echo '  see: myblog'
	@echo ''
	@grep -E '^[%/0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: format
format: ## format実行
	docker compose exec -T app sh -c "go fmt ./..."

.PHONY: show-rdb-log
show-rdb-log: ## rdbコンテナで実行されたSQLのログを表示
	docker compose exec rdb tail -fn ${TAIL_LINES} /var/tmp/mysqld.log

# ==========================
# テスト
# ==========================

.PHONY: lint
lint: ## lint実行
	docker compose exec -T app sh -c "go run github.com/golangci/golangci-lint/cmd/golangci-lint run";

.PHONY: create-empty-test-db
create-empty-test-db: ## テスト用のDBを作成
	@docker compose exec -e MYSQL_PWD=root rdb-test sh -c "mysql -u root -e 'CREATE DATABASE IF NOT EXISTS test_empty'"

.PHONY: test
UPDATE_SNAPSHOTS ?= ''
VERBOSE ?= 0
UNIT_TEST_DIR = ./app/...

test: create-empty-test-db ## test実行 (特定のテストケースだけを実行したい場合は`CASE=TestFoo make test`のように実行する。)
	@docker compose exec -e APP_ENV=test -e UPDATE_SNAPSHOTS=${UPDATE_SNAPSHOTS} app sh -c '\
		GOTEST_OPTS="-short"; \
		if [ -n "${CASE}" ] || [ "${VERBOSE}" = "1" ]; then \
			GOTEST_OPTS="$${GOTEST_OPTS} -v"; \
		fi; \
		if [ -n "${CASE}" ]; then \
			TARGET_FILE_PATH=$$(grep -rl "func ${CASE}(" $$(find . -name "*_test.go")); \
			TARGET_DIR=$$(dirname "$$TARGET_FILE_PATH"); \
			go test $${GOTEST_OPTS} -run=${CASE} $$TARGET_DIR; \
		else \
			go test $${GOTEST_OPTS} ${UNIT_TEST_DIR}; \
		fi'

# ==========================
# データベース
# ==========================

.PHONY: migrate
migrate: ## マイグレーション実行($SQL_FILEをdb/migrations以下にあるsqlファイル名に指定)
	docker compose exec rdb sh -c "mysql -u root -proot myblog < ./docker-entrypoint-initdb.d/${SQL_FILE}"

.PHONY: migrate-new
migrate-new: ## マイグレーションファイルを新規作成
	touch "db/migrations/$$(date +"%Y%m%d%H%M%S")-table_name.sql"

.PHONY: migrate-reset
migrate-reset: ## DBをリセットしてからマイグレーション再実行（並列実行）
	@echo "並列でDBリセットとマイグレーションを実行中..."
	@for CONTAINER in rdb rdb-test; do \
		( \
			DB_NAME=$$(if [ "$$CONTAINER" = "rdb" ]; then echo "myblog"; else echo "test"; fi); \
			docker compose exec -T -e MYSQL_PWD=root $$CONTAINER sh -c "mysql -u root -e 'drop database $$DB_NAME'"; \
			docker compose exec -T -e MYSQL_PWD=root $$CONTAINER sh -c "mysql -u root -e 'create database $$DB_NAME'"; \
			docker compose exec -T -e MYSQL_PWD=root $$CONTAINER sh -c "for SQL_FILE in \$$(ls /docker-entrypoint-initdb.d/*.sql | xargs -n1 basename); do \
				echo \"migrate: container:$$CONTAINER, file:\$$SQL_FILE\"; \
				mysql -u root $$DB_NAME < /docker-entrypoint-initdb.d/\$$SQL_FILE; \
			done"; \
		) & \
	done; \
	wait
	@echo "すべてのDBリセットとマイグレーションが完了しました"

.PHONY: import-seed-data
import-seed-data: ## db/seed以下にあるsqlファイルを実行
	for SQL_FILE in `ls db/seed/*.sql | sed 's|db/seed/||'`; do \
		echo "import: $$SQL_FILE"; \
		cat db/seed/$$SQL_FILE | docker compose exec -T -e MYSQL_PWD=root rdb mysql -u root myblog; \
	done
