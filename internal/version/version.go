package version

// Version is set with ldflags -X
var Version = "This was not set"
var Build = "This was not set"
var UIVersion = "This was not set"
var UIBuild = "This was not set"

type VersionResponse struct {
	Version   string `json:"apiVersion"`
	Build     string `json:"apiBuild"`
	UIVersion string `json:"uiVersion"`
	UIBuild   string `json:"uiBuild"`
}

var DefaultVersionResponse = VersionResponse{
	Version:   Version,
	Build:     Build,
	UIVersion: UIVersion,
	UIBuild:   UIBuild,
}
