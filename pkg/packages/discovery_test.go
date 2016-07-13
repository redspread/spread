package packages

import (
	"fmt"
	"net"
	"net/http"
	"testing"
)

var expandTestData = []struct {
	in    string
	out   string
	error bool
}{
	{
		in:    "",
		error: true,
	},
	// TODO: should disallow
	//{
	//	in:    "{",
	//	error: true,
	//},

	{
		in:  "hadoop",
		out: DefaultDomain + "/" + DefaultNamespace + "/hadoop",
	},
	{
		in:  "library/hadoop",
		out: DefaultDomain + "/library/hadoop",
	},
	{
		in:  "redspread.com/library/hadoop",
		out: "redspread.com/library/hadoop",
	},
	{
		in:  "github.com/redspread/hadoop",
		out: "github.com/redspread/hadoop",
	},
}

func TestExpandPackageName(t *testing.T) {
	for i, test := range expandTestData {
		pkg, err := ExpandPackageName(test.in)
		// check for issues with errors
		if err == nil && test.error {
			t.Errorf("test %d (input: %s): should have errored", i, test.in)
		} else if err != nil && !test.error {
			t.Errorf("test %d errored: %v", i, err)
			// check for issues with value
		} else if pkg != test.out {
			t.Errorf("test %d: expected '%s', got '%s'", i, test.out, pkg)
		}
	}
}

func TestDiscoverPackage(t *testing.T) {
	expected := packageInfo{
		prefix:  "redspread.com/halp",
		repoURL: "http://104.155.154.203/test.git",
	}
	server := NewDServer(t, expected)
	go server.Start()
	defer server.Stop()

	// here DNS would resolve from package name to host
	// we will only check if package data matches
	// another test should be used to ensure that prefix matches package
	importURL := fmt.Sprintf("%s/halp", server.Addr())
	actual, err := DiscoverPackage(importURL, true, true)
	if err != nil {
		t.Errorf("could not discover package: %v", err)
	} else if actual.repoURL != expected.repoURL {
		t.Errorf("repoURL did not match: \"%s\" (expected \"%s\")", actual.repoURL, expected.repoURL)
	} else if actual.prefix != expected.prefix {
		t.Errorf("prefix did not match: \"%s\" (expected \"%s\")", actual.prefix, expected.prefix)
	}
}

// testDServer mocks a server with discovery info.
type testDServer struct {
	info packageInfo
	net.Listener
	*testing.T
}

func NewDServer(t *testing.T, pkgInfo packageInfo) *testDServer {
	return &testDServer{
		info:     pkgInfo,
		Listener: randListener(t),
		T:        t,
	}
}

func (s *testDServer) Addr() string {
	return s.Listener.Addr().String()
}

func (s *testDServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handler)
	http.Serve(s.Listener, mux)
}

func (s *testDServer) Stop() error {
	return s.Listener.Close()
}

func (s *testDServer) handler(w http.ResponseWriter, r *http.Request) {
	msg := "<!DOCTYPE html><html><head><meta name=\"%s\" content=\"%s %s\"><title>Discovery Test Page</title></head><body><h1>Nothing to see here!</h1></body></html>"
	if _, err := fmt.Fprintf(w, msg, DiscoveryMetaName, s.info.prefix, s.info.repoURL); err != nil {
		s.Fatalf("Encountered error mocking discovery response: %v", err)
	}
}

// randomListener returns a listener for an available port.
func randListener(t *testing.T) *net.TCPListener {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	return lis
}
