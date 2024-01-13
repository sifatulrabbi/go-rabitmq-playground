.PHONE: rabbitmq
rabbitmq:
	CONF_ENV_FILE="/opt/homebrew/etc/rabbitmq/rabbitmq-env.conf" /opt/homebrew/opt/rabbitmq/sbin/rabbitmq-server
dev:
	go run ./main.go
build:
	go build -o ./build/app ./main.go
start:
	./build/app
