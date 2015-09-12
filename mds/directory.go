package mds

import (
	"errors"
	"sync"

	"github.com/satori/go.uuid"
)

// Directory is a place to store Directory entries and File entries
type Directory struct {
	uuid        uuid.UUID
	name        string
	parent      *Directory
	files       map[uuid.UUID]*File
	directories map[uuid.UUID]*Directory
	entries     map[string]interface{}
	sync.RWMutex
}

// NewDirectory returns a Directory
func (d *Directory) NewDirectory(name string) *Directory {
	return NewDirectory(name, d)
}

// NewDirectory returns a Directory
// parent can be nil, but only if it's a root directory.
func NewDirectory(name string, parent *Directory) *Directory {
	d := &Directory{
		uuid:        uuid.NewV4(),
		name:        name,
		parent:      parent,
		files:       make(map[uuid.UUID]*File),
		directories: make(map[uuid.UUID]*Directory),
		entries:     make(map[string]interface{}),
	}
	err := parent.addChildDirectory(d)
	if err != nil {
		panic(err)
	}
	return d
}

func (d *Directory) addChildDirectory(child *Directory) error {
	if d != nil {
		d.Lock()
		defer d.Unlock()
		if d.entries[child.name] == nil {
			d.directories[child.uuid] = child
			d.entries[child.name] = child
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return nil
}

func (d *Directory) removeChildDirectory(child *Directory) error {
	if d != nil {
		d.Lock()
		defer d.Unlock()
		if d.entries[child.name] == child {
			delete(d.directories, child.uuid)
			delete(d.entries, child.name)
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return nil
}

func (d *Directory) renameChildDirectory(child *Directory, oldname string, newname string) error {
	if d != nil {
		d.Lock()
		defer d.Unlock()
		if d.entries[oldname] == child {
			delete(d.entries, oldname)
			d.entries[newname] = child
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return nil
}

func (d *Directory) addFile(f *File) error {
	if d != nil {
		d.Lock()
		defer d.Unlock()
		if d.entries[f.name] == nil {
			d.files[f.uuid] = f
			d.entries[f.name] = f
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return nil
}

func (d *Directory) removeFile(f *File) error {
	if d != nil {
		d.Lock()
		defer d.Unlock()
		if d.entries[f.name] == f {
			delete(d.files, f.uuid)
			delete(d.entries, f.name)
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return nil
}

func (d *Directory) renameFile(f *File, oldname string, newname string) error {
	if d != nil {
		d.Lock()
		defer d.Unlock()
		if d.entries[oldname] == f {
			delete(d.entries, oldname)
			d.entries[newname] = f
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return nil
}

// Delete removes this Directory
func (d *Directory) Delete(deleteChildren bool) bool {
	if d != nil {
		if deleteChildren {
			for _, child := range d.directories {
				if !child.Delete(deleteChildren) {
					return false
				}
			}
		} else if len(d.directories)+len(d.files) > 0 {
			return false
		}
		d.parent.removeChildDirectory(d)
		return true
	}
	return false
}

// Parent returns parent Directory or nil
func (d *Directory) Parent() *Directory {
	if d == nil {
		return nil
	}
	return d.parent
}

// Directories returns all child Directory entries
func (d *Directory) Directories() []*Directory {
	if d == nil {
		return nil
	}
	dirs := make([]*Directory, len(d.directories))
	i := 0
	for _, dir := range d.directories {
		dirs[i] = dir
		i++
	}
	return dirs
}

// Files returns all child File entries
func (d *Directory) Files() []*File {
	if d == nil {
		return nil
	}
	files := make([]*File, len(d.files))
	i := 0
	for _, file := range d.files {
		files[i] = file
		i++
	}
	return files
}

// UUID returns UUID
func (d *Directory) UUID() uuid.UUID {
	return d.uuid
}

// Name returns name
func (d *Directory) Name() string {
	return d.name
}

// GetByName return entry by name; it'll be either a Directory or a File
func (d *Directory) GetByName(name string) (*Directory, *File, error) {
	if d != nil {
		d.Lock()
		defer d.Unlock()
		switch entry := d.entries[name].(type) {
		case nil:
			return nil, nil, errors.New("404 " + name)
		case *Directory:
			return entry, nil, nil
		case *File:
			return nil, entry, nil
		default:
			return nil, nil, errors.New("ERROR A")
		}
	}
	return nil, nil, errors.New("ERROR B")
}
