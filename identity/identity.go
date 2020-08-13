// Package identity provides build identity data set at compile time by govvv.
// See https://github.com/ahmetb/govvv for details.
package identity

// Define variables for build metadata set at compile time by govvv.
var (
	BuildDate  string
	GitBranch  string
	GitCommit  string
	GitState   string
	GitSummary string
	Version    string
)

// Metadata includes useful data about a build.
type Metadata struct {
	BuildDate  string `json:"buildDate"`
	GitBranch  string `json:"gitBranch"`
	GitCommit  string `json:"gitCommit"`
	GitState   string `json:"gitState"`
	GitSummary string `json:"gitSummary"`
	Version    string `json:"version"`
}

// GetMetadata returns the build metadata.
func GetMetadata() Metadata {
	return Metadata{
		BuildDate:  BuildDate,
		GitBranch:  GitBranch,
		GitCommit:  GitCommit,
		GitState:   GitState,
		GitSummary: GitSummary,
		Version:    Version,
	}
}
