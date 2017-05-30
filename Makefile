PLUGIN_NAME=eyz77/docker-ipam-proxy-plugin

clean:
	rm -rf ./plugin ./bin
	rm -f docker-ipam-proxy-plugin
	docker plugin disable ${PLUGIN_NAME} || true
	docker plugin rm ${PLUGIN_NAME} || true
	docker rm -vf tmp || true
	docker rmi ipam-proxy-build-image || true
	docker rmi ${PLUGIN_NAME}:rootfs || true

build:
	docker build -t ipam-proxy-build-image -f Dockerfile.build .
	docker create --name tmp ipam-proxy-build-image
	docker cp tmp:/go/bin/docker-ipam-proxy-plugin .
	docker rm -vf tmp
	docker rmi ipam-proxy-build-image
	docker build -t ${PLUGIN_NAME}:rootfs .
	mkdir -p ./plugin/rootfs
	docker create --name tmp ${PLUGIN_NAME}:rootfs
	docker export tmp | tar -x -C ./plugin/rootfs
	cp config.json ./plugin/
	docker rm -vf tmp
	rm -f docker-ipam-proxy-plugin

create-plugin:
	docker plugin create ${PLUGIN_NAME} ./plugin

push-plugin:
	docker plugin push ${PLUGIN_NAME}

rm-plugin:
	docker plugin rm ${PLUGIN_NAME}

push:  clean build create-plugin push-plugin rm-plugin clean
