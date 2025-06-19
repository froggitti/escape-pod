package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/DDLbots/api/go/jdocspb"
	"github.com/DDLbots/api/go/tokenpb"
	"github.com/DDLbots/escape-pod/internal/bluey"
	"github.com/DDLbots/escape-pod/internal/debug"
	"github.com/DDLbots/escape-pod/internal/flags"
	"github.com/DDLbots/escape-pod/internal/license/interceptor"
	"github.com/DDLbots/escape-pod/internal/logtracer"
	"github.com/DDLbots/escape-pod/internal/ui"
	"github.com/DDLbots/escape-pod/internal/version"
//	"github.com/DDLbots/escape-pod/internal/vte"

	flicensemanager "github.com/DDLbots/escape-pod/internal/license/interceptor/file"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	ep_license "github.com/DDLbots/internal-api/go/ep_licensepb"
	"github.com/DDLbots/internal-api/go/sttpb"
	"github.com/DDLbots/jdocs/pkg/jdocs"
	"github.com/DDLbots/jdocs/pkg/model/file"
	"github.com/DDLbots/saywhatnow/pkg/memory"
	"github.com/DDLbots/saywhatnow/pkg/saywhatnow"
	"github.com/DDLbots/saywhatnow/pkg/server/stt/dispatcher"
	"github.com/DDLbots/token/pkg/tokenservices/noop"
	"github.com/digital-dream-labs/vector-bluetooth/ble"

	chipperpb "github.com/DDLbots/api/go/chipperpb"
	chipper "github.com/DDLbots/chipper/pkg/server"
	saywhatnowClient "github.com/DDLbots/ddl-chipper/pkg/tts/saywhatnow"
	logger "github.com/DDLbots/go-logger"
	loggerf "github.com/DDLbots/go-logger/flags"
	serverf "github.com/DDLbots/go-server/flags"
	grpcp "github.com/DDLbots/go-server/grpc"
	ep_logtracer "github.com/DDLbots/internal-api/go/ep_logtracerpb"
	token "github.com/DDLbots/token/pkg/server"

	// FIXME: these need to be removed
	"log"
)

