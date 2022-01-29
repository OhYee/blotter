package avatar

import "testing"

func TestGet(t *testing.T) {
	tests := []struct {
		email      string
		wantAvatar string
	}{
		{
			email:      "oyohyee@oyohyee.com",
			wantAvatar: "https://avatars.githubusercontent.com/u/13498329?v=4",
		},
		{
			email:      "896817156@qq.com",
			wantAvatar: "https://cravatar.cn/avatar/c92b45d4385ba3d97a777e63bb519908?s=640&d=404",
		},
		{
			email:      "404@oyohyee.com",
			wantAvatar: DefaultAvatar,
		},
	}
	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if gotAvatar := Get(tt.email); gotAvatar != tt.wantAvatar {
				t.Errorf("Get() = %v, want %v", gotAvatar, tt.wantAvatar)
			}
		})
	}
}
