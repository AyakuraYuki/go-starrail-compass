package ng

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	compassRegexpStr = `(?P<outerRing>[0-9-+]+),` +
		`(?P<middleRing>[0-9-+]+),` +
		`(?P<innerRing>[0-9-+]+)/` +
		`(?P<ringGroups>[imo,]+)`
	ringRegexpStr = `(?P<location>[0-5])(?P<speed>(?:\+|-)[1-4])`
)

var (
	compassRegexp = regexp.MustCompile(compassRegexpStr)
	ringRegexp    = regexp.MustCompile(ringRegexpStr)
)

// ParseCompass 解析罗盘信息表达式
// 罗盘信息表达式满足如下格式:
//
//	{oLoc}{oSpeed},{mLoc}{mSpeed},{iLoc}{iSpeed}/{rg1},{rg2},{rg3}
//	{oLoc}, {mLoc}, {iLoc}：分别表示外圈、中圈、内圈的初始刻度
//	{oSpeed}, {mSpeed}, {iSpeed}：分别表示外圈、中圈、内圈的旋转速度
//	{rg1}, {rg2}, {rg3}：分别表示三种转动方案，可选值如下
//		o
//		m
//		i
//		om 或 mo
//		oi 或 io
//		im 或 mi
func ParseCompass(expression string) (Compass, error) {
	compass := Compass{}

	// 正则解析获得捕获组
	groups := compassRegexp.FindStringSubmatch(expression)
	if len(groups) == 0 {
		return compass, fmt.Errorf(`invalid compass expression: "%s" (not match "%s")`, expression, compassRegexpStr)
	}

	// 解析各捕获组的表达式
	outer, err := parseRing(groups[compassRegexp.SubexpIndex("outerRing")])
	if err != nil {
		return compass, fmt.Errorf(`parse outer ring error: %w`, err)
	}
	compass.OuterRing = outer

	middle, err := parseRing(groups[compassRegexp.SubexpIndex("middleRing")])
	if err != nil {
		return compass, fmt.Errorf(`parse middle ring error: %w`, err)
	}
	compass.MiddleRing = middle

	inner, err := parseRing(groups[compassRegexp.SubexpIndex("innerRing")])
	if err != nil {
		return compass, fmt.Errorf(`parse inner ring error: %w`, err)
	}
	compass.InnerRing = inner

	ringGroups, err := parseRingGroups(groups[compassRegexp.SubexpIndex("ringGroups")])
	if err != nil {
		return compass, fmt.Errorf("parse ring groups error: %w", err)
	}
	compass.RingGroups = ringGroups

	return compass, nil
}

// parseRing 解析罗盘圈表达式
func parseRing(expression string) (Ring, error) {
	ring := Ring{}

	// 正则解析
	groups := ringRegexp.FindStringSubmatch(expression)
	if len(groups) == 0 {
		return ring, fmt.Errorf(`invalid ring expression: "%s" (not match "%s")`, expression, ringRegexpStr)
	}

	locationPart := groups[ringRegexp.SubexpIndex("location")]
	location, err := strconv.ParseInt(locationPart, 10, 8)
	if err != nil {
		return ring, fmt.Errorf(`parse ring location "%s" error: %w`, locationPart, err)
	}
	ring.Location = int(location)

	speedPart := groups[ringRegexp.SubexpIndex("speed")]
	speed, err := strconv.ParseInt(speedPart, 10, 8)
	if err != nil {
		return ring, fmt.Errorf(`parse ring speed "%s" error: %w`, speedPart, err)
	}
	ring.Speed = int(speed)

	return ring, nil
}

// parseRingGroups 解析罗盘转动方案表达式
func parseRingGroups(expression string) ([]RingGroup, error) {
	ringGroups := make([]RingGroup, 0)
	parts := strings.Split(expression, ",")
	for i, expr := range parts {
		ringGroup, err := parseRingGroup(expr)
		if err != nil {
			return nil, fmt.Errorf(`parse ring groups at index %d error: %w`, i, err)
		}
		ringGroups = append(ringGroups, ringGroup)
	}
	return ringGroups, nil
}

// parseRingGroup 解析罗盘圈转动方案表达式
func parseRingGroup(expression string) (RingGroup, error) {
	switch expression {
	case "o":
		return Outer, nil
	case "m":
		return Middle, nil
	case "i":
		return Inner, nil
	case "om", "mo":
		return OuterMiddle, nil
	case "oi", "io":
		return OuterInner, nil
	case "mi", "im":
		return MiddleInner, nil
	}
	return 0, fmt.Errorf(`unknown ring group: %s`, expression)
}
