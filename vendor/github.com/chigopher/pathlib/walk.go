package pathlib

import (
	"errors"
	"fmt"
	"os"
)

// WalkOpts is the struct that defines how a walk should be performed
type WalkOpts struct {
	// Depth defines how far down a directory we should recurse. A value of -1 means
	// infinite depth. 0 means only the direct children of root will be returned, etc.
	Depth int

	// Algorithm specifies the algoritm that the Walk() function should use to
	// traverse the directory.
	Algorithm Algorithm

	// FollowSymlinks defines whether symlinks should be dereferenced or not. If True,
	// the symlink itself will never be returned to WalkFunc, but rather whatever it
	// points to. Warning!!! You are exposing yourself to substantial risk by setting this
	// to True. Here be dragons!
	FollowSymlinks bool

	// MinimumFileSize specifies the minimum size of a file for visitation.
	// If negative, there is no minimum size.
	MinimumFileSize int64

	// MaximumFileSize specifies the maximum size of a file for visitation.
	// If negative, there is no maximum size.
	MaximumFileSize int64

	// VisitFiles specifies that we should visit regular files during
	// the walk.
	VisitFiles bool

	// VisitDirs specifies that we should visit directories during the walk.
	VisitDirs bool

	// VisitSymlinks specifies that we should visit symlinks during the walk.
	VisitSymlinks bool

	// VisitFirst specifies that, in the algorithms where it is appropriate,
	// a node's contents should be visited first, before recursing down. If false,
	// a node's subdirectories will be recursed first before visiting any of its
	// other children.
	//
	// This option is not appropriate in the Basic algorithm, where ordering is
	// explicitly forbidden.
	// VisitFirst bool
}

// DefaultWalkOpts returns the default WalkOpts struct used when
// walking a directory.
func DefaultWalkOpts() *WalkOpts {
	return &WalkOpts{
		Depth:           -1,
		Algorithm:       AlgorithmBasic,
		FollowSymlinks:  false,
		MinimumFileSize: -1,
		MaximumFileSize: -1,
		VisitFiles:      true,
		VisitDirs:       true,
		VisitSymlinks:   true,
	}
}

// MeetsMinimumSize returns whether size is at least the minimum specified.
func (w *WalkOpts) MeetsMinimumSize(size int64) bool {
	if w.MinimumFileSize < 0 {
		return true
	}
	return size >= w.MinimumFileSize
}

// MeetsMaximumSize returns whether size is less than or equal to the maximum specified.
func (w *WalkOpts) MeetsMaximumSize(size int64) bool {
	if w.MaximumFileSize < 0 {
		return true
	}
	return size <= w.MaximumFileSize
}

// Algorithm represents the walk algorithm that will be performed.
type Algorithm int

const (
	// AlgorithmBasic is a walk algorithm. It iterates over filesystem objects in the
	// order in which they are returned by the operating system. It guarantees no
	// ordering of any kind. This is the most efficient algorithm and should be used
	// in all cases where ordering does not matter.
	AlgorithmBasic Algorithm = iota
	// AlgorithmDepthFirst is a walk algorithm. It iterates over a filesystem tree
	// by first recursing as far down as it can in one path. Each directory is visited
	// only after all of its children directories have been recursed.
	AlgorithmDepthFirst
)

// Walk is an object that handles walking through a directory tree
type Walk struct {
	Opts *WalkOpts
	root *Path
}

type WalkOptsFunc func(config *WalkOpts)

func WalkDepth(depth int) WalkOptsFunc {
	return func(config *WalkOpts) {
		config.Depth = depth
	}
}

func WalkAlgorithm(algo Algorithm) WalkOptsFunc {
	return func(config *WalkOpts) {
		config.Algorithm = algo
	}
}

func WalkFollowSymlinks(follow bool) WalkOptsFunc {
	return func(config *WalkOpts) {
		config.FollowSymlinks = follow
	}
}

func WalkMinimumFileSize(size int64) WalkOptsFunc {
	return func(config *WalkOpts) {
		config.MinimumFileSize = size
	}
}

func WalkMaximumFileSize(size int64) WalkOptsFunc {
	return func(config *WalkOpts) {
		config.MaximumFileSize = size
	}
}

func WalkVisitFiles(value bool) WalkOptsFunc {
	return func(config *WalkOpts) {
		config.VisitFiles = value
	}
}

func WalkVisitDirs(value bool) WalkOptsFunc {
	return func(config *WalkOpts) {
		config.VisitDirs = value
	}
}

func WalkVisitSymlinks(value bool) WalkOptsFunc {
	return func(config *WalkOpts) {
		config.VisitSymlinks = value
	}
}

// NewWalk returns a new Walk struct with default values applied
func NewWalk(root *Path, opts ...WalkOptsFunc) (*Walk, error) {
	config := DefaultWalkOpts()
	for _, opt := range opts {
		opt(config)
	}
	return NewWalkWithOpts(root, config)
}

