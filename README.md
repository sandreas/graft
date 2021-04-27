```bash
GH_USER="sandreas"
GH_NAME="graft"
GH_URL="git@github.com:$GH_USER/$GH_NAME.git"

git clone "$GH_URL"
cd "$GH_NAME"

go mod init
GO111MODULE=on go get github.com/urfave/cli/v2
```
