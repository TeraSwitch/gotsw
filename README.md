# GoTSW

The official Teraswitch Go client.

## Install

```sh
go get github.com/teraswitch/gotsw/v2
```

## Usage

```go
import "github.com/teraswitch/gotsw/v2"
```

### Authentication

```go
package main

import (
    "github.com/teraswitch/gotsw/v2"
)

func main() {
    client := gotsw.New("id:secret")
}
```

For more examples, see the [examples](examples) directory.
