### Build apps
```shell
make build
```

### Run apps
```shell
make run
```

### Build and run apps
```shell
make all
```

### Testing (open 3 windows in terminal)
#### Window 1
```shell
tail -f logs/downstream.log
```
#### Window 2
```shell
tail -f logs/downstream.log
```
#### Window 3
```shell
curl 127.0.0.1:9010
```
