rabbitmq:
	CONF_ENV_FILE="/opt/homebrew/etc/rabbitmq/rabbitmq-env.conf" /opt/homebrew/opt/rabbitmq/sbin/rabbitmq-server
dev:
	go run ./*.go
build:
	go build -o ./build/app ./*.go
start:
	./build/app

.PHONE: rabbitmq, dev, build, start
