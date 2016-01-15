package image

// Image contains configuration necessary to deploy an image or if necessary, built it.
type Image struct {
	Name     string
	Tag      string
	Registry string
	Build    *Build
}
