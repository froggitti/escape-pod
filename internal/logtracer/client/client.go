package main

// Deprecated
// func main() {
// 	c, _ := client.New(
// 		client.WithTarget(":8084"),
// 		client.WithInsecureSkipVerify(),
// 	)

// 	_ = c.Connect()

// 	logcl := ep_logtracer.NewLogTracerClient(c.Conn())

// 	ctx, cancel := context.WithCancel(context.Background())

// 	done := make(chan bool)

// 	stream, _ := logcl.Trace(ctx, &ep_logtracer.TraceReq{})

// 	go func() {
// 		for {
// 			_, err := stream.Recv()
// 			if err == io.EOF {
// 				done <- true
// 				return
// 			}
// 			if err != nil {
// 				log.Fatalf("can not receive %v", err)
// 			}
// 		}
// 	}()
// 	<-done

// 	cancel()
// 	_ = c.Close()
// }
