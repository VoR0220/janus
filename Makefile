.PHONY: install
install: 
	go install github.com/qtumproject/janus/cli/janus

.PHONY: release
release: darwin linux

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./build/janus-darwin-amd64 github.com/qtumproject/janus/cli/janus

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build -o ./build/janus-linux-amd64 github.com/qtumproject/janus/cli/janus

.PHONY: quick-start
quick-start:
	cd docker && ./spin_up.sh && cd ..

.PHONY: docker-dev
docker-dev:
	docker build --no-cache -f ./docker/standalone/Dockerfile -t qtum/janus:dev .
	
.PHONY: local-dev
local-dev:
	docker run --name qtum_testchain -d -p 3889:3889 qtum:qtum qtumd -regtest -rpcbind=0.0.0.0:3889 -rpcallowip=0.0.0.0/0 -logevents -rpcuser=qtum -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole
	sleep 3
	docker cp ${GOPATH}/src/github.com/qtumproject/janus/docker/fill_user_account.sh qtum_testchain:.
	docker exec qtum_testchain /bin/sh -c ./fill_user_account.sh
	QTUM_RPC=http://qtum:testpasswd@localhost:3889 QTUM_NETWORK=regtest janus --accounts ./docker/standalone/myaccounts.txt --dev

.PHONY: local-dev-logs
local-dev-logs:
	docker run --name qtum_testchain -d -p 3889:3889 qtum:qtum qtumd -regtest -rpcbind=0.0.0.0:3889 -rpcallowip=0.0.0.0/0 -logevents -rpcuser=qtum -rpcpassword=testpasswd -deprecatedrpc=accounts -printtoconsole
	sleep 3
	docker cp ${GOPATH}/src/github.com/qtumproject/janus/docker/fill_user_account.sh qtum_testchain:.
	docker exec qtum_testchain /bin/sh -c ./fill_user_account.sh
	QTUM_RPC=http://qtum:testpasswd@localhost:3889 QTUM_NETWORK=regtest janus --accounts ./docker/standalone/myaccounts.txt --dev > janus_dev_logs.txt