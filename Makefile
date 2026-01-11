.PHONY: batch-migrate batch-seed batch-rankings batch-all db-start db-stop db-status

# Batch jobs (run on Cloud Run)
batch-migrate:
	./scripts/batch.sh migrate

batch-seed:
	./scripts/batch.sh seed

batch-rankings:
	./scripts/batch.sh rankings

batch-all:
	./scripts/batch.sh all

# Database management
db-start:
	curl -s "https://db-manager-671942133800.asia-northeast1.run.app?action=start"

db-stop:
	curl -s "https://db-manager-671942133800.asia-northeast1.run.app?action=stop"

db-status:
	curl -s "https://db-manager-671942133800.asia-northeast1.run.app?action=status"

# Deploy
deploy-backend:
	gh workflow run "Backend CI" --ref main || echo "Triggering via push..."

# Import CSV data
import-csv:
	@if [ -z "$(FILE)" ] || [ -z "$(TOKEN)" ]; then \
		echo "Usage: make import-csv FILE=/path/to/data.csv TOKEN=your-auth-token"; \
		echo "  Get auth token from browser DevTools after logging in"; \
		exit 1; \
	fi
	cd scripts && go run import-csv.go "$(FILE)" "$(TOKEN)"

# Help
help:
	@echo "Available commands:"
	@echo "  make batch-migrate   - Run database migrations"
	@echo "  make batch-seed      - Create initial admin user"
	@echo "  make batch-rankings  - Update keyword rankings"
	@echo "  make batch-all       - Run all batch jobs"
	@echo "  make db-start        - Start Cloud SQL"
	@echo "  make db-stop         - Stop Cloud SQL"
	@echo "  make db-status       - Check Cloud SQL status"
	@echo "  make import-csv FILE=path TOKEN=token - Import CSV data"
