package services

type StackService struct{}

func (s *StackService) Up(name string) error {
	// TODO: move logic from internal/api + orchestrator
	return nil
}

func (s *StackService) Down(name string) error {
	return nil
}
