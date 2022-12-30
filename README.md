# JS Promise in Golang

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
    if err != nil {
        fmt.Println("error is not expected")
    }
    if value != true {
        fmt.Println("value is not true")
    }
}
```