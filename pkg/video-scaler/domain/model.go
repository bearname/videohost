package domain

import "strconv"

type Quality int

const (
	Q2160p Quality = iota
	Q1440p
	Q1080p
	Q720p
	Q480p
	Q320p
)

func (q Quality) String() string {
	return [...]string{"2160", "1440", "1080", "720", "480", "360"}[q]
}

func (q Quality) Values() int {
	return [...]int{2160, 1440, 1080, 720, 480, 360}[q]
}

func IsSupportedQuality(quality int) bool {
	ints := []int{2160, 1440, 1080, 720, 480, 360}
	for _, a := range ints {
		if a == quality {
			return true
		}
	}

	return false
}

type Resolution struct {
	width  int
	height int
}

func (r *Resolution) String() string {
	return strconv.Itoa(r.width) + `x` + strconv.Itoa(r.height)
}

func QualityToResolution(quality Quality) Resolution {
	m := map[Quality]Resolution{
		Q1080p: {width: 1920, height: 1080},
		Q720p:  {width: 1280, height: 720},
		Q480p:  {width: 850, height: 480},
		Q320p:  {width: 640, height: 360},
	}
	return m[quality]
}
