package jdocs

import (
	pb "github.com/DDLbots/api/go/jdocspb"
	"github.com/DDLbots/jdocs/pkg/model"
	"google.golang.org/grpc"

	"log"

	"github.com/DDLbots/jdocs/pkg/jdocs"
)

// Start starts the jdocs service
func Start(db model.DB, transport *grpc.Server) {

	s, err := jdocs.New(
		// jdocs.WithLogger(log.Base()),
		jdocs.WithDB(db),
	)
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterJdocsServer(transport, s)
}
