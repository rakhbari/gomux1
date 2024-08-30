package utils

import "testing"

func TestRomanToInt(t *testing.T) {
	type args struct {
		roman string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "III",
			args:    args{roman: "III"},
			want:    3,
			wantErr: false,
		},
		{
			name:    "LVIII",
			args:    args{roman: "LVIII"},
			want:    58,
			wantErr: false,
		},
		{
			name:    "MCMXCIV",
			args:    args{roman: "MCMXCIV"},
			want:    1994,
			wantErr: false,
		},
		{
			name:    "More than 15 char numeral MCMXCIVMCMXCIVIII",
			args:    args{roman: "MCMXCIVMCMXCIVIII"},
			want:    -1,
			wantErr: true,
		},
		{
			name:    "Invalid roman numeral character",
			args:    args{roman: "IIaI"},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RomanToInt(tt.args.roman)
			if (err != nil) != tt.wantErr {
				t.Errorf("RomanToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RomanToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
