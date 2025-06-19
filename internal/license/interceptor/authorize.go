package interceptor

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

const (
	accessTokenMetadataKey = "anki-access-token"
)

func checkBypass(method string) bool {
	bypass := map[string]bool{
		"/tokenpb.Token/AssociatePrimaryUser": true,
		"/tokenpb.Token/RefreshToken":         true,
		"/jdocspb.Jdocs/WriteDoc":             true,
		"/jdocspb.Jdocs/ReadDocs":             true, // TODO: remove
		"/sttpb.STT/Parse":                    true,
		"/sttpb.STTManager/ListIntents":       true,
		"/sttpb.STTManager/GetIntent":         true,
		"/sttpb.STTManager/AddIntent":         true,
		"/sttpb.STTManager/DeleteIntent":      true,
		"/sttpb.STTManager/Match":             true,
		"/sttpb.STTManager/AddIntentsList":    true,
		"/logtracer.LogTracer/Trace":          true,
		"/interceptor.LicenseManager/Add":     true,
		"/interceptor.LicenseManager/List":    true,
		"/interceptor.LicenseManager/Delete":  true,
		"/bluey.Bluey/Init":                   true,
		"/bluey.Bluey/Close":                  true,
		"/bluey.Bluey/Scan":                   true,
		"/bluey.Bluey/Connect":                true,
		"/bluey.Bluey/SendPin":                true,
		"/bluey.Bluey/WifiScan":               true,
		"/bluey.Bluey/Status":                 true,
		"/bluey.Bluey/WifiConnect":            true,
		"/bluey.Bluey/Auth":                   true,
		"/bluey.Bluey/Configure":              true,
		"/bluey.Bluey/OTAStart":               true,
		"/bluey.Bluey/OTACancel":              true,
		"/bluey.Bluey/FetchLogs":              true,
		"/bluey.Bluey/ListLogs":               true,
		"/bluey.Bluey/DeleteLogs":             true,
	}

	if ok := bypass[method]; !ok {
		return false
	}
	return true
}

func (s *Interceptor) authorize(ctx context.Context, method string) bool {
	if checkBypass(method) {
		return true
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	t, ok := md[accessTokenMetadataKey]
	if !ok {
		return false
	}

	tok, err := jwt.Parse(t[0], func(t *jwt.Token) (interface{}, error) {
		return (s.key), nil
	})
	if err != nil {
		return false
	}

	robot, ok := checkToken(tok)
	if !ok {
		return false
	}

	if _, ok := s.bots[robot]; !ok {
		return false
	}

	return true
}

func checkToken(t *jwt.Token) (string, bool) {
	if t == nil {
		return "", false
	}

	if !t.Valid {
		return "", false
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}

	exp, ok := claims["expires"].(string)
	if !ok {
		return "", false
	}
	tz, _ := time.LoadLocation("UTC")
	expires, err := time.ParseInLocation(time.RFC3339, exp, tz)
	if err != nil {
		return "", false
	}

	if time.Now().UTC().After(expires) {
		return "", false
	}

	requestor, ok := claims["requestor_id"].(string)
	if !ok {
		return "", false
	}

	return requestor, true
}
