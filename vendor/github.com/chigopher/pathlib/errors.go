package pathlib

import "fmt"

var (
	// ErrDoesNotImplement indicates that the afero filesystem doesn't
	// implement the required interface.
	ErrDoesNotImplement = fmt.Errorf("doesn't implement required interface")
	// ErrInfoIsNil indicates that a nil os.FileInfo object was provided
	ErrInfoIsNil = fmt.Errorf("provided os.Info object was nil")
	// ErrInvalidAlgorithm specifies that an unknown algorithm was given for Walk
	ErrInvalidAlgorithm = fmt.Errorf("invalid algorithm specified")
	// ErrLstatNotPossible specifies that the filesystem does not support lstat-ing
	ErrLstatNotPossible = fmt.Errorf("lstat is not possible")
	// ErrRelativeTo indicates that we could not make one path relative to another
	ErrRelativeTo = fmt.Errorf("failed to make path relative to other")
	// ErrStopWalk indicates to the Walk function that the walk should be aborted
	ErrStopWalk = fmt.Errorf("stop filesystem walk")
)
