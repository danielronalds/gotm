package services

type FilewatcherServiceFilesystem interface {
	DirReader
	FileReader
}

type FilewatcherService struct {
	filesystem FilewatcherServiceFilesystem
	// Cache for keeping track of previous reads of a file
	cache map[string]string
}

func NewFilewatcherService(filesystem FilewatcherServiceFilesystem) FilewatcherService {
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

func (s FilewatcherService) HaveFilesChanged(directory string) ([]string, error) {
	filesChanged := make([]string, 0)

	projectFiles, err := s.filesystem.ReadDirRecursive(directory)
	if err != nil {
		return filesChanged, err
	}

	for _, file := range projectFiles {
		changed, err := s.hasFileChanged(file)
		if err != nil {
			return filesChanged, err
		}

		if changed {
			filesChanged = append(filesChanged, file)
		}
	}

	return filesChanged, nil
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
