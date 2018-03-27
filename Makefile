VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git rev-parse --short HEAD)

ifndef VERSION
	VERSION := latest
endif
	
LDFLAGS := -X main.commit=$(COMMIT) -X main.branch=$(BRANCH) -X main.version=$(VERSION)

prepare:
	go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
	go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli

migrate:
	migrate -database=$(DSN) -path=etc/migrations/ up

migrate-force:
	migrate -database=$(DSN) -force 2

httpd:
	go install -ldflags "-w -s $(LDFLAGS)" ./cmd/httpd

docker-image:
	docker build -t dumbwall:$(VERSION) .
