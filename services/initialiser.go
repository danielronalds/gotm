package services

// Service for hanlding initialising new projects
type InitialiserService struct{}

func NewInitialiserService() InitialiserService {
	return InitialiserService{}
}

func (s InitialiserService) CreateProject(directoryName string) error {
	// TODO: Write this functionality
	return nil
}
