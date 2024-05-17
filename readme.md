compile
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o test -trimpath -ldflags "-s -w -buildid=" .

```
