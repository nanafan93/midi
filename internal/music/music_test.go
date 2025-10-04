package music

import "testing"

func TestIntegerToNoteName(t *testing.T) {
	type args struct {
		noteNumber int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Very Low C", args{24}, "(C)1"},
		{"Low C", args{36}, "(C)2"},
		{"Sub Contra C", args{48}, "(C)3"},
		{"Middle C", args{60}, "(C)4"},
		{"Tenor C", args{72}, "(C)5"},
		{"Soprano C", args{84}, "(C)6"},
		{"Double High C", args{96}, "(C)7"},
		{"Eighth Octave C", args{108}, "(C)8"},

		{"G Sharp 4", args{68}, "(G♯, A♭)4"},
		{"B1", args{35}, "(B)1"},
		{"F3", args{53}, "(F)3"},
		{"A5", args{81}, "(A)5"},
		{"D7", args{98}, "(D)7"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntegerToNoteName(tt.args.noteNumber); got != tt.want {
				t.Errorf("IntegerToNoteName() = %v, want %v", got, tt.want)
			}
		})
	}
}