// NewWalkWithOpts returns a Walk object with the given WalkOpts applied
func NewWalkWithOpts(root *Path, opts *WalkOpts) (*Walk, error) {
	if root == nil {
		return nil, fmt.Errorf("root path can't be nil")
	}
	if opts == nil {
		return nil, fmt.Errorf("opts can't be nil")
	}
	return &Walk{
		Opts: opts,
		root: root,
	}, nil
}

func (w *Walk) maxDepthReached(currentDepth int) bool {
	if w.Opts.Depth >= 0 && currentDepth > w.Opts.Depth {
		return true
	}
	return false
}

type dfsObjectInfo struct {
	path *Path
	info os.FileInfo
	err  error
}

func (w *Walk) walkDFS(walkFn WalkFunc, root *Path, currentDepth int) error {
	if w.maxDepthReached(currentDepth) {
		return nil
	}

	var children []*dfsObjectInfo

	if err := w.iterateImmediateChildren(root, func(child *Path, info os.FileInfo, encounteredErr error) error {
		// Since we are doing depth-first, we have to first recurse through all the directories,
		// and save all non-directory objects so we can defer handling at a later time.
		if IsDir(info.Mode()) {
			if err := w.walkDFS(walkFn, child, currentDepth+1); err != nil {
				return err
			}
		}

		children = append(children, &dfsObjectInfo{
			path: child,
			info: info,
			err:  encounteredErr,
		})

		return nil
	}); err != nil {
		return err
	}

	// Iterate over all children after all subdirs have been recursed
	for _, child := range children {
		passesQuery, err := w.passesQuerySpecification(child.info)
		if err != nil {
			return err
		}

		if passesQuery {
			if err := walkFn(child.path, child.info, child.err); err != nil {
				return err
			}
		}

	}
	return nil
}

// iterateImmediateChildren is a function that handles discovering root's immediate children,
// and will run the algorithm function for every child. The algorithm function is essentially
// what differentiates how each walk behaves, and determines what actions to take given a
// certain child.
func (w *Walk) iterateImmediateChildren(root *Path, algorithmFunction WalkFunc) error {
	children, err := root.ReadDir()
	if err != nil {
		return err
	}

	var info os.FileInfo
	for _, child := range children {
		if child.String() == root.String() {
			continue
		}
		if w.Opts.FollowSymlinks {
			info, err = child.Stat()
			if err != nil {
				return err
			}
		} else {
			info, err = child.Lstat()
		}

		if info == nil {
			if err != nil {
				return err
			}
			return ErrInfoIsNil
		}

		if algoErr := algorithmFunction(child, info, err); algoErr != nil {
			return algoErr
		}
	}
	return nil
}

// passesQuerySpecification returns whether or not the object described by
// the os.FileInfo passes all of the query specifications listed in
// the walk options.
func (w *Walk) passesQuerySpecification(info os.FileInfo) (bool, error) {
	if IsFile(info.Mode()) {
		if !w.Opts.VisitFiles {
			return false, nil
		}

		if !w.Opts.MeetsMinimumSize(info.Size()) ||
			!w.Opts.MeetsMaximumSize(info.Size()) {
			return false, nil
		}
	} else if IsDir(info.Mode()) && !w.Opts.VisitDirs {
		return false, nil
	} else if IsSymlink(info.Mode()) && !w.Opts.VisitSymlinks {
		return false, nil
	}

	return true, nil
}

func (w *Walk) walkBasic(walkFn WalkFunc, root *Path, currentDepth int) error {
	if w.maxDepthReached(currentDepth) {
		return nil
	}

	err := w.iterateImmediateChildren(root, func(child *Path, info os.FileInfo, encounteredErr error) error {
		if IsDir(info.Mode()) {
			if err := w.walkBasic(walkFn, child, currentDepth+1); err != nil {
				return err
			}
		}

		passesQuery, err := w.passesQuerySpecification(info)
		if err != nil {
			return err
		}

		if passesQuery {
			if err := walkFn(child, info, encounteredErr); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// WalkFunc is the function provided to the Walk function for each directory.
type WalkFunc func(path *Path, info os.FileInfo, err error) error

// Walk walks the directory using the algorithm specified in the configuration.
func (w *Walk) Walk(walkFn WalkFunc) error {

	switch w.Opts.Algorithm {
	case AlgorithmBasic:
		if err := w.walkBasic(walkFn, w.root, 0); err != nil {
			if errors.Is(err, ErrStopWalk) {
				return nil
			}
			return err
		}
		return nil
	case AlgorithmDepthFirst:
		if err := w.walkDFS(walkFn, w.root, 0); err != nil {
			if errors.Is(err, ErrStopWalk) {
				return nil
			}
			return err
		}
		return nil
	default:
		return ErrInvalidAlgorithm
	}
}
