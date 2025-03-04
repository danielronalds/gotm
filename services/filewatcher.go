package services

type FilewatcherService struct {
	filesystem FilesystemReader
	// Cache for keeping track of previous reads of a file
	cache map[string]string
}

func NewFilewatcherService(filesystem FilesystemReader) FilewatcherService {
	cache := make(map[string]string)
	return FilewatcherService{filesystem, cache}
}

func (s FilewatcherService) UpdateCache(directory string) error {
	// Cleaning cache
	for k := range s.cache {
		delete(s.cache, k)
	}

	projectFiles, err := s.filesystem.ReadDirRecursive(directory)
	if err != nil {
		return err
	}

	for _, file := range projectFiles {
		contents, err := s.filesystem.ReadFile(file)
		if err != nil {
			return err
		}

		s.cache[file] = contents
	}

	return nil
}

func (s FilewatcherService) HaveFilesChanged(directory string) (bool, error) {
	projectFiles, err := s.filesystem.ReadDirRecursive(directory)
	if err != nil {
		return false, err
	}

	projectChanged := false

	for _, file := range projectFiles {
		changed, err := s.hasFileChanged(file)
		if err != nil {
			return projectChanged, err
		}

		if changed {
			projectChanged = true
		}
	}

	return projectChanged, nil
}

func (s FilewatcherService) hasFileChanged(filename string) (bool, error) {
	contents, err := s.filesystem.ReadFile(filename)
	if err != nil {
		return false, err
	}

	oldContents, ok := s.cache[filename]
	s.cache[filename] = contents

	if !ok {
		return true, nil
	}

	return contents != oldContents, nil
}
