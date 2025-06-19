package saywhatnowmanager

// Deprecated
// Start starts the saywhatnow management service
// func Start(db *mongo.Database, transport *grpcserver.Server) {
// 	if err := transport.RegisterHTTPService(
// 		[]func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{
// 			sttpb.RegisterSTTManagerHandlerFromEndpoint,
// 		},
// 	); err != nil {
// 		log.Fatal(err)
// 	}

// 	mgr, err := management.New(
// 		management.WithMongoBackend(db),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	sttpb.RegisterSTTManagerServer(transport.Transport(), mgr)
// }
