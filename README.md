# docker-ipam-proxy-plugin
Docker IPAM Plugin (v2) proxy to remote HTTP IPAM service

## Configuration and installation and from source
```
$ git clone https://github.com/eyz/docker-ipam-proxy-plugin.git
$ cd docker-ipam-proxy-plugin
$ make build
```
NOTE: in plugin/config.json, change the IPAM_HTTP_PROXY_HOST default from localhost:8080 to your desired target before creating and enabling this plugin; you can only pass KEY=VALUE pairs to Docker plugins when installing (not when enabling), so when you build from source this is a typical step
```
$ docker plugin create eyz/docker-ipam-proxy-plugin:latest .
$ docker plugin enable eyz/docker-ipam-proxy-plugin:latest
```
## Installation from Docker Hub
```
# where IPAMSERVER:PORT is your remote IPAM HTTP server and port
$ docker plugin install eyz77/docker-ipam-proxy-plugin IPAM_HTTP_PROXY_HOST=IPAMSERVER:PORT 

# where IPAMSERVER is your remote IPAM HTTP server (defaulting to TCP 80 for the port)
$ docker plugin install eyz77/docker-ipam-proxy-plugin IPAM_HTTP_PROXY_HOST=IPAMSERVER 
```
NOTE: in my testing, the plugin requires use of the host network, so you will need to grant the "network: [host]" permission

## Example usage
### Install, pointing to IPAM HTTP server IP and port "192.168.0.100:8080"
```
$ docker plugin install eyz77/docker-ipam-proxy-plugin:latest IPAM_HTTP_PROXY_HOST=192.168.0.100:8080
Plugin "eyz77/docker-ipam-proxy-plugin:latest" is requesting the following privileges:
 - network: [host]
Do you grant the above permissions? [y/N] y
latest: Pulling from eyz77/docker-ipam-proxy-plugin
21d1727448ec: Download complete 
Digest: sha256:4858d0f4953e2b285d6ff8ae9741c9da0cbbbd72b917efad84405f0c8de4e027
Status: Downloaded newer image for eyz77/docker-ipam-proxy-plugin:latest
Installed plugin eyz77/docker-ipam-proxy-plugin:latest
```
### Creation of test network "ipamTestNetwork" using the installed and configured IPAM driver
```
$ docker network create --ipam-driver=eyz77/docker-ipam-proxy-plugin:latest ipamTestNetwork
cbe2edfa59b78d5361213f9a9d7333dda73749dd4b24e76c4397f9209a5b5916
```
### Verification that a subnet and gateway was allocated from the remote IPAM service (at 192.168.0.100:8080)
```
$ docker network inspect cbe
[
    {
        "Name": "ipamTestNetwork",
        "Id": "cbe2edfa59b78d5361213f9a9d7333dda73749dd4b24e76c4397f9209a5b5916",
        "Created": "2017-05-30T11:40:55.19686084-07:00",
        "Scope": "local",
        "Driver": "bridge",
        "EnableIPv6": false,
        "IPAM": {
            "Driver": "eyz77/docker-ipam-proxy-plugin:latest",
            "Options": {},
            "Config": [
                {
                    "Subnet": "10.100.100.0/24",
                    "Gateway": "10.100.100.1"
                }
            ]
        },
        "Internal": false,
        "Attachable": false,
        "Ingress": false,
        "Containers": {},
        "Options": {},
        "Labels": {}
    }
]
```
### Launch of busybox container in the test network, and verification that an IP was allocated from the remote IPAM service
```
$ docker run -ti --network=ipamTestNetwork busybox 
/ # ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
78: eth0@if79: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue 
    link/ether 02:42:0a:64:64:02 brd ff:ff:ff:ff:ff:ff
    inet 10.100.100.2/24 scope global eth0
       valid_lft forever preferred_lft forever
```
