package ng

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"sort"
)

// NewHungerSolver 创建穷举求解器
func NewHungerSolver(opts SolverOptions) (Solver, error) {
	return &hungerSolver{logger: opts.Logger}, nil
}

// hungerSolver 穷举求解器的实现
type hungerSolver struct {
	logger logr.Logger
}

var _ Solver = &hungerSolver{}

// Solve 求解引航罗盘
func (s *hungerSolver) Solve(ctx context.Context, compass Compass) (Steps, error) {
	if err := compass.Validate(); err != nil {
		return nil, fmt.Errorf(`invalid compass, error: %w`, err)
	}

	// 对所有可能的解法试错
	for _, solution := range s.getPossibleSolutions(compass) {
		if ok, _ := CheckSolution(compass, solution); ok {
			return solution.Standardize(), nil
		}
		s.logger.V(1).Info(fmt.Sprintf(`try solution "%s" failed`, solution.String()))
	}

	return nil, errors.New(`the compass has no solution`)
}

// getPossibleSolutions 获取所有可能的解法
func (s *hungerSolver) getPossibleSolutions(compass Compass) []Steps {
	possibleSolutions := make([]Steps, 0)

	for _, rg := range compass.RingGroups {
		temp := make([]Steps, 0)
		// 转动 6 次保证任何方案都可以回归原点
		for i := 0; i < SCALES; i++ {
			if len(possibleSolutions) == 0 {
				temp = append(temp, Steps{{RingGroup: rg, Count: i}})
				continue
			}
			for _, cur := range possibleSolutions {
				temp = append(temp, append(cur, Step{RingGroup: rg, Count: i}))
			}
		}
		possibleSolutions = temp
	}

	sort.SliceStable(possibleSolutions, func(i, j int) bool {
		sumI := 0
		for _, step := range possibleSolutions[i] {
			sumI += step.Count
		}
		sumJ := 0
		for _, step := range possibleSolutions[j] {
			sumJ += step.Count
		}
		return sumI < sumJ
	})

	return possibleSolutions
}
