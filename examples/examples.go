package examples

// TODO convert to doc at some point
// config := tlswrap.NewConfig("./", []string{"domain.tld"})
// isDevMode := true
// exampleGinServer(config, isDevMode)
// exampleNativeServer1(config, isDevMode)
// exampleNativeServer2(config, isDevMode)
/*
func exampleNativeServer1(config tlswrap.Config, isDevMode bool) {
	// routes
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "lol")
	})

	// setup prod or dev server depending on isDevMode
	if isDevMode {
		err := tlswrap.StartDevServer("127.0.0.1:8443", config)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		config := tlswrap.NewStageConfigFromConfig(config)
		err := tlswrap.StartServer("0.0.0.0:8443", config)
		if err != nil {
			log.Fatalf("failed start server: %s", err)
		}
	}
}

func exampleNativeServer2(config tlswrap.Config, isDevMode bool) {
	// routes
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "lol")
	})

	// setup stage or dev server depending on isDevMode
	var server http.Server
	var httpHandler func()
	var err error
	if isDevMode {
		server, err = tlswrap.NewDevServer("127.0.0.1:8443", nil, config)
		if err != nil {
			log.Fatalf("failed to create dev server: %s", err)
		}
	} else {
		config := tlswrap.NewStageConfigFromConfig(config)
		server, httpHandler = tlswrap.NewServer("0.0.0.0:8443", nil, config)
		go httpHandler()
	}

	// serv it up
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal("failed to listen", err)
	}
}

func exampleNativeServer(config tlswrap.Config, isDevMode bool) {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "Hello Log")
	})
	if isDevMode {
		err := tlswrap.StartDevServer("127.0.0.1:8443", config)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// remove this line to use prod config
		config := tlswrap.NewStageConfigFromConfig(config)
		err := tlswrap.StartServer("0.0.0.0:8443", config)
		if err != nil {
			log.Fatalf("failed start server: %s", err)
		}
	}
}

func exampleGinServer(config tlswrap.Config, isDevMode bool) {
	handler := gin.Default()
	handler.GET("/", func(c *gin.Context) {
		fmt.Println("Hello World")
	})

	if isDevMode {
		address := "0.0.0.0:8443"
		if err := tlswrap.StartDevServerWithHandler(address, config, handler); err != nil {
			log.Fatal(err)
		}
	} else {
		address := "127.0.0.1:8443"
		if err := tlswrap.StartServerWithHandler(address, config, handler); err != nil {
			log.Fatal(err)
		}
	}
}
*/
