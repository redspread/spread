package data

import (
	"fmt"
)

const (
	MaxObjectIDLen = 40
	MinObjectIDLen = 7
)

// A SRL represents a parsed Spread Resource Locator (SRL), a globally unique address for an object or field stored within a repository.
// This is represented as:
//
// 	treeish/path[?fieldpath]
type SRL struct {
	// Treeish is a Git Object ID to either a Commit or a Tree Git object.
	// The use of a Git OID (treeish) allows for any object or field to be addressed regardless if it is accessible.
	// The Object ID may be truncated down to a minimum of 7 characters.
	Treeish string

	// Path to the Spread Object being addressed. If omitted, SRL refers to treeish.
	// This will be traversed starting from the given Treeish.
	Path string

	// Field specifies a path to the field within the Object that is being referred to. If omitted, the SRL refers to the entire object.
	// Path must be given to have a Field.
	// Fieldpaths are specified by name using the character “.” to specify sub-fields.
	// Fieldpath of arrays are addressed using their 0 indexed position wrapped with parentheses.
	// The use of parentheses is due to restrictions in the syntax of URLs.
	Field string
}

func (s *SRL) String() string {
	return ""
}

// ParseSRL parses rawsrl into SRL struct.
func ParseSRL(rawsrl string) (*SRL, error) {
	return nil, nil
}

var (
	// ErrOIDTooLong is returned when the length of the Git ObjectID is above MaxObjectIDLen
	ErrOIDTooLong = fmt.Errorf("git object ID was too long, must be %d chars at most.", MaxObjectIDLen)

	// ErrOIDTooShort is returned when the length of the Git ObjectID is below MinObjectIDLen
	ErrOIDTooShort = fmt.Errorf("git object ID was too short, must be at least %d chars.", MinObjectIDLen)
)
