# Auto TLS Server
A auto TLS wrap around your server.

Auto provision TLS for development, stage or production

# How to use

### Development Server
To start a development server that creates a self-signed certification and starts a server with it

```golang
import "github.com/ronniskansing/tlswrap"

func main() {
  http.HandleFunc("/", func(rw http.  ResponseWriter, r *http.Request) {
      fmt.Fprintln(rw, "Hello World")
  })
}
  config := tlswrap.NewConfig("./", []string{""})
  tlswrap.StartDevServer("127.0.0.1:8443", config)
}
  
```

### Stage Server
To start a stage server that uses a ACME stage server for certificates.

```golang
import "github.com/ronniskansing/tlswrap"

func main() {
  http.HandleFunc("/", func(rw http.  ResponseWriter, r *http.Request) {
      fmt.Fprintln(rw, "Hello World")
  })

  config := tlswrap.NewConfig(tlswrap.NewConfig("./", []string{""}))
  tlswrap.StartStageServer(":8443", config)
}
```

### Production Server
To start a production server that self-provisions certificates
```golang
import "github.com/ronniskansing/tlswrap"

func main() {
  http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "Hello World")
  })

  config := tlswrap.NewConfig(tlswrap.NewConfig("./", []string{"domain.tld"}))
  tlswrap.StartServer(":443", config)
}
```