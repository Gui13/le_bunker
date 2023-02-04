.PHONY: all pech client clean stop

all: client/client pech


pech:
	GOOS=linux go build

client/client:
	cd client && GOOS=linux go build client.go

run: all
	docker-compose down && docker-compose build && docker-compose up -d && docker-compose logs -f

stop:
	docker-compose down

clean:
	rm -f pech client/client
