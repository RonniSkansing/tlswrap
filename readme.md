# Auto TLS Server
A auto TLS wrap around your server.

Auto provision TLS for development, stage or production

# How to use

### Development Server
To start a development server that creates a self-signed certification and starts a server with it

```golang
config := tlswrap.NewConfig("./", []string{""})
tlswrap.StartServerWithHandler("127.0.0.1:8443", config, nil)

http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "Hello World")
})
```

### Stage Server
To start a stage server that uses a ACME stage server for certificates.

```golang
config := tlswrap.NewStageConfigFromConfig(tlswrap.NewConfig("./", []string{""}))
tlswrap.StartServerWithHandler("127.0.0.1:8443", config, nil)

http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "Hello World")
})
```

### Production Server
To start a production server that self-provisions certificates
```golang
config := tlswrap.NewStageConfigFromConfig(tlswrap.NewConfig("./", []string{"domain.tld"}))
tlswrap.StartServerWithHandler(":443", config, nil)

http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(rw, "Hello World")
})
```