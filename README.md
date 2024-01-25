# treblle-go-gin
The offical Treblle SDK for Go using Gin


## Installation

```shell
go get github.com/treblle/treblle-go-gin
```

Trebble uses [Go Modules](https://github.com/golang/go/wiki/Modules) to manage dependencies.


## Basic configuration

Configure Treblle at the start of your `main()` function:

```go
import "github.com/treblle/treblle-go-gin"

func main() {
	treblle.Configure(treblle.Configuration{
		APIKey:     "YOUR API KEY HERE",
		ProjectID:  "YOUR PROJECT ID HERE",
		KeysToMask: []string{"password", "card_number"}, // optional, mask fields you don't want sent to Treblle
		ServerURL:  "https://rocknrolla.treblle.com",    // optional, don't use default server URL
	}

    // rest of your program.
}

```


After that, just use the middleware with any of your Gin handlers:
 ```go
r := gin.Default()
r.Use(treblle.GinMiddleware())
```


