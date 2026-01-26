.PHONY: new-migration

new-migration:
	@echo "This command requires goloang migrate; https://github.com/golang-migrate/migrate"
	@echo "Usage example, make new-migration name=create_urls_table"
	@migrate create -ext sql -dir migration -seq $(name)
