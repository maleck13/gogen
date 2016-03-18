
## Generate a basic golang REST api server

Ensure you have $GOPATH set

```

go get github.com/maleck13/gogen

go install .

gogen generate --package=github.com/example/app

cd $GOPATH/src/github.com/example/app

go get .

go build .

./app serve

curl http://localhost:3000

```
