package data

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	pb "rsprd.com/spread/pkg/spreadproto"
)

const (
	MaxObjectIDLen = 40
	MinObjectIDLen = 7
	ValidOIDChars  = "abcdef0123456789"

	PathDelimiter  = "/"
	FieldDelimiter = "?"
)

var (
	OIDRegex   = regexp.MustCompile(`[^a-f0-9]+`)
	PathRegex  = regexp.MustCompile(`[^a-zA-Z0-9./]+`)
	FieldRegex = regexp.MustCompile(`[^a-zA-Z0-9./()]+`)
)

// A SRI represents a parsed Spread Resource Identifier (SRI), a globally unique address for an document or field stored within a repository.
// This is represented as:
//
// 	treeish/path[?field]
type SRI struct {
	// Treeish is a Git Object ID to either a Commit or a Tree Git object.
	// The use of a Git OID (treeish) allows for any document or field to be addressed regardless if it is accessible.
	// The Object ID may be truncated down to a minimum of 7 characters.
	// A single character “*” indicates a relative reference, this intentionally can’t be formed into a URL.
	Treeish string

	// Path to the Spread Document being addressed. If omitted, SRI refers to treeish.
	// This will be traversed starting from the given Treeish.
	Path string

	// Field specifies a path to the field within the Document that is being referred to. If omitted, the SRI refers to the entire document.
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

// Proto returns the protobuf representation of an SRI
func (s *SRI) Proto() *pb.SRI {
	return &pb.SRI{
		Treeish: s.Treeish,
		Path:    s.Path,
		Field:   s.Field,
	}
}

// IsTreeish is true if identifier points to tree.
func (s *SRI) IsTree() bool {
	return !s.IsDocument() && !s.IsField()
}

// IsDocument is true if points to document.
func (s *SRI) IsDocument() bool {
	return len(s.Field) == 0 && len(s.Path) > 0
}

// IsField is true if points to field.
func (s *SRI) IsField() bool {
	return len(s.Field) > 0 && len(s.Path) > 0
}

// ParseSRI parses rawsri into SRI struct.
func ParseSRI(rawsri string) (*SRI, error) {
	oid, path, field := parts(rawsri)
	var err error
	if oid, err = ValidateOID(oid); err != nil {
		return nil, err
	}

	if path, err = ValidatePath(path); err != nil {
		return nil, err
	}

	if field, err = ValidateField(field); err != nil {
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

func ValidateOID(oidStr string) (string, error) {
	if len(oidStr) == 1 && oidStr[0] == '*' {
		return "*", nil
	} else if len(oidStr) < MinObjectIDLen {
		return "", fmt.Errorf("git object ID was too short (%d chars), must be at least %d chars.", len(oidStr), MinObjectIDLen)
	} else if len(oidStr) > MaxObjectIDLen {
		return "", fmt.Errorf("git object ID was too long (%d chars), must be %d chars at most.", len(oidStr), MaxObjectIDLen)
	}

	// check has valid chars
	if OIDRegex.MatchString(oidStr) {
		return "", fmt.Errorf("invalid Treeish, invalid character in '%s' (only can contain '%s')", oidStr, ValidOIDChars)
	}
	return oidStr, nil
}

func ValidatePath(pathStr string) (string, error) {
	if len(pathStr) == 0 {
		return "", nil
	}

	// check has valid chars
	if PathRegex.MatchString(pathStr) {
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

func ValidateField(fieldStr string) (string, error) {
	if len(fieldStr) == 0 {
		return "", nil
	}

	if fieldStr[0] == '.' {
		return "", errors.New("invalid Field: cannot begin with '.'")
	}

	if fieldStr[len(fieldStr)-1] == '.' {
		return "", errors.New("invalid Field: cannot end with '.'")
	}

	if strings.Contains(fieldStr, "..") {
		return "", errors.New("invalid Field: cannot repeat '.'")
	}

	// check has valid chars
	if FieldRegex.MatchString(fieldStr) {
		return "", fmt.Errorf("invalid Field, invalid character in '%s' (must match regex '%s')", fieldStr, FieldRegex)
	}

	if err := checkIllegalParens(fieldStr); err != nil {
		return "", fmt.Errorf("invalid Field: %v", err)
	}

	return fieldStr, nil
}

func checkIllegalParens(fieldStr string) error {
	inParen := -1
	for i, c := range fieldStr {
		if c == '(' {
			inParen = i
		} else if c == ')' {
			if inParen == -1 {
				return errors.New("closed parenthese when one hasn't been opened")
			} else if inParen == i-1 {
				return errors.New("must specify array position, cannot have '()'")
			}
			inParen = -1
		} else if inParen != -1 && !unicode.IsNumber(c) {
			return errors.New("only numeric characters can be used in parentheses")
		}
	}

	if inParen != -1 {
		return errors.New("unclosed parentheses")
	}

	return nil
}
