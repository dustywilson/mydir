package mds

import (
	"errors"
	"sync"

	"github.com/satori/go.uuid"
)

// File is an entry that may have Version entries
type File struct {
	uuid      uuid.UUID
	name      string
	directory *Directory
	versions  map[uuid.UUID]*Version
	entries   map[string]interface{}
	sync.RWMutex
}

// NewFile returns a File in the referenced Directory
func (d *Directory) NewFile(name string) *File {
	return NewFile(name, d)
}

// NewFile returns a File in the referenced Directory
func NewFile(name string, directory *Directory) *File {
	f := &File{
		uuid:      uuid.NewV4(),
		name:      name,
		directory: directory,
		versions:  make(map[uuid.UUID]*Version),
	}
	err := directory.addFile(f)
	if err != nil {
		panic(err)
	}
	return f
}

func (f *File) addVersion(v *Version) error {
	if f != nil {
		f.Lock()
		defer f.Unlock()
		if f.versions[v.uuid] == nil {
			f.versions[v.uuid] = v
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return errors.New("ERROR") // FIXME: error should be made useful
}

func (f *File) removeVersion(v *Version) error {
	if f != nil {
		f.Lock()
		defer f.Unlock()
		if f.versions[v.uuid] == v {
			delete(f.versions, v.uuid)
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return errors.New("ERROR") // FIXME: error should be made useful
}

// Delete removes this File
func (f *File) Delete() error {
	if f != nil {
		for _, v := range f.versions {
			v.Delete()
		}
		if err := f.directory.removeFile(f); err != nil {
			return err
		}
	}
	return nil
}

// Directory returns parent Directory or nil
func (f *File) Directory() *Directory {
	if f == nil {
		return nil
	}
	return f.directory
}

// Versions returns all child Version entries
func (f *File) Versions() []*Version {
	if f == nil {
		return nil
	}
	versions := make([]*Version, len(f.versions))
	i := 0
	for _, version := range f.versions {
		versions[i] = version
		i++
	}
	return versions
}

// UUID returns UUID
func (f *File) UUID() uuid.UUID {
	return f.uuid
}

// Name returns name
func (f *File) Name() string {
	return f.name
}
