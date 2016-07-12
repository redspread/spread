package data

import (
	"net/http"
)

// DiscoveryQueryParam is the query parameter that is appended to URLs to signal the request is looking for
// repository information.
const DiscoveryQueryParam = "spread-get=1"

// httpClient is a copy of DefaultClient for testing purposes.
var httpClient = http.DefaultClient

// packageInfo contains the data retrieved in the discovery process.
type packageInfo struct {
	// prefix is the package contained in the repo. It should be an exact match or prefix to the requested package name.
	prefix string
	// repoURL is the location of the repository where package data is stored.
	repoURL string
}

// DiscoverPackage uses the package name to fetch a Spread URL to the repository. Set insecure when HTTP is allowed.
func DiscoverPackage(packageName string, insecure bool) (packageInfo, error) {
	return packageInfo{}, nil
}

func fetch(scheme, packageName string) (*http.Response, error) {
	return nil, nil
}
