package token

// Start starts the token service
// func Start(signingKey string, transport *grpc.Server) {
// 	tok, err := noop.New(
// 		noop.WithSigningKey(signingKey),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// create the jdocs client
// 	c, err := client.New(
// 		// Bit of a hack here to guarantee the port is the same
// 		client.WithTarget(fmt.Sprintf("127.0.0.1:%s", os.Getenv("DDL_RPC_PORT"))),
// 		client.WithInsecureSkipVerify(),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if err := c.Connect(); err != nil {
// 		log.Fatal(err)
// 	}

// 	s, err := server.New(
// 		// server.WithLogger(log.Base()),
// 		server.WithTokenizer(tok),
// 		server.WithJdocsClient(
// 			jdocspb.NewJdocsClient(c.Conn()),
// 		),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	pb.RegisterTokenServer(transport, s)
// }
