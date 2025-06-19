package chipper

import (
	"fmt"
	"os"

	pb "github.com/DDLbots/api/go/chipperpb"
	"github.com/DDLbots/ddl-chipper/pkg/intentextender"
	"github.com/DDLbots/ddl-chipper/pkg/intentextender/escapepod"
	"github.com/DDLbots/ddl-chipper/pkg/tts/saywhatnow"
	"github.com/DDLbots/internal-api/go/sttpb"

	"log"

	"github.com/DDLbots/chipper/pkg/server"
	"google.golang.org/grpc"
)

// Start starts the chipper service
func Start(transport *grpc.Server, sttClient sttpb.STTClient) {
	// initialize the saywhatnow processor
	p, err := saywhatnow.New(
		saywhatnow.WithSTTClient(sttClient),
		saywhatnow.WithTarget(fmt.Sprintf("127.0.0.1:%s", os.Getenv("DDL_RPC_PORT"))),
		saywhatnow.WithViper(),
		saywhatnow.WithFunctionServer(getFunctionServer()),
	)
	if err != nil {
		log.Fatalf("new saywhatnow client %v", err)
	}

	s, _ := server.New(
		// server.WithLogger(log.Base()),
		server.WithIntentProcessor(p),
		server.WithKnowledgeGraphProcessor(p),
		server.WithIntentGraphProcessor(p),
	)

	pb.RegisterChipperGrpcServer(transport, s)
}

func getFunctionServer() intentextender.Server {
	e := os.Getenv("ENABLE_EXTENSIONS")
	if e == "" || e == "false" || e == "FALSE" {
		return nil
	}

	esc, err := escapepod.New()
	if err != nil {
		log.Fatal(err)
	}
	return esc
}
