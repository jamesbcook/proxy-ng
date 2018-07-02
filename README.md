# Proxy-ng

## Purpose

Wanted to make use of multiple socks proxies in a random order. Along with using socks proxies in a random order, I wanted each request to have a random useragent.

## Setup

### Build

* `git clone https://github.com/jamesbcook/proxy-ng.git`
* `make`
* `./proxy-ng -help`

### Pre-Built binary

* Download from [here](https://github.com/jamesbcook/proxy-ng/releases)
* `./proxy-ng -help`

## Running

Proxy-ng opens the following ports:

* 9292 for the local socks proxy
* 9293 for the local http proxy

Unless specified by a flag, Proxy-ng will use look for the following files in its running directory:

* useragents.json
* socks5-proxies.json

## Example

I use [cloud-proxy](https://github.com/tomsteele/cloud-proxy) to setup multiple socks proxies

```bash
./cloud-proxy -token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -key 'xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx' -count 5
2018/07/01 18:12:22 Creating 1 droplets to region nyc1
2018/07/01 18:12:23 Creating 1 droplets to region sgp1
2018/07/01 18:12:24 Creating 1 droplets to region lon1
2018/07/01 18:12:25 Creating 1 droplets to region nyc3
2018/07/01 18:12:26 Creating 1 droplets to region ams3
2018/07/01 18:12:27 Droplets deployed. Waiting 100 seconds...
2018/07/01 18:14:08 SSH proxy started on port 55555 on droplet name: cloud-proxy-OucHZHiV IP: 159.65.236.62
2018/07/01 18:14:09 SSH proxy started on port 55556 on droplet name: cloud-proxy-q4AZAYVN IP: 178.128.95.250
2018/07/01 18:14:10 SSH proxy started on port 55557 on droplet name: cloud-proxy-RQLq0UQm IP: 139.59.173.24
2018/07/01 18:14:11 SSH proxy started on port 55558 on droplet name: cloud-proxy-eVd59B6d IP: 209.97.153.98
2018/07/01 18:14:11 SSH proxy started on port 55559 on droplet name: cloud-proxy-oSRSfFO1 IP: 188.166.6.62
2018/07/01 18:14:11 proxychains config
socks5 127.0.0.1 55555
socks5 127.0.0.1 55556
socks5 127.0.0.1 55557
socks5 127.0.0.1 55558
socks5 127.0.0.1 55559
2018/07/01 18:14:11 socksd config
"upstreams": [
{"type": "socks5", "address": "127.0.0.1:55555"},
{"type": "socks5", "address": "127.0.0.1:55556"},
{"type": "socks5", "address": "127.0.0.1:55557"},
{"type": "socks5", "address": "127.0.0.1:55558"},
{"type": "socks5", "address": "127.0.0.1:55559"}
]
2018/07/01 18:14:11 Please CTRL-C to destroy droplets
```

After pointing the browser to the local http listener setup by Proxy-ng and heading to whatsmyip.org it showed my hostname and IP to be different.

![whatsmyip](https://github.com/jamesbcook/proxy-ng/raw/master/media/whatsmyip-results.png)

## Help Output

```text
Usage of ./proxy-ng:
  -http string
        HTTP listener to accept connections, this changes the useragent on each request (default "localhost:9293")
  -socks string
        Local socks listener to accept connections (default "localhost:9292")
  -socksFile string
        Socks file that contains socks proxies to use (default "socks5-proxies.json")
  -uaFile string
        Json file that contains useragents to use (default "useragents.json")
  -verbose
        Verbose output from proxy
  -version
        Current Version
```