package ng

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// SCALES 总刻度值
const SCALES = 6

// Ring 定义引航罗盘中的一圈
type Ring struct {
	// 位置
	// 指针从罗盘正左方∠0°沿顺时针方向旋转至当前位置所需的刻度
	// 刻度以∠60°为一度
	// 例如：0 表示正左方∠0°位置，3 表示正右方∠180°位置
	// 有效范围是 0-6
	Location int
	// 旋转速度
	// 单位为刻度，符号表示旋转方向，正数是顺时针，负数是逆时针
	// 旋转速度不会是 0（目前游戏内没有发现不旋转的圈）
	// 例如：-1 表示每次逆时针旋转 1 个刻度；2 表示每次顺时针旋转 2 个刻度
	// 有效范围是 1-4
	Speed int
}

// String 转为字符串表述
func (r *Ring) String() string {
	return fmt.Sprintf("%d%+d", r.Location, r.Speed)
}

// RingGroup 引航罗盘方案分组
type RingGroup uint8

// RingGroup 的合法值
const (
	Outer       RingGroup = 0b100
	Middle      RingGroup = 0b010
	Inner       RingGroup = 0b001
	OuterMiddle           = Outer | Middle
	OuterInner            = Outer | Inner
	MiddleInner           = Middle | Inner
)

// Name 返回方案名称
func (rg RingGroup) Name() string {
	switch rg {
	case Outer:
		return "Outer"
	case Middle:
		return "Middle"
	case Inner:
		return "Inner"
	case OuterMiddle:
		return "OuterMiddle"
	case OuterInner:
		return "OuterInner"
	case MiddleInner:
		return "MiddleInner"
	}
	return ""
}

// ShortName 返回缩写
func (rg RingGroup) ShortName() string {
	switch rg {
	case Outer:
		return "o"
	case Middle:
		return "m"
	case Inner:
		return "i"
	case OuterMiddle:
		return "om"
	case OuterInner:
		return "oi"
	case MiddleInner:
		return "mi"
	}
	return ""
}

// String 转为字符串表述
func (rg RingGroup) String() string {
	return rg.Name()
}

// Compass 引航罗盘
type Compass struct {
	OuterRing  Ring        // 外圈
	MiddleRing Ring        // 中圈
	InnerRing  Ring        // 内圈
	RingGroups []RingGroup // 方案，可以同时旋转的一个或多个圈组成的一个分组
}

// Validate 合法化
func (c *Compass) Validate() error {
	if c.OuterRing.Speed == 0 || c.MiddleRing.Speed == 0 || c.InnerRing.Speed == 0 {
		return errors.New("ring speed must be declared")
	}
	return nil
}

// IsRingGroupSupported 判断指定方案是否受当前罗盘支持
func (c *Compass) IsRingGroupSupported(ringGroup RingGroup) bool {
	if len(c.RingGroups) == 0 {
		return false
	}
	for _, v := range c.RingGroups {
		if v == ringGroup {
			return true
		}
	}
	return false
}

// Standardize 标准化罗盘
func (c *Compass) Standardize() *Compass {
	// 拷贝原始罗盘的方案，并排序
	sortedRingGroups := make([]RingGroup, len(c.RingGroups))
	copy(sortedRingGroups, c.RingGroups)
	sort.Slice(sortedRingGroups, func(i, j int) bool {
		return sortedRingGroups[i] < sortedRingGroups[j]
	})

	// 对拷贝并且排序后的方案去重
	deduplicatedRingGroups := make([]RingGroup, 0)
	for _, v := range sortedRingGroups {
		if len(deduplicatedRingGroups) > 0 && deduplicatedRingGroups[len(deduplicatedRingGroups)-1] == v {
			continue
		}
		deduplicatedRingGroups = append(deduplicatedRingGroups, v)
	}

	return &Compass{
		OuterRing: Ring{
			Location: Mod(c.OuterRing.Location, SCALES),
			Speed:    c.OuterRing.Speed % SCALES,
		},
		MiddleRing: Ring{
			Location: Mod(c.MiddleRing.Location, SCALES),
			Speed:    c.MiddleRing.Speed % SCALES,
		},
		InnerRing: Ring{
			Location: Mod(c.InnerRing.Location, SCALES),
			Speed:    c.InnerRing.Speed % SCALES,
		},
		RingGroups: deduplicatedRingGroups,
	}
}

// String 转为字符串表述
func (c *Compass) String() string {
	// 标准化罗盘
	std := c.Standardize()
	// 转换方案获取简称
	ringGroups := make([]string, len(std.RingGroups))
	for i := range ringGroups {
		ringGroups[i] = std.RingGroups[i].ShortName()
	}
	content := strings.Join(ringGroups, ",")
	// 组合罗盘信息
	return fmt.Sprintf("%s,%s,%s/%s",
		std.OuterRing.String(),
		std.MiddleRing.String(),
		std.InnerRing.String(),
		content,
	)
}
