package flags

import "github.com/urfave/cli/v2"

var (
	AppEnv = "app-env"

	UIPort = "ui-port"

	// File Server Locations
	RootDirectory   = "root-directory"
	BleLogDirectory = "ble-log-directory"
	OTADirectory    = "ota-directory"
	UIDirectory     = "ui-directory"

	JdocsFilepath   = "jdocs-filepath"
	LicenesFilepath = "licenses-filepath"
	IntentsFilepath = "intents-filepath"

	// Speech To Speech
	STTModel                    = "stt-model"
	STTScorer                   = "stt-scorer"
	NumOfAudioStreamDispatchers = "num-of-audio-stream-dispatchers"

	// Extension Engine Stuff
	EscapePodExtender           = "escapepod-extender"
	EscapePodExtenderTarget     = "escapepod-extender-target"
	EscapePodExtenderDisableTLS = "escapepod-extender-disable-tls"

	DefaultIntentsFilepath = "default-intents-filepath"

	// TODO: add the vad timeout
	// vad-timeout

	EnableProfiler = "enable-profiler"
)

var AppFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    AppEnv,
		EnvVars: flagNamesToEnv(AppEnv),
	},
	&cli.StringFlag{
		Name:    UIPort,
		EnvVars: flagNamesToEnv(UIPort),
	},

	&cli.StringFlag{
		Name:    RootDirectory,
		EnvVars: flagNamesToEnv(RootDirectory),
	},
	&cli.StringFlag{
		Name:    BleLogDirectory,
		EnvVars: flagNamesToEnv(BleLogDirectory),
	},
	&cli.StringFlag{
		Name:    OTADirectory,
		EnvVars: flagNamesToEnv(OTADirectory),
	},
	&cli.StringFlag{
		Name:    UIDirectory,
		EnvVars: flagNamesToEnv(UIDirectory),
	},
	&cli.StringFlag{
		Name:    JdocsFilepath,
		EnvVars: flagNamesToEnv(JdocsFilepath),
	},
	&cli.StringFlag{
		Name:    LicenesFilepath,
		EnvVars: flagNamesToEnv(LicenesFilepath),
	},
	&cli.StringFlag{
		Name:    STTModel,
		EnvVars: flagNamesToEnv(STTModel),
	},
	&cli.StringFlag{
		Name:    STTScorer,
		EnvVars: flagNamesToEnv(STTScorer),
	},
	&cli.StringFlag{
		Name:    EscapePodExtender,
		EnvVars: flagNamesToEnv(EscapePodExtender),
	},
	&cli.StringFlag{
		Name:    EscapePodExtenderTarget,
		EnvVars: flagNamesToEnv(EscapePodExtenderTarget),
	},
	&cli.StringFlag{
		Name:    EscapePodExtenderDisableTLS,
		EnvVars: flagNamesToEnv(EscapePodExtenderDisableTLS),
	},
	&cli.StringFlag{
		Name:    IntentsFilepath,
		EnvVars: flagNamesToEnv(IntentsFilepath),
	},
	&cli.StringFlag{
		Name:    DefaultIntentsFilepath,
		EnvVars: flagNamesToEnv(DefaultIntentsFilepath),
	},
	&cli.BoolFlag{
		Name:    EnableProfiler,
		EnvVars: flagNamesToEnv(EnableProfiler),
		Aliases: []string{"p"},
	},
	&cli.IntFlag{
		Name:    NumOfAudioStreamDispatchers,
		EnvVars: flagNamesToEnv(NumOfAudioStreamDispatchers),
		Value:   1,
	},
}
