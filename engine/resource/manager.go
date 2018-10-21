package resource

import (
	"strings"
)

// Very generic file storage.
// If the struct came from a file, it should be obtainable from here
type manager struct {
	resources map[string]IFile
}

// Add a new file
func (m *manager) Add(file IFile) {
	m.resources[strings.ToLower(file.GetFilePath())] = file
}

// Remove an open file
func (m *manager) Remove(filePath string) {
	delete(m.resources, strings.ToLower(filePath))
}

// Find a specific file
func (m *manager) Get(filePath string) IFile {
	return m.resources[strings.ToLower(filePath)]
}

func (m *manager) Has(filePath string) bool {
	return (m.resources[strings.ToLower(filePath)] != nil)
}

var resourceManager manager

func Manager() *manager {
	if resourceManager.resources == nil {
		resourceManager.resources = map[string]IFile{}
	}

	return &resourceManager
}
