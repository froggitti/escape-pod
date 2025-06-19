package saywhatnow

import (
	"log"

	vp "github.com/DDLbots/internal-api/go/sttpb"

	"github.com/DDLbots/saywhatnow/pkg/server/memorymatcher"
	"github.com/DDLbots/saywhatnow/pkg/server/stt"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

// Deprecated
// Start starts the saywhatnow service
func Start(db *mongo.Database, transport *grpc.Server) {
	// initialize the matcher
	m, err := memorymatcher.New(
		memorymatcher.WithDatabase(db),
	)
	if err != nil {
		log.Fatal(err)
	}

	// initialize the speech to text engine
	st, err := stt.New(
		stt.WithViper(),
		stt.WithIntentMatcher(m),
	)
	if err != nil {
		log.Fatal("stt initialization error: ", err)
	}

	vp.RegisterSTTServer(transport, st)
}
