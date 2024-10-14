package versionrange

import (
	"testing"
)

func TestVersionRange_Match(t *testing.T) {
	tests := []struct {
		name    string
		vr      VersionRange
		version string
		want    bool
	}{
		{
			name: "version within range",
			vr: VersionRange{
				Min: "1.0.0",
				Max: "2.0.0",
			},
			version: "1.5.0",
			want:    true,
		},
		{
			name: "version below minimum",
			vr: VersionRange{
				Min: "1.0.0",
				Max: "2.0.0",
			},
			version: "0.9.9",
			want:    false,
		},
		{
			name: "version equals minimum",
			vr: VersionRange{
				Min: "1.0.0",
				Max: "2.0.0",
			},
			version: "1.0.0",
			want:    true,
		},
		{
			name: "version equals maximum",
			vr: VersionRange{
				Min: "1.0.0",
				Max: "2.0.0",
			},
			version: "2.0.0",
			want:    false,
		},
		{
			name: "version above maximum",
			vr: VersionRange{
				Min: "1.0.0",
				Max: "2.0.0",
			},
			version: "2.1.0",
			want:    false,
		},
		{
			name: "no max version, version within range",
			vr: VersionRange{
				Min: "1.0.0",
				Max: "",
			},
			version: "3.0.0",
			want:    true,
		},
		{
			name: "invalid version format",
			vr: VersionRange{
				Min: "1.0.0",
				Max: "2.0.0",
			},
			version: "invalid.version",
			want:    false,
		},
		{
			name: "invalid min format",
			vr: VersionRange{
				Min: "invalid.version",
				Max: "2.0.0",
			},
			version: "1.5.0",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vr.Match(tt.version); got != tt.want {
				t.Errorf("VersionRange.Match() = %v, want %v, case name: %v", got, tt.want, tt.name)
			}
		})
	}
}
