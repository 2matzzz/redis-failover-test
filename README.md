# redis-failover-test

A small CLI to check if the Redis server failover is functioning.

## Usage

```
$ go run main.go --help
Usage of /var/folders/n3/6zd4mt6n70ldwbhy03s3jt940000gn/T/go-build2972800864/b001/exe/main:
  -h string
        Redis server IP address (default "127.0.0.1")
  -p string
        Redis server port (default "6379")
```

## Test

Assuming you have [Multipass](https://multipass.run/) set up to test the program locally.

1. Launch Redis server with a new VM:
   ```console
   $ multipass launch -n redis-server --disk 5G --memory 2G --cpus 5 --cloud-init cloud-init.yml
   ```
2. Run the program:
   ```console
   $ go run main.go -h $(multipass info redis-server --format json | jq -r '.info."redis-server".ipv4[0]')
   2024/10/29 16:26:35 Connecting to Redis at 192.168.205.27:6379
   2024/10/29 16:26:35 maintainConnectionPing: PONG
   2024/10/29 16:26:35 newConnectionPing: PONG
   2024/10/29 16:26:36 maintainConnectionPing: PONG
   2024/10/29 16:26:36 newConnectionPing: PONG
   2024/10/29 16:26:37 maintainConnectionPing: PONG
   2024/10/29 16:26:37 newConnectionPing: PONG
   ...
   ```

## Clean Up

```
$ multipass delete redis-server
$ multipass purge
```
