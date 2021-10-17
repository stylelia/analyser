package analyser

// Interface for all exec.Command stuff
type CommandRunner interface {
	Run() error
	Output() ([]byte, error)
}
