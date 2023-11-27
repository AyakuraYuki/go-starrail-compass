package ng

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Step 声明了转动方案的操作次数
type Step struct {
	// 转动方案
	RingGroup RingGroup
	// 转动次数
	Count int
}

// Validate 合法化
func (s *Step) Validate() error {
	if s.RingGroup == 0 {
		return errors.New("invalid ring group")
	}
	return nil
}

// String 转为字符串表述
func (s *Step) String() string {
	if s.Count <= 0 {
		return ""
	}
	return fmt.Sprintf("%s%d", s.RingGroup.ShortName(), s.Count)
}

// Steps 引航罗盘解谜步骤组合
type Steps []Step

// Validate 合法化
func (s Steps) Validate() error {
	return nil
}

// Standardize 标准化
func (s Steps) Standardize() Steps {
	if len(s) == 0 {
		return nil
	}

	// 拷贝并排序
	sortedSteps := make(Steps, len(s))
	copy(sortedSteps, s)
	sort.Slice(sortedSteps, func(i, j int) bool {
		return sortedSteps[i].RingGroup < sortedSteps[j].RingGroup
	})

	// 组合
	var simplified Steps
	for _, step := range sortedSteps {
		if step.Count <= 0 {
			continue
		}
		if len(simplified) > 0 && simplified[len(simplified)-1].RingGroup == step.RingGroup {
			simplified[len(simplified)-1].Count += step.Count
		} else {
			simplified = append(simplified, step)
		}
	}
	return simplified
}

// String 转为字符串表述
func (s Steps) String() string {
	// 标准化
	std := s.Standardize()
	if len(std) == 0 {
		return ""
	}

	// 逐个转换成字符串
	steps := make([]string, len(std))
	for i := range std {
		steps[i] = std[i].String()
	}
	return strings.Join(steps, ",")
}

// CheckSolution 检查罗盘解谜步骤
// 对穷举的解法试错，尝试捕获最小能够解谜的转动操作方案
func CheckSolution(compass Compass, solution Steps) (bool, error) {
	if err := compass.Validate(); err != nil {
		return false, fmt.Errorf(`invalid compass, error: %w`, err)
	}
	if err := solution.Validate(); err != nil {
		return false, fmt.Errorf(`invalid solution, error: %w`, err)
	}

	// 各圈的初始位置
	outer := compass.OuterRing.Location
	middle := compass.MiddleRing.Location
	inner := compass.InnerRing.Location

	// 转一下
	for _, step := range solution {
		if !compass.IsRingGroupSupported(step.RingGroup) {
			return false, fmt.Errorf(`steps contains an unexpected ring group, which is not supported by compass: %s (must be one of %v)`, step.RingGroup.Name(), compass.RingGroups)
		}

		if step.RingGroup&Outer > 0 {
			outer += step.Count * compass.OuterRing.Speed
		}
		if step.RingGroup&Middle > 0 {
			middle += step.Count * compass.MiddleRing.Speed
		}
		if step.RingGroup&Inner > 0 {
			inner += step.Count * compass.InnerRing.Speed
		}
	}

	// 检查转动后的最终位置
	if inner%SCALES != 0 || middle%SCALES != 0 || outer%SCALES != 0 {
		return false, nil
	}
	return true, nil
}
