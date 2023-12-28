package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	hourInMilliseconds     = 3600000
	minuteInMilliseconds   = 60000
	secondInMilliseconds   = 1000
	gametickInMilliseconds = 600
)

func TimeToMs(time string) int {
	tc := strings.Split(time, ":")
	nc := len(tc)

	ms := 0
	if nc == 3 {
		h, _ := strconv.Atoi(tc[0])
		m, _ := strconv.Atoi(tc[1])
		s, _ := strconv.Atoi(tc[2])
		ms = h*hourInMilliseconds + m*minuteInMilliseconds + s*secondInMilliseconds
	} else if nc == 2 {
		m, _ := strconv.Atoi(tc[0])
		s_and_ms := strings.Split(tc[1], ".")
		s, _ := strconv.Atoi(s_and_ms[0])
		ms_str := "0"
		if len(s_and_ms) > 1 {
			ms_str = s_and_ms[1]
			ms_str = ms_str + strings.Repeat("0", 3-len(ms_str))
		}
		ms_part, _ := strconv.Atoi(ms_str)
		ms = m*minuteInMilliseconds + s*secondInMilliseconds + ms_part
	} else if nc == 1 {
		ms_str := tc[0]
		ms, _ = strconv.Atoi(ms_str)
	}

	return ms
}

func TimeToTicks(time string) int {
	totalms := TimeToMs(time)
	ticks := totalms / gametickInMilliseconds
	if totalms%gametickInMilliseconds != 0 {
		ticks++
	}
	return ticks
}

func TicksToTime(ticks int) string {
	ms := ticks * gametickInMilliseconds
	h := ms / hourInMilliseconds
	ms %= hourInMilliseconds
	m := ms / minuteInMilliseconds
	ms %= minuteInMilliseconds
	s := ms / secondInMilliseconds
	ms %= secondInMilliseconds
	ms_str := fmt.Sprintf("%02d", ms/10)

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%s", h, m, s, ms_str)
	} else if m > 0 {
		return fmt.Sprintf("%02d:%02d.%s", m, s, ms_str)
	} else {
		return fmt.Sprintf("%02d.%s", s, ms_str)
	}
}
