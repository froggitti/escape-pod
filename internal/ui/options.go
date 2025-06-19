package ui

import (
	"crypto/tls"
	"fmt"

	"github.com/DDLbots/escape-pod/internal/ui/statuswebsocket"
	"github.com/digital-dream-labs/vector-bluetooth/ble"
)

type Option interface {
	applyOption(*Server) error
}

type applyOptionFunc func(*Server) error

func (f applyOptionFunc) applyOption(s *Server) error {
	return f(s)
}

func WithX509KeyPairString(pubCert, privCert string) Option {
	return applyOptionFunc(func(s *Server) error {
		crt, err := tls.X509KeyPair(
			[]byte(pubCert),
			[]byte(privCert),
		)
		if err != nil {
			return fmt.Errorf("load x509 key pair: %v", err)
		}
		s.crt = crt
		return nil
	})
}

func WithDLStatusChannel(ch chan ble.StatusChannel) Option {
	return applyOptionFunc(func(s *Server) error {
		s.dlCh = ch
		s.dlWs = statuswebsocket.New(ch)
		return nil
	})
}

func WithRootDir(rootDir string) Option {
	return applyOptionFunc(func(s *Server) error {
		s.rootDir = rootDir
		return nil
	})
}

func WithOTADir(otaDir string) Option {
	return applyOptionFunc(func(s *Server) error {
		s.otaDir = otaDir
		return nil
	})
}

func WithLogsDir(logsDir string) Option {
	return applyOptionFunc(func(s *Server) error {
		s.logsDir = logsDir
		return nil
	})
}

func WithUIDir(uiDir string) Option {
	return applyOptionFunc(func(s *Server) error {
		s.uiDir = uiDir
		return nil
	})
}

// Deprecated
func WithOTAPort(otaPort string) Option {
	return applyOptionFunc(func(s *Server) error {
		s.otaPort = otaPort
		return nil
	})
}

func WithPort(port string) Option {
	return applyOptionFunc(func(s *Server) error {
		s.port = port
		return nil
	})
}

func WithWebsocketPort(websocketPort string) Option {
	return applyOptionFunc(func(s *Server) error {
		s.websocketPort = websocketPort
		return nil
	})
}

func WithParsedIntent(parsed chan interface{}) Option {
	return applyOptionFunc(func(s *Server) error {
		s.parsed = parsed
		return nil
	})
}