const (
	tlsCert = `-----BEGIN CERTIFICATE-----
MIIFBTCCAu2gAwIBAgIUFn4HuHNv0I48LRGei5Yf5DpKTNUwDQYJKoZIhvcNAQEL
BQAwgYExCzAJBgNVBAYTAlVTMRUwEwYDVQQIDAxQZW5uc3lsdmFuaWExEzARBgNV
BAcMClBpdHRzYnVyZ2gxGzAZBgNVBAoMEkRpZ2l0YWwgRHJlYW0gTGFiczEpMCcG
CSqGSIb3DQEJARYaYnJldHRAZGlnaXRhbGRyZWFtbGFicy5jb20wIBcNMjAxMjEw
MTgzNDU0WhgPMjIyMDEwMjMxODM0NTRaMIGbMQswCQYDVQQGEwJVUzEVMBMGA1UE
CAwMUGVubnN5bHZhbmlhMRMwEQYDVQQHDApQaXR0c2J1cmdoMRswGQYDVQQKDBJE
aWdpdGFsIERyZWFtIExhYnMxGDAWBgNVBAMMD2VzY2FwZXBvZC5sb2NhbDEpMCcG
CSqGSIb3DQEJARYaYnJldHRAZGlnaXRhbGRyZWFtbGFicy5jb20wggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQClNNnh/88lNt6VyoulGL8iDHae+o6wCFr9
glCsJv5WO2SFIGX7dZJ4CuCuo7GUOu7P74aqXQfoisDK8Hqr3iRftlZRl800N1Vr
khUNAx49obDBQctmjq0syMrTDDCIf/q3flmSqW0veT08zXOTgT7mdBG4e7Wuvwq1
NagU5jv35IpxWnNTdrU7HWKr8DHogeP1eI/09KR3bA9xtwsEh8USdguHxLY3V6bB
xlC6vN3gItNf0HsinSjU2qRW47P5m5aoNfYmmjvjUpcvH/qHjcCvvAMYZ7+31jWd
4kUhbpWW4zYgIa37fo3dQiquit7IfLUSyT7Bt9eO3eLQCDPdNm8FAgMBAAGjVzBV
MB8GA1UdIwQYMBaAFLEGVfbWGssezxPrNlderfGJaU0eMAkGA1UdEwQCMAAwCwYD
VR0PBAQDAgTwMBoGA1UdEQQTMBGCD2VzY2FwZXBvZC5sb2NhbDANBgkqhkiG9w0B
AQsFAAOCAgEAJsVRYGxoqYhp6vc1jLUVqInH9RPF17B9CwHg7ECyvt/z+IKSCGEz
xvYx4/PwrSMYqss9MiOpWC/ZNI6rU7ruks/KQ15DYPAZJl06hvz+LYuLrphdtaO+
CnLUzpcn4MQ9NdxsOFqWiZelIRdEAiq/rHGpsOOSdIfARvoNN3uOyvhmDCFtzUTi
0V/bA7AqshEHArWLCVI+UOnDpgB7bYibTPBS1RFpHskPR0ynsEEsMQQMM0howwBG
iYi7AqdlrjHxo1J6AzX31EReXIi/v1HK+b+aCoQYKPqiAJ6kkYt3ffnuQ3W7pCnd
wr3iM/g+TkWNqRf9f9jxmG4ozwKI9q7/vQdNkJfu534T3ohmrgzGSKQ99IX3VN9P
ePFrxsy8w3OcSflv/Vjo/nHTUETnuB4EdeUOPNf/7euOtznUFhsyLmDuFHA6a2yU
PpI+g2HNgQUbJOtIRV3owc1ihnXwgeuj8Xt2eZvt3VaUE0333uC6Yc0EngIqf4gA
CQWbnsQM8HMh/+npFzbZlcM4RZoFChObb4O0FqVlYjD7/xnJE78dqSLwP8VVSD4g
y+WRDN4PJXgGZw4omtpC5OTReDOZmmZJUN2YTla1SZXyJcZv/x+P1th0qnmZ4Hf3
E4oNSE7wWX9/lZRTb3fegWL4AlW3IsosyIpYTHZ8Qqy1YfHo3zdxXqU=
-----END CERTIFICATE-----`

	tlsKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEApTTZ4f/PJTbelcqLpRi/Igx2nvqOsAha/YJQrCb+VjtkhSBl
+3WSeArgrqOxlDruz++Gql0H6IrAyvB6q94kX7ZWUZfNNDdVa5IVDQMePaGwwUHL
Zo6tLMjK0wwwiH/6t35ZkqltL3k9PM1zk4E+5nQRuHu1rr8KtTWoFOY79+SKcVpz
U3a1Ox1iq/Ax6IHj9XiP9PSkd2wPcbcLBIfFEnYLh8S2N1emwcZQurzd4CLTX9B7
Ip0o1NqkVuOz+ZuWqDX2Jpo741KXLx/6h43Ar7wDGGe/t9Y1neJFIW6VluM2ICGt
+36N3UIqroreyHy1Esk+wbfXjt3i0Agz3TZvBQIDAQABAoIBAHeIzRm72OrJT7Y8
LlxPkoQVVoLjMfjmoseI0cwuDprgMHQuo/uU71ySKk3SPTvOhFrJqbt8wqscMjDk
XS4b9l+Wc9BnsN9WJiVGNpsKpYfchSLf80cKdvzPcAnSaQ9q4kKAVllK46iU5Z0n
3rdcreFbHDNKt4Nv0VSaNTqh98P9W5CK1pIZbGHJGqFaE5VtUxvBzUGc7jAXj+eD
I8Ajum4AG3ZcJBBMiGCX6TzCoHZYBRUicLRvKkANDqNP4/L9wqbAlNmpP7IUD/yC
OjlYlEKjEXyk5hmv/QHRqN2yK2qpCBddTxmwiXAa3BSu17hTv0zUlDTR5TjO2TEG
tsNDnaECgYEA2q547mBRIt3/D5VrjXZMOGYFgt306nzFg3XJYM5ArOxR99tZFWRv
oHdwBMIvVI6RgHhw6LuferZI4n5T0UUaSUUuHXvyD9T20S8buDrC6TsmNruymMM+
jvry6yGF3+AMlIHrayq0U6FHBrhfGsH91Is1/RF6Hu4Rpfa85GVASYkCgYEAwWY3
P/BVho5CkEu1H3PMH4HiDaYo/VrA+pCaNe8KH7xQJ5IOL4IywT2W2YVKsxU9fwmE
ZEmXJm7nzDLYwsyuk9pgxSjud2SVBUvnnBDtbX0sNHB/Qi7nLjbs5ZjHISH2RNhP
daOGI33ZRKXOfvbsXUNSA2r8fR1bE3JvA1YnJp0CgYBEfLH5Dgc7IUWZbtVxR2RV
oXYGZ1cl/Q+qvT/lZpMQ1S5Srsq2jW78VYuqodpK5B+jmZTa/q/SsbYf4SqE9txl
qBnqOAA2fx8RomxPBXA3tUOhjqU/fJ5iDyv3Ade4pqWp+Qpu1MAHFRJ2g1WdvrWt
VDADYu7ZMvwp+x1rdl5s6QKBgAWGNe3Nn6PITH5yqynK1PnRa/OX23PhM8H0f3Mq
8M8XQfLfaShSP8DlUXnFJO0YnjkSvIVg1MB0Soq6qRZnYlU216zKDoW6iccs8+Cx
WxbVjH2y+O+bB196kim8w3Ne1PoCc8KYeSxqW9pqIgveYcIIOj9+vteUDxXvHtyp
iVTBAoGASI0u1CTISUMuTn2Dknd750UcKXkfFa1yTWG2NbsB64CFx8vRcM78rHVo
bmgD2V1X25m3Pdqk1OD1Q2+ohb+p7ss+PygV2tFMryaJ1MbaPnoU+7DewptFJ7Do
wMmnGhTrRbpZouBkZx2OUXhOkhKui63c82eAF8h60zJosbPXQ8w=
-----END RSA PRIVATE KEY-----`

	escapepodRootPEM = `-----BEGIN CERTIFICATE-----
MIIF5zCCA8+gAwIBAgIUd7+iytZTXyFdQRwAtMc28TtXutQwDQYJKoZIhvcNAQEL
BQAwgYExCzAJBgNVBAYTAlVTMRUwEwYDVQQIDAxQZW5uc3lsdmFuaWExEzARBgNV
BAcMClBpdHRzYnVyZ2gxGzAZBgNVBAoMEkRpZ2l0YWwgRHJlYW0gTGFiczEpMCcG
CSqGSIb3DQEJARYaYnJldHRAZGlnaXRhbGRyZWFtbGFicy5jb20wIBcNMjAxMjEw
MTgxODQwWhgPMjIyMDEwMjMxODE4NDBaMIGBMQswCQYDVQQGEwJVUzEVMBMGA1UE
CAwMUGVubnN5bHZhbmlhMRMwEQYDVQQHDApQaXR0c2J1cmdoMRswGQYDVQQKDBJE
aWdpdGFsIERyZWFtIExhYnMxKTAnBgkqhkiG9w0BCQEWGmJyZXR0QGRpZ2l0YWxk
cmVhbWxhYnMuY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAvtk7
ccn6Y/1xQlynkvvon86fgj4wQuqYbKmUOeNZkNv80QhW+WcVEGuNkDctXNGAwTZI
kBM4nKl1CEq24OgT6I6skYMqEOEGJdaGld08LJwlf4PHG9V7qdNOiWv36GpK2blp
PB0UuAgxVgEF5DC6V+RWxbFYfaQcfa9OMeMCvzaL5yQdK1gVYZLjWS3yeFn33oh9
H2mhTelW/24+UXV7buneAvyEYXJEMdRmEjKTtsWktMGDp6d9Eg6cZjwQzJrlhuQj
FpHdM+FdpFkQ/22wGi+L27f2m+lCEE+G37bkPB+L2EXuYNSDv9mVSyqb6NbqpwX5
sPyEzPitDROP7+sUgeTjg83biQAjQyedKxOGab8BwKSJsbmGcodAgHS/sMdDAuV/
Uw2iy1zN8uQaOIvXxWCWpvXAcwA4Aro2ruEkMR7h/+Cc26dt6EQ2VrSyxVEzq+Sz
5fxAYtwFAsd023Tpztu+UewmNOhQHPb1Xj99/nW2lR4fH0kSfsyBFKlKXgFwDSZD
vbIQ0ZEEFrIpqnsgmwc/C2w1ztD22R5Jnfrb8s23xj4awv3/aPl+4NeVBwQMgxS6
8n/JWRuFvtIDj4gQWbKq5wjIdcGaPEhHRYsLiWr0NCDwn/R+CCdNIY8SCtMllrP8
PUBxwnlQ5NfBw4azwBO9tW47iD0pbVrxx48UfNECAwEAAaNTMFEwHQYDVR0OBBYE
FLEGVfbWGssezxPrNlderfGJaU0eMB8GA1UdIwQYMBaAFLEGVfbWGssezxPrNlde
rfGJaU0eMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIBAKsD9769
YV23Ak3RBd5IBk1LNxwQLbp4oLes4k4XdSNRWEddMFOLHcfJgtjAugoJwRza1OIz
7I/2Pe0AhUpvtMqDn9zquDACx0Eno2/3xOGdjxCH8rXTF1Hpc/NEWjyzhJIj4Oz+
DHtNicWkJG4e/yLuHD6E4qlcqYW8JEQt3j156IZIQgJxkkSDHS7sxpDMVYX06AKd
b7Msmagde9gm4jeFIOqa3pCoGABUdPCG4poTIwAmuobDDji08bo7PBYFWvB8Ye6y
kcfVML0OzWRJgcG8AWltz1rA4GSVG91oSr0DXNyUTHR0GFc7OqX26KkxuHV6EoDz
SkqvwRxceOjTf49JobHlFbKuj6zN+t608ZEVPUq8g/uG68ecat9MBnH3mbenMFZ/
7AVW3dpEbxqrIK8wLGleoQQ2+l1Jkfu7LNbPqiMUIWF02FNK5PoD9W5KFl/geYqD
3cA4AKYtrHnvEql6DABZILfmEnkfxYipZkwBaOoUSLuRSLP/kdNwkZjdg8O+U9qA
ac9sgHFZ/wegv3etK46TtLTb0vrMWm5pFFcFRCiYeqNpAhAwEb1b8YQR0sjT9rlL
o6yMR+1kuDXYzB1REuuTOZHGDj7Zt1yiHv+nr1dPY/EBV+69is4e2lIXyHlI2ov8
HEtdAlxwZiakfw9zndxIweAw9YapnyT+4xp2
-----END CERTIFICATE-----`

	signingKey = "testKey12345"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
	}()

	if err := NewApp().RunContext(ctx, os.Args); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		log.Fatal(err)
	}
}

func NewApp() *cli.App {
	return &cli.App{
		Name:  "escape-pod",
		Usage: "",
		Flags: flags.Join(
			flags.AppFlags,
			flags.JdocsFlags,
			flags.VersionFlags,
			serverf.GRPCFlags,
			loggerf.LogFlags,
		),
		Action: func(c *cli.Context) error {

			if c.Bool(flags.VersionKey) {
				fmt.Fprintf(os.Stdout, "api version: %v\nui version: %s\n", version.Version, version.UIVersion)
				return nil
			}
			if c.Bool(flags.BuildKey) {
				fmt.Fprintf(os.Stdout, "api build: %s\nui build: %s\n", version.Build, version.UIBuild)
				return nil

			}
			if c.Bool(flags.BuildKey) {
				fmt.Fprintf(os.Stdout,
					"api version: %s\nui version: %s\napi build: %s\nui build: %s\n",
					version.Version, version.UIVersion, version.Build, version.UIBuild,
				)
				return nil
			}
			// Get out needed env vars
			rootDir := c.String(flags.RootDirectory)
			bleLogsDir := c.String(flags.BleLogDirectory)
			uiDir := c.String(flags.UIDirectory)
			otaDir := c.String(flags.OTADirectory)
			licensesFilepath := c.String(flags.LicenesFilepath)
			intentsFilepath := c.String(flags.IntentsFilepath)
			defaultIntentsFilepath := c.String(flags.DefaultIntentsFilepath)
			grpcPort := c.String(serverf.GRPCPort)
			uiPort := c.String(flags.UIPort)

			enableProfiler := c.Bool(flags.EnableProfiler)
			if enableProfiler {
				profiler := debug.NewDebugServer(fmt.Sprintf(":6060"))
				go func() {
					profiler.ListenAndServe()
				}()
			}

			if !c.IsSet(serverf.GRPCPubCert) && !c.IsSet(serverf.GRPCPrivCert) {
				////
				//
				// TLS CERTS GET SET IF NOT ALREADY SET; SUBSEQUENT READS STILL NEED TO BE FROM CONTEXT
				//
				////
				c.Set(serverf.GRPCPubCert, tlsCert)
				c.Set(serverf.GRPCPrivCert, tlsKey)
			}
			pubCert := c.String(serverf.GRPCPubCert)
			privCert := c.String(serverf.GRPCPrivCert)

			certs, err := grpcp.ParseCertificates(pubCert, privCert)
			if err != nil {
				return err
			}

			cp, err := x509.SystemCertPool()
			if err != nil {
				return fmt.Errorf("system cert pool: %v", err)
			}
			if ok := cp.AppendCertsFromPEM([]byte(escapepodRootPEM)); !ok {
				return errors.New("append to system ca")
			}

			creds := credentials.NewTLS(&tls.Config{ServerName: "escapepod.local", RootCAs: cp, Certificates: []tls.Certificate{certs}})

			// This needs to go to the grpc server config for the WithGatewayDialOptions
			tlsCreds := grpc.WithTransportCredentials(creds)
			// Create GRPC Dailer
			conn, err := grpc.DialContext(c.Context, fmt.Sprintf(":%s", c.String(serverf.GRPCPort)), tlsCreds)
			if err != nil {
				return fmt.Errorf("setup grpc conn: %v", err)
			}

			// SETUP LOGGER
			logr, err := setupLogger(c)
			if err != nil {
				return fmt.Errorf("setup logger: %v", err)
			}

			// SETUP BLUETOOTH SERVICE
			// This is the channel that the BLE library uses to provide status for the UI
			dlstatchan := make(chan ble.StatusChannel)
			blueyService, err := bluey.New(filepath.Join(rootDir, bleLogsDir), dlstatchan)
			if err != nil {
				return fmt.Errorf("setup bluetooth: %v", err)
			}

			// SETUP LICENSE MANAGER
			lm, err := flicensemanager.New(
				flicensemanager.WithFilePath(filepath.Join(rootDir, licensesFilepath)),
				flicensemanager.WithDebugger(logr),
			)
			if err != nil {
				return fmt.Errorf("new license manager: %v", err)
			}

			// SETUP INTERCEPTOR/LICENSE SERVICE
			licenseService, err := interceptor.New(signingKey, lm)
			if err != nil {
				return fmt.Errorf("new interceptor: %v", err)
			}

			// SETUP TOKEN SERVICE
			tok, err := noop.New(
				noop.WithSigningKey(signingKey),
			)
			if err != nil {
				return fmt.Errorf("new tokenizer: %v", err)
			}

			tokenService, err := token.New(
				token.WithLogger(logr.With("entity", "token_service")),
				token.WithTokenizer(tok),
				token.WithJdocsClient(
					jdocspb.NewJdocsClient(conn),
				),
			)
			if err != nil {
				return fmt.Errorf("setup token service: %v", err)
			}

			// This interface is used be json.Marshal to get the values to json to send to the front end
			parsed := make(chan interface{})

			// SETUP UI SERVER
			uiService, err := ui.New(
				ui.WithX509KeyPairString(tlsCert, tlsKey),
				ui.WithRootDir(rootDir),
				ui.WithLogsDir(bleLogsDir),
				ui.WithOTADir(otaDir),
				ui.WithUIDir(uiDir),
				ui.WithPort(uiPort),
				ui.WithWebsocketPort(uiPort),
				ui.WithDLStatusChannel(dlstatchan),
				ui.WithParsedIntent(parsed),
			)
			if err != nil {
				return fmt.Errorf("setup ui server: %v", err)
			}

			// SETUP LOGTRACER
			logtrcr := &logtracer.LogTracer{}

			//// SAYWHATNOW SERVICE: Converts Text To Speech Then Parses to Match Intent
			////// MEMORY
			logr.Debug(fmt.Sprintf("default intents file path: %s", defaultIntentsFilepath))

			memoryIntents, err := memory.New(
				memory.WithDefaultIntentsFilepath(defaultIntentsFilepath),
				memory.WithFilePath(filepath.Join(rootDir, intentsFilepath)),
				// memory.WithDebugger(logr),
			)
			if err != nil {
				return fmt.Errorf("setup memory intents manager: %v", err)
			}

			dsp, err := dispatcher.New(
				c.Int(flags.NumOfAudioStreamDispatchers),
				c.String(flags.STTModel),
				c.String(flags.STTScorer),
				// dispatcher.WithDebugger(logr),
			)
			if err != nil {
				return fmt.Errorf("setup dispatcher: %v", err)
			}

			saywhatnowService, err := saywhatnow.New(
				saywhatnow.WithDispatcher(dsp),
				saywhatnow.WithIntentManager(memoryIntents),
				saywhatnow.WithIntentMatcher(memoryIntents),
				saywhatnow.WithLogger(logr),
				saywhatnow.WithParsedChan(parsed),
			)

			// SETUP CHIPPER SERVICE
			//// SETUP CHIPPER STT CLIENT !!NOTE!! THIS IS ACCTUALLY THE SPEECH TO INTENT CLIENT
			p, err := saywhatnowClient.New(
				saywhatnowClient.WithSTTClient(sttpb.NewSTTClient(conn)),
				// TODO: get all of these config out of viperize
				saywhatnowClient.WithViper(),
				// saywhatnowClient.WithFunctionServer(getFunctionServer()),
			)
			if err != nil {
				return fmt.Errorf("new saywhatnow client %v", err)
			}

			chipperService, err := chipper.New(
				// chipper.WithLogger(log.Base()),
				chipper.WithIntentProcessor(p),
				chipper.WithKnowledgeGraphProcessor(p),
				chipper.WithIntentGraphProcessor(p),
			)
			if err != nil {
				return fmt.Errorf("setup chipper service %v", err)
			}
			logr.Debug("jdocs file path:", filepath.Join(c.String(flags.RootDirectory), c.String(flags.JdocsFilepath)))
			// SETUP JDOCS SERVICE
			fileDB, err := file.New(
				file.WithFilename(filepath.Join(c.String(flags.RootDirectory), c.String(flags.JdocsFilepath))),
				file.WithDebugger(logr),
			)
			if err != nil {
				return fmt.Errorf("setup jdocs ds: %v", err)
			}

			jdocsService, err := jdocs.New(
				jdocs.WithDB(fileDB),
				// FIXME: find out if it need authorization
				jdocs.WithAuthorization(false),
				jdocs.WithLogger(logr),
			)
			if err != nil {
				return fmt.Errorf("setup jdocs service: %v", err)
			}

			grpcOpts := setupBaseGRPC(c, logr,
				uiService,
				func(s *grpc.Server) {
					// These are needed for the Vector Robot
					// Jdocs: note needed tls
					jdocspb.RegisterJdocsServer(s, jdocsService)
					// Chipper: note needed tls
					chipperpb.RegisterChipperGrpcServer(s, chipperService)
					// SayWhatNow Parser: note needed tls
					sttpb.RegisterSTTServer(s, saywhatnowService)
					// LogTracer: note needed tls
					ep_logtracer.RegisterLogTracerServer(s, logtrcr)
					// Token: note needed tls
					tokenpb.RegisterTokenServer(s, tokenService)

					// These are needed for the UI
					// SayWhatNow Manager:
					sttpb.RegisterSTTManagerServer(s, saywhatnowService)
					// Bluetooth:
					ep_bluetooth.RegisterBlueyServer(s, blueyService)
					// License:
					ep_license.RegisterLicenseManagerServer(s, licenseService)
				},
				[]grpcp.GatewayServiceHandler{
					ep_bluetooth.RegisterBlueyHandler,

					ep_license.RegisterLicenseManagerHandler,

					sttpb.RegisterSTTManagerHandler,
				},
				[]grpc.StreamServerInterceptor{
					licenseService.StreamServerInterceptor(),
				},
				[]grpc.UnaryServerInterceptor{
					licenseService.UnaryServerInterceptor(),
				},
			)

			grpcOpts = append(grpcOpts,
				grpcp.WithPrivCert(c.String(serverf.GRPCPrivCert)),
				grpcp.WithPubCert(c.String(serverf.GRPCPubCert)),
				grpcp.WithDialCerts([]tls.Certificate{certs}),
				grpcp.WithGatewayDialOptions(tlsCreds),
				grpcp.WithInsecureSkipVerify(),
				grpcp.WithClientAuthType(tls.RequestClientCert),
			)

			grpcServer, err := grpcp.New(append([]grpcp.Option{
				grpcp.WithPort(grpcPort),
				grpcp.WithTLS(true),
			}, grpcOpts...)...)
			if err != nil {
				return fmt.Errorf("new grpc server: %v", err)
			}

			uiServer, err := grpcp.New(append([]grpcp.Option{
				grpcp.WithGatewayAddr("0.0.0.0", grpcPort),
				grpcp.WithPort(uiPort),
			}, grpcOpts...)...)
			if err != nil {
				return fmt.Errorf("new ui server: %v", err)
			}

			// run the http servers in an error group so that if one goes down they can both to down
			ctx, cancel := context.WithCancel(c.Context)
			eg, ctx := errgroup.WithContext(ctx)

			eg.Go(func() error {
				defer cancel()
				if err := uiServer.StartWithContext(c.Context); err != nil {
					return fmt.Errorf("start ui server: %v", err)
				}
				return nil
			})

			eg.Go(func() error {
				defer cancel()
				if err := grpcServer.StartWithContext(c.Context); err != nil {
					return fmt.Errorf("start grpc server: %v", err)
				}
				return nil
			})

			return eg.Wait()
		},
	}
}

func setupBaseGRPC(
	c *cli.Context,
	logr logger.Logger[*slog.Logger, slog.Level],
	handler grpcp.Handler,
	registerServiceFunc func(*grpc.Server),
	gatewayServiceHandlers []grpcp.GatewayServiceHandler,
	streamInterceptors []grpc.StreamServerInterceptor,
	unaryInterceptors []grpc.UnaryServerInterceptor,
) []grpcp.Option {

	opts := []grpcp.Option{
		grpcp.WithRegisterService(registerServiceFunc),
		grpcp.WithServerOptions(
			grpc.ChainStreamInterceptor(
				append([]grpc.StreamServerInterceptor{
					logger.LoggingStreamServerInterceptor(logr),
				}, streamInterceptors...)...,
			),
			grpc.ChainUnaryInterceptor(append([]grpc.UnaryServerInterceptor{
				logger.LoggingUnaryServerInterceptor(logr),
			}, unaryInterceptors...)...),
		),
		grpcp.WithGatewayServiceHandlers(gatewayServiceHandlers...),
		grpcp.WithLogger(logr),
	}

	if handler != nil {
		opts = append(opts, grpcp.WithHandler(handler))
	}

	return opts
}

func setupLogger(c *cli.Context) (logger.Logger[*slog.Logger, slog.Level], error) {
	// Setup Logger
	logOpts := []logger.Option{
		logger.WithEnv(c.String(loggerf.LogEnv)),
		logger.WithLevel(c.String(loggerf.LogLevel)),
		logger.WithLogStacktrace(c.Bool(loggerf.LogStacktrace)),
	}
	logr, err := logger.New(logOpts...)
	if err != nil {
		return nil, fmt.Errorf("new logger: %v", err)
	}
	return logr, nil
}
