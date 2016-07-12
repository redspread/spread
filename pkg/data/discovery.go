package data

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	// DiscoveryQueryParam is the query parameter that is appended to URLs to signal the request is looking for
	// repository information.
	DiscoveryQueryParam = "spread-get=1"

	// DiscoveryMetaName is the the 'name' of the <meta> tag that contains Spread package information.
	DiscoveryMetaName = "spread-ref"
)

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
// Verbose will print information to STDOUT.
func DiscoverPackage(packageName string, insecure, verbose bool) (packageInfo, error) {
	// first try HTTPS
	urlStr, res, err := fetch("https", packageName, verbose)
	if err != nil || res.StatusCode != 200 {
		if verbose {
			if err != nil {
				fmt.Fprint(Out, "https fetch failed")
			} else {
				fmt.Fprintf(Out, "ignoring https fetch with status code %d", res.StatusCode)
			}
		}
		// fallback to HTTP if insecure is allowed
		if insecure {
			urlStr, res, err = fetch("http", packageName, verbose)
		}
	}
	if err != nil {
		return packageInfo{}, err
	}

	// close body when done
	if res != nil {
		defer res.Body.Close()
	}

	if verbose {
		fmt.Fprintf(Out, "Parsing meta information from '%s' (status code %d)", urlStr, res.StatusCode)
	}

	pkgs, err := parseSpreadRefs(res.Body)
	if err != nil {
		return packageInfo{}, fmt.Errorf("could not parse for Spread references: %v", err)
	} else if len(pkgs) < 1 {
		return packageInfo{}, fmt.Errorf("no reference found at '%s'", urlStr)
	} else if len(pkgs) > 1 && verbose {
		fmt.Fprintf(Out, "found more than one reference at '%s', using first found", urlStr)
	}
	return pkgs[0], nil
}

// fetch retrieves the package using the given scheme and returns the response and a string of the URL of the fetch.
func fetch(scheme, packageName string, verbose bool) (string, *http.Response, error) {
	u, err := url.Parse(scheme + "://" + packageName)
	if err != nil {
		return "", nil, err
	}
	u.RawQuery = DiscoveryQueryParam
	urlStr := u.String()
	if verbose {
		fmt.Fprintf(Out, "fetching %s", urlStr)
	}

	res, err := httpClient.Get(urlStr)
	return urlStr, res, err
}

// parseSpreadRefs reads an HTML document from r and uses it to return information about the package.
// Information is currently stored in a <meta> tag with the name "spread-ref". Based on Go Get parsing code.
func parseSpreadRefs(r io.Reader) (pkgs []packageInfo, err error) {
	d := xml.NewDecoder(r)
	// only support documents encoded with ASCII
	d.CharsetReader = func(charset string, in io.Reader) (io.Reader, error) {
		switch strings.ToLower(charset) {
		case "ascii":
			return in, nil
		default:
			return nil, fmt.Errorf("cannot decode Spread package information encoded in %q", charset)
		}
	}
	d.Strict = false
	var t xml.Token
	for {
		if t, err = d.Token(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		if e, ok := t.(xml.StartElement); ok && strings.EqualFold(e.Name.Local, "body") {
			return
		}
		if e, ok := t.(xml.EndElement); ok && strings.EqualFold(e.Name.Local, "head") {
			return
		}
		e, ok := t.(xml.StartElement)
		if !ok || !strings.EqualFold(e.Name.Local, "meta") {
			continue
		}
		if attrValue(e.Attr, "name") != DiscoveryMetaName {
			continue
		}

		if f := strings.Fields(attrValue(e.Attr, "content")); len(f) == 2 {
			pkgs = append(pkgs, packageInfo{
				prefix:  f[0],
				repoURL: f[1],
			})
		}
	}
}

// attrValue returns the attribute value for the case-insensitive key
// `name', or the empty string if nothing is found.
func attrValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if strings.EqualFold(a.Name.Local, name) {
			return a.Value
		}
	}
	return ""
}
