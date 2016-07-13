package packages

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"rsprd.com/spread/pkg/config"
)

var domainRegexp *regexp.Regexp

func init() {
	domainRegexp = regexp.MustCompile(DomainRegexpStr)
}

const (
	// DiscoveryQueryParam is the query parameter that is appended to URLs to signal the request is looking for
	// repository information.
	DiscoveryQueryParam = "spread-get=1"

	// DiscoveryMetaName is the the 'name' of the <meta> tag that contains Spread package information.
	DiscoveryMetaName = "spread-ref"

	// DefaultDomain is the domain assumed for packages if one is not given.
	DefaultDomain = "redspread.com"

	// DefaultNamespace is the namespace used if no domain and namespace is given.
	DefaultNamespace = "library"

	// DomainRegexpStr is a regular expression string to validate domains.
	DomainRegexpStr = "^([a-z0-9]+(-[a-z0-9]+)*\\.)+[a-z]{2,}$"
)

// ExpandPackageName returns a retrievable package name for packageName by adding the Redspread domain where a domain isn't specified.
func ExpandPackageName(packageName string) (string, error) {
	if len(packageName) == 0 {
		return "", errors.New("empty package name")
	}

	// if single segment, assume domain and namespace
	pkgArr := strings.Split(packageName, "/")
	if len(pkgArr) == 1 {
		return DefaultDomain + "/" + DefaultNamespace + "/" + packageName, nil
	}

	if isDomain(pkgArr[0]) {
		return packageName, nil
	}

	// if no domain, assume it
	return DefaultDomain + "/" + packageName, nil
}

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
				fmt.Fprint(config.Out, "https fetch failed")
			} else {
				fmt.Fprintf(config.Out, "ignoring https fetch with status code %d", res.StatusCode)
			}
		}
		// fallback to HTTP if insecure is allowed
		if insecure {
			if res != nil {
				res.Body.Close()
			}
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
		fmt.Fprintf(config.Out, "Parsing meta information from '%s' (status code %d)", urlStr, res.StatusCode)
	}

	pkgs, err := parseSpreadRefs(res.Body)
	if err != nil {
		return packageInfo{}, fmt.Errorf("could not parse for Spread references: %v", err)
	} else if len(pkgs) < 1 {
		return packageInfo{}, fmt.Errorf("no reference found at '%s'", urlStr)
	} else if len(pkgs) > 1 && verbose {
		fmt.Fprintf(config.Out, "found more than one reference at '%s', using first found", urlStr)
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
		fmt.Fprintf(config.Out, "fetching %s", urlStr)
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

// isDomain returns whether given string is a domain.
// It first checks the TLD, and then uses a regular expression.
func isDomain(s string) bool {
	if strings.HasSuffix(s, ".") {
		s = s[:len(s)-1]
	}

	split := strings.Split(s, ".")
	tld := split[len(split)-1]

	if len(tld) < 2 {
		return false
	}

	s = strings.ToLower(s)
	return domainRegexp.MatchString(s)
}
