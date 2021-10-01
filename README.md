# GoTSW

The official Teraswitch Go client.  

## Install
```sh
go get -u github.com/teraswitch/gotsw
```

## Usage
```go
import "github.com/teraswitch/gotsw"
```

### Authentication
```go
package main

import (
    "github.com/teraswitch/gotsw"
)

func main() {
    client := gotsw.NewFromToken("id:secret")
}
```