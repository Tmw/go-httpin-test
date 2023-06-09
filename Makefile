.PHONY: run
run:
	go run main.go

.PHONY: call
call:
	curl -XPOST \
	-H 'content-type: application/json' \
	-d '{"name": "bobby", "age": 42, "extra_field": true}' \
	localhost:8080/wow
