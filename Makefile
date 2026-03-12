# Vars
APP_NAME=credit-risk-mvp
MAIN_PATH=cmd/main.go

# default make
all: build

# Run application
run:
	@echo "🚀 Starting app.."
	go run $(MAIN_PATH)

# Build binary
build:
	@echo "📦 Build binary.."
	go build -o $(APP_NAME) $(MAIN_PATH)

# TESTS
test:
	@echo "🧪 Starting tests.."
	go test -v ./...

clean:
	@echo "🧹 Cleaning..."
	rm -f $(APP_NAME)

# Migration up
migrate-up:
	@echo "🦆 Накатываем миграции..."
	goose -dir migrations postgres "user=app password=pass dbname=creditrisk sslmode=disable" up

# Migration down
migrate-down:
	@echo "🔙 Откат последней миграции..."
	goose -dir migrations postgres "user=app password=pass dbname=creditrisk sslmode=disable" down

.PHONY: all run build test clean migrate-up migrate-down