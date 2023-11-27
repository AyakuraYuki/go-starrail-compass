package ng

import (
	"context"

	"github.com/go-logr/logr"
)

// Solver 引航罗盘求解器
type Solver interface {
	// Solve 求解引航罗盘
	Solve(ctx context.Context, compass Compass) (Steps, error)
}

// SolverOptions 求解器的选项
type SolverOptions struct {
	Logger logr.Logger
}
