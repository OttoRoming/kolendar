all: bin/kolendar

SQL_FILES := $(wildcard *.sql)
GO_FILES := $(wildcard *.go) $(wildcard */*.go)

db: $(SQL_FILES)
	sqlc generate

bin/kolendar: db $(GO_FILES)
	go build -o bin/kolendar .

clean:
	rm -rv db/ bin/

.PHONY: all clean

