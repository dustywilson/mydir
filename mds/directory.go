package mds

import (
	"errors"
	"sync"

	"github.com/satori/go.uuid"
)

// Directory is a place to store Directory entries and File entries
type Directory struct {
	uuid              uuid.UUID
	name              string
	parent            *Directory
	files             map[uuid.UUID]*File
	directories       map[uuid.UUID]*Directory
	entries           map[string]interface{}
	masterFiles       map[uuid.UUID]*File
	masterDirectories map[uuid.UUID]*Directory
	masterVersions    map[uuid.UUID]*Version
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
	if d.parent == nil {
		d.masterFiles = make(map[uuid.UUID]*File)
		d.masterDirectories = make(map[uuid.UUID]*Directory)
		d.masterVersions = make(map[uuid.UUID]*Version)
		d.masterDirectories[d.uuid] = d
	} else {
		err := parent.addChildDirectory(d)
		if err != nil {
			panic(err)
		}
	}
	return d
}

func (d *Directory) addChildDirectory(child *Directory) error {
	if d != nil {
		rootDir := d.Root()
		if !d.IsRoot() {
			rootDir.Lock()
			defer rootDir.Unlock()
		}
		d.Lock()
		defer d.Unlock()
		if d.entries[child.name] == nil {
			d.directories[child.uuid] = child
			d.entries[child.name] = child
			rootDir.masterDirectories[child.uuid] = child
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return nil
}

func (d *Directory) removeChildDirectory(child *Directory) error {
	if d != nil {
		rootDir := d.Root()
		if !d.IsRoot() {
			rootDir.Lock()
			defer rootDir.Unlock()
		}
		d.Lock()
		defer d.Unlock()
		if d.entries[child.name] == child {
			delete(d.directories, child.uuid)
			delete(d.entries, child.name)
			delete(rootDir.masterDirectories, child.uuid)
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
		rootDir := d.Root()
		if !d.IsRoot() {
			rootDir.Lock()
			defer rootDir.Unlock()
		}
		d.Lock()
		defer d.Unlock()
		if d.entries[f.name] == nil {
			d.files[f.uuid] = f
			d.entries[f.name] = f
			rootDir.masterFiles[f.uuid] = f
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return nil
}

func (d *Directory) removeFile(f *File) error {
	if d != nil {
		rootDir := d.Root()
		if !d.IsRoot() {
			rootDir.Lock()
			defer rootDir.Unlock()
		}
		d.Lock()
		defer d.Unlock()
		if d.entries[f.name] == f {
			delete(d.files, f.uuid)
			delete(d.entries, f.name)
			delete(rootDir.masterFiles, f.uuid)
			return nil
		}
		return errors.New("ERROR") // FIXME: error should be made useful
	}
	return nil
}

func (d *Directory) addFileVersion(v *Version) error {
	if d != nil {
		rootDir := d.Root()
		rootDir.Lock()
		defer rootDir.Unlock()
		rootDir.masterVersions[v.uuid] = v
	}
	return nil
}

func (d *Directory) removeFileVersion(v *Version) error {
	if d != nil {
		rootDir := d.Root()
		rootDir.Lock()
		defer rootDir.Unlock()
		delete(rootDir.masterVersions, v.uuid)
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

// IsRoot returns bool if this Directory is root
func (d *Directory) IsRoot() bool {
	return d.parent == nil
}

// Root returns the root directory
func (d *Directory) Root() *Directory {
	rootDir := d
	for rootDir.parent != nil {
		rootDir = rootDir.parent
	}
	return rootDir
}

// GetDirectoryByUUID returns a Directory by UUID
// This gets it from the Directory's Root()
func (d *Directory) GetDirectoryByUUID(u uuid.UUID) *Directory {
	d.RLock()
	defer d.RUnlock()
	return d.Root().masterDirectories[u]
}

// GetFileByUUID returns a File by UUID
// This gets it from the Directory's Root()
func (d *Directory) GetFileByUUID(u uuid.UUID) *File {
	d.RLock()
	defer d.RUnlock()
	return d.Root().masterFiles[u]
}

// GetVersionByUUID returns a File by UUID
// This gets it from the Directory's Root()
func (d *Directory) GetVersionByUUID(u uuid.UUID) *Version {
	d.RLock()
	defer d.RUnlock()
	return d.Root().masterVersions[u]
}
