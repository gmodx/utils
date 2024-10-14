package versionrange

import "github.com/blang/semver/v4"

type VersionRange struct {
	Min string
	Max string
}

// Match: check version match, logic: min <= current < max
func (vr VersionRange) Match(version string) bool {
	if version == "" {
		return false
	}

	minV, err := semver.Make(vr.Min)
	if err != nil {
		return false
	}

	currV, err := semver.Make(version)
	if err != nil {
		return false
	}

	if currV.Compare(minV) < 0 { // current < min
		return false
	}

	if vr.Max == "" {
		return true
	}

	maxV, err := semver.Make(vr.Max)
	if err != nil {
		return false
	}

	if currV.Compare(maxV) >= 0 { // current >= max
		return false
	}

	return true
}
