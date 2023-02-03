[![Go Reference](https://pkg.go.dev/badge/github.com/Ompluscator/go-promise.svg)](https://pkg.go.dev/github.com/Ompluscator/go-promise)
![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/ompluscator/go-promise)](https://goreportcard.com/report/github.com/ompluscator/go-promise)

# JS Promise in Golang

This library introduces the [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) feature from JavaScript in Go. 
It relies on generics and requires minimal version of Go to be 1.18.

## Features

The library supports following features.

### Wide Promise creation

By using standard creation method with _resolve_ and _reject_ functions
as arguments:

```go
import (
    "github.com/ompluscator/go-promise"
)


func main() {
    promise := go_promise.New(func(resolve go_promise.ResolveFunc[int], reject go_promise.RejectFunc) {
        resolve(10)
    })

    value, err := go_promise.Await[int](promise)
    fmt.Println(value, err)
    // Output: 10, nil
}
```

By using simple function:

```go
import (
    "github.com/ompluscator/go-promise"
)


func main() {
    promise := go_promise.Function(func() (int, error) {
        return 10, nil
    })

    value, err := go_promise.Await[int](promise)
    fmt.Println(value, err)
    // Output: 10, nil
}
```

By using _resolve_ function directly:

```go
import (
    "github.com/ompluscator/go-promise"
)


func main() {
    promise := go_promise.Resolve(10)

    value, err := go_promise.Await[int](promise)
    fmt.Println(value, err)
    // Output: 10, nil
}
```

By using _reject_ function directly:

```go
import (
    "github.com/ompluscator/go-promise"
)


func main() {
    promise := go_promise.Resolve(10)

    value, err := go_promise.Await[int](promise)
    fmt.Println(value, err)
    // Output: 10, nil
}
```

### Chaining

Chaining with _then_ and _catch_ methods is also supported:

```go
import (
    "github.com/ompluscator/go-promise"
)


func main() {
    promise := go_promise.Function(func() (int, error) {
        return 10, nil
    }).With(go_promise.Then(func(value int) (float64, error) {
        return 0, errors.New("error")
    })).With(go_promise.Catch(func(err error) float64 {
        return 11.1
    })).With(go_promise.Then(func(value float64) (bool, error) {
        return value > 11, nil
    }))

    value, err := go_promise.Await[bool](promise)
    fmt.Println(value, err)
    // Output: true, nil
}
```

### Wrappers

We can wrap a promise with timeout:

```go
import (
	"time"
	
    "github.com/ompluscator/go-promise"
)


func main() {
    promise := go_promise.Timeout[int](go_promise.Function(func() (int, error) {
        time.Sleep(time.Second)
        return 10, nil
    }), 500 * time.Milisecond)

    value, err := go_promise.Await[int](promise)
    fmt.Println(value, err)
    // Output: 0, promise.timeout
}
```

We can wrap a promise with retrial policy:

```go
import (
	"github.com/ompluscator/go-promise"
)


func main() {
	try := 0
	
    promise := go_promise.WithRetry[int](go_promise.Function(func() (int, error) {
        if try < 3 {
            try++
            return 0, errors.New("error")
        }
        return 10, nil
    }), 3)

    value, err := go_promise.Await[int](promise)
    fmt.Println(value, err)
    // Output: 10, nil
}
```

We can wrap a promise with pre-execution policy:

```go
import (
	"time"
	
	"github.com/ompluscator/go-promise"
)


func main() {
	promise := go_promise.AsPreExecuted[int](go_promise.Function(func() (int, error) {
        time.Sleep(time.Second)
        return 10, nil
    }))

    time.Sleep(time.Second)
	// promise is already executed, so there is no additional waiting time

    value, err := go_promise.Await[int](promise)
    fmt.Println(value, err)
    // Output: 10, nil
}
```

### Resolvers

Waiting for results of all promises with method _all_:

```go
import (
    "github.com/ompluscator/go-promise"
)


func main() {
    promise := go_promise.All[int](go_promise.Promises{
        go_promise.Function(func() (int, error) {
            return 10, nil
        }),
        go_promise.Function(func() (int, error) {
            return 11, nil
        }),
        go_promise.Function(func() (int, error) {
            return 12, nil
        }),
    })

    value, err := go_promise.Await[[]int](promise)
    fmt.Println(value, err)
    // Output: [11, 10, 12], nil
}
```

Waiting for results of all promises with method _allSettled_:

```go
import (
	"errors"
	
    "github.com/ompluscator/go-promise"
)


func main() {
    promise := go_promise.AllSettled[int](go_promise.Promises{
        go_promise.Function(func() (int, error) {
            return 10, nil
        }),
        go_promise.Function(func() (int, error) {
            return 0, errors.New("error")
        }),
        go_promise.Function(func() (int, error) {
            return 12, nil
        }),
    })

    settled, err := go_promise.Await[go_promise.SettledResults[int]](promise)
    fmt.Println(settled.Values(), settled.Errors())
    // Output: [10, 12], [error]
}
```

Waiting for results of all promises with method _any_:

```go
import (
    "errors"

    "github.com/ompluscator/go-promise"
)


func main() {
    promise := go_promise.Any[int](go_promise.Promises{
        go_promise.Function(func() (int, error) {
            return 0, errors.New("error")
        }),
        go_promise.Function(func() (int, error) {
            return 11, nil
        }),
        go_promise.Function(func() (int, error) {
            return 0, errors.New("error")
        }),
    })

    value, err := go_promise.Await[int](promise)
    fmt.Println(value, err)
    // Output: 11, nil
}
```

Waiting for results of all promises with method _race_:

```go
import (
    "errors"

    "github.com/ompluscator/go-promise"
)


func main() {
    promise := go_promise.Race[int](go_promise.Promises{
        go_promise.Function(func() (int, error) {
            return 10, nil
        }),
        go_promise.Function(func() (int, error) {
            return 11, nil
        }),
        go_promise.Function(func() (int, error) {
            return 12, nil
        }),
    })

    value, err := go_promise.Await[int](promise)
    fmt.Println(value, err)
    // Output: 11, nil
}
```