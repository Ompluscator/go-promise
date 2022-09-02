# JS Promise in Golang

```go
func main() {
	first := New(func(resolve ResolveFunc[int], reject RejectFunc) {
		resolve(10)
	})

	second := WithThen[int, float64](first, func(value int) (float64, error) {
		return 0, errors.New("error")
	})

	catched := WithCatch[float64, float64](second, func(err error) float64 {
		return 11.1
	})

	third := WithThen[float64, bool](catched, func(value float64) (bool, error) {
		return value > 11, nil
	})

	value, err := third.Await()
	if err != nil {
		fmt.Println("error is not expected")
	}
	if value != true {
        fmt.Println("value is not true")
	}
}
```