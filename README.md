# Project setup

## Basics
```bash
GH_USER="sandreas"
GH_NAME="graft"
GH_URL="git@github.com:$GH_USER/$GH_NAME.git"
CHECKOUT_PATH="$GOPATH/pkg/mod/github.com/$GH_USER"
mkdir -p "$CHECKOUT_PATH"
cd "$CHECKOUT_PATH"
git clone "$GH_URL"
cd "$GH_NAME"

go mod init
```

## Commandline app boilerplate
```bash
# create directory
mkdir -p "cmd/$GH_NAME/"

# create main.go
cat <<EOF > "cmd/$GH_NAME/main.go"
package main

func main() {
  println("hello world")
}
EOF

# test the app
go run ./...
```





