package image

// Image contains configuration necessary to deploy an image or if necessary, built it.
type Image struct {
	Name     string
	Tag      string
	Registry string
	Build    *Build
}

func (i Image) DockerName() string {
	// TODO: implement
	return i.Name
}

// FromString creates an Image using a string representation
func FromString(str string) (*Image, error) {
	// TODO: Implement
	return &Image{Name: str}, nil
}
