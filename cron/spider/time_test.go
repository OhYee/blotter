package spider

import (
	"testing"
)

func Test_parseTime(t *testing.T) {
	tests := []struct {
		s    string
		want int64
	}{
		{
			s:    "2021-09-05",
			want: 1630771200,
		},
		{
			s:    "2021-9-5",
			want: 1630771200,
		},
		{
			s:    "2021-09-5",
			want: 1630771200,
		},
		{
			s:    "2021-9-05",
			want: 1630771200,
		},
		{
			s:    "2021/09/05",
			want: 1630771200,
		},
		{
			s:    "2021/9/5",
			want: 1630771200,
		},
		{
			s:    "2021/09/5",
			want: 1630771200,
		},
		{
			s:    "2021/9/05",
			want: 1630771200,
		},
		{
			s:    "2021-09-05 17:30:25",
			want: 1630834225,
		},
		{
			s:    "21-09-05 17:30:25",
			want: 1630834225,
		},
		{
			s:    "2021-09-05 25:30:25",
			want: 1630771200, // 只处理日期部分
		},
		{
			s:    "2021-13-05 25:30:25",
			want: 0,
		},
		{
			s:    "21-09-05 17:30:25",
			want: 1630834225,
		},
		{
			s:    "time is: 21-09-05",
			want: 1630771200,
		},
		{
			s:    "0000-00-00",
			want: 0,
		},
		{
			s:    "21|09|05",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := parseTime(tt.s)
			ts := int64(0)
			if got != nil {
				ts = got.Unix()
			}
			if tt.want != ts {
				t.Errorf("parseTime() = %v, want %v", ts, tt.want)
			}
		})
	}
}
