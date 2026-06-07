all: bin/kolendar

SQL_FILES := $(wildcard *.sql)
GO_FILES := $(wildcard *.go) $(wildcard */*.go)
TEMPL_FILES := $(wildcard templates/*.templ)

db: $(SQL_FILES)
	sqlc generate

templ: $(TEMPL_FILES)
	templ generate

bin/kolendar: db templ $(GO_FILES)
	go build -o bin/kolendar .

clean:
	rm -rv db/ bin/

.PHONY: all clean templ
