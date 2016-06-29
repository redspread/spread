package data

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

const (
	MaxObjectIDLen = 40
	MinObjectIDLen = 7
	ValidOIDChars  = "abcdef0123456789"

	OIDRegex   = `[^a-f0-9]+`
	PathRegex  = `[^a-zA-Z0-9./]+`
	FieldRegex = `[^a-zA-Z0-9./()]+`

	PathDelimiter  = "/"
	FieldDelimiter = "?"
)

// A SRI represents a parsed Spread Resource Identifier (SRI), a globally unique address for an object or field stored within a repository.
// This is represented as:
//
// 	treeish/path[?field]
type SRI struct {
	// Treeish is a Git Object ID to either a Commit or a Tree Git object.
	// The use of a Git OID (treeish) allows for any object or field to be addressed regardless if it is accessible.
	// The Object ID may be truncated down to a minimum of 7 characters.
	Treeish string

	// Path to the Spread Object being addressed. If omitted, SRI refers to treeish.
	// This will be traversed starting from the given Treeish.
	Path string

	// Field specifies a path to the field within the Object that is being referred to. If omitted, the SRI refers to the entire object.
	// Path must be given to have a Field.
	// Fieldpaths are specified by name using the character “.” to specify sub-fields.
	// Fieldpath of arrays are addressed using their 0 indexed position wrapped with parentheses.
	// The use of parentheses is due to restrictions in the syntax of URLs.
	Field string
}

// String returns a textual representation of the SRI which will be similar to the input.
func (s *SRI) String() string {
	str := s.Treeish
	if len(s.Path) > 0 {
		str += "/" + s.Path
	}
	if len(s.Field) > 0 {
		str += "?" + s.Field
	}
	return str
}

// ParseSRI parses rawsri into SRI struct.
func ParseSRI(rawsri string) (*SRI, error) {
	oid, path, field := parts(rawsri)
	var err error
	if oid, err = ParseOID(oid); err != nil {
		return nil, err
	}

	if path, err = ParsePath(path); err != nil {
		return nil, err
	}

	if field, err = ParseField(field); err != nil {
		return nil, err
	}

	return &SRI{
		Treeish: oid,
		Path:    path,
		Field:   field,
	}, nil
}

func parts(rawsri string) (oid, path, field string) {
	if len(rawsri) == 0 {
		return
	}

	// OID
	pathDelim := strings.Index(rawsri, PathDelimiter)
	// check if only OID
	if pathDelim == -1 {
		oid = rawsri
		return
	}
	oid = rawsri[:pathDelim]

	// Path
	fieldDelim := strings.LastIndex(rawsri, FieldDelimiter)
	if fieldDelim == -1 {
		if len(rawsri) > pathDelim+1 {
			path = rawsri[pathDelim+1:]
		}
		return
	}

	path, field = rawsri[pathDelim+1:fieldDelim], rawsri[fieldDelim+1:]
	return
}

func ParseOID(oidStr string) (string, error) {
	if len(oidStr) < MinObjectIDLen {
		return "", fmt.Errorf("git object ID was too short (%d chars), must be at least %d chars.", len(oidStr), MinObjectIDLen)
	} else if len(oidStr) > MaxObjectIDLen {
		return "", fmt.Errorf("git object ID was too long (%d chars), must be %d chars at most.", len(oidStr), MaxObjectIDLen)
	}

	// check has valid chars
	if regexp.MustCompile(OIDRegex).MatchString(oidStr) {
		return "", fmt.Errorf("invalid Treeish, invalid character in '%s' (only can contain '%s')", oidStr, ValidOIDChars)
	}
	return oidStr, nil
}

func ParsePath(pathStr string) (string, error) {
	if len(pathStr) == 0 {
		return "", nil
	}

	// check has valid chars
	if regexp.MustCompile(PathRegex).MatchString(pathStr) {
		return "", fmt.Errorf("invalid Path, invalid character in '%s' (must match regex '%s')", pathStr, PathRegex)
	}

	pathStr = filepath.Clean(pathStr)
	if pathStr[0] == '/' {
		if len(pathStr) == 1 {
			return "", nil
		}
		pathStr = pathStr[1:]
	}
	return pathStr, nil
}

func ParseField(fieldStr string) (string, error) {
	if len(fieldStr) == 0 {
		return "", nil
	}

	if fieldStr[0] == '.' {
		return "", errors.New("invalid Field: cannot begin with '.'")
	}

	if strings.Contains(fieldStr, "..") {
		return "", errors.New("invalid Field: cannot repeat '.'")
	}

	// check has valid chars
	if regexp.MustCompile(FieldRegex).MatchString(fieldStr) {
		return "", fmt.Errorf("invalid Field, invalid character in '%s' (must match regex '%s')", fieldStr, FieldRegex)
	}

	if err := checkIllegalParens(fieldStr); err != nil {
		return "", fmt.Errorf("invalid Field: %v", err)
	}

	return fieldStr, nil
}

func checkIllegalParens(fieldStr string) error {
	inParen := false
	for _, c := range fieldStr {
		if c == '(' {
			inParen = true
		} else if c == ')' {
			if !inParen {
				return errors.New("closed parenthese when one hasn't been opened")
			}
			inParen = false
		} else if inParen && !unicode.IsNumber(c) {
			return errors.New("only numeric characters can be used in parentheses")
		}
	}

	if inParen {
		return errors.New("unclosed parentheses")
	}

	return nil
}
