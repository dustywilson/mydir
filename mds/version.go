package mds

import (
	"sync"

	"github.com/satori/go.uuid"
)

// Version is a version of File
type Version struct {
	uuid uuid.UUID
	file *File
	sync.RWMutex
}

// NewVersion returns a Version of the referenced File
func (f *File) NewVersion() *Version {
	return NewVersion(f)
}

// NewVersion returns a Version of the referenced File
func NewVersion(file *File) *Version {
	v := &Version{
		uuid: uuid.NewV4(),
		file: file,
	}
	err := file.addVersion(v)
	if err != nil {
		panic(err)
	}
	return v
}

// Delete removes this Version
func (v *Version) Delete() error {
	if v != nil {
		if err := v.file.removeVersion(v); err != nil {
			return err
		}
	}
	return nil
}

// File returns parent File or nil
func (v *Version) File() *File {
	if v == nil {
		return nil
	}
	return v.file
}

// UUID returns UUID
func (v *Version) UUID() uuid.UUID {
	return v.uuid
}
