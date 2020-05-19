lumix webcam
===

## What's this

convert from Lumix UDP streaming to HTTP MJPEG streaming.

## How to build

```
go build
```

## How to use

1. connect Lumix to wifi.
2. get Lumix's IP.
3. start streaming

```
./keep.sh Lumix's IP
```

4. start server

```
./lumix-webcam -addr :8080
```

5. open browser

```
open http://localhost:8080
```


## Tested Devices

- Lumix DMC-GX7mk2
