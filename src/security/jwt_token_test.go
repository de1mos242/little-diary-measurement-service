package security

import (
	"little-diary-measurement-service/src/config"
	_ "little-diary-measurement-service/src/test_data"
	"testing"
)

func TestJwtTokenReader_ReadUserUuid(t *testing.T) {
	type fields struct {
		PublicKey string
	}
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test valid token",
			fields:  fields{PublicKey: config.Config.GetAuthServerJwtPublicKey()},
			args:    args{tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0ODk5NTgsIm5iZiI6MTU4NTQ4OTk1OCwianRpIjoiNWE4OGNlYTctMzQ5NC00MThiLTlmMWYtMGQ5NjI0NDkwNjU1IiwiZXhwIjoyNTg1NDkwODU4LCJpZGVudGl0eSI6MSwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6ImFkbWluIiwidXVpZCI6IjNhZWY5YjVjLTg4M2YtNDEyZC1hZGU3LTQ3YmU0MzgyN2Q2OCIsInJlc291cmNlcyI6W119fQ.m_tC6rROZWPU1quds3lf07r8toOX5Gwg-FRKL8UIo6We_KT5dVYeU23330619BLaaNgSOObpTNysl31hxLBBXgyxngnyqRroXXWIz0lcieY76X18Z7xfq8twKavyaAjN4XyUmTVtTzBuN1VKknBIzGGN9QpKfJrjpFWGFVNXj3bOmJRqHyDO49uGQ5CgJdidnQ_OFVRXugXuLnkP94fe4bEAQxW-vE53zkrgdNFqElCvj6t5IA7yKX216SJxUs6vj-nkHsjJ-_sw1QdDCQ-ZKvb-kRAU8x9UHYiMaeQifIg3SVQt2jHqLzZpF1S1uOmyfJ1GOVIZt1MMy9wjYrtl-A"},
			want:    "3aef9b5c-883f-412d-ade7-47be43827d68",
			wantErr: false,
		},
		{
			name:    "test invalid token",
			fields:  fields{PublicKey: config.Config.GetAuthServerJwtPublicKey()},
			args:    args{tokenString: "invalid token"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "test expired token",
			fields:  fields{PublicKey: config.Config.GetAuthServerJwtPublicKey()},
			args:    args{tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0ODk5NTgsIm5iZiI6MTU4NTQ4OTk1OCwianRpIjoiNWE4OGNlYTctMzQ5NC00MThiLTlmMWYtMGQ5NjI0NDkwNjU1IiwiZXhwIjoxNTg1NDkwODU4LCJpZGVudGl0eSI6MSwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6ImFkbWluIiwidXVpZCI6IjNhZWY5YjVjLTg4M2YtNDEyZC1hZGU3LTQ3YmU0MzgyN2Q2OCIsInJlc291cmNlcyI6W119fQ.VC0_icto48HJC1OxezdzeccxyqqpSh6LB1xiSSPNKaiboY3A7bdPA2Zsv7zmUMfkVIeOSaqHkwE2YgjT_LRjTL1pb4j3T0Q_Zgkg5f2eDOi3M-HCVtE4NuZrT8bQJbBO5Fb8HqpSxWwSW-PJAX8jHHMPYYF5dytNOfRJ12amqPNvLk2qW-8CcxMXoi2i4lEF4NDfQqe1kx3xCQI9n1njBP0mrpHuz6O17M13gsj5oeTDiAIlmEQom8Yg7CHqqBjN6Pnihsrx4oghON5v_UGnMCnegKeisIjuwAbMk1yTzSArf_FSl9VSznV_vcsJkgd0_YacrCoqiJ7sBrZYlPlEoQ"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "test token without user uuid",
			fields:  fields{PublicKey: config.Config.GetAuthServerJwtPublicKey()},
			args:    args{tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0ODk5NTgsIm5iZiI6MTU4NTQ4OTk1OCwianRpIjoiNWE4OGNlYTctMzQ5NC00MThiLTlmMWYtMGQ5NjI0NDkwNjU1IiwiZXhwIjoxNTg1NDkwODU4LCJpZGVudGl0eSI6MSwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIn0.cd7vYYNGSBpBWc-miVSrwIyx4s5GtzKFvfPMmL76MxUpiLGQOPXeK7Q2jDz3LFRl0wGuD6ldm_nspGVgWXApII95a2nstqcBcwxR-HainxC5hr2mF0-D2tHc6IGSLyUmZ985blgqZVZDx1RnpXnZ6WPn5Nlumg1Az4lxssI-tlWjfq-gr5iEEbvrj1fARYt4aOXd-ppghQwQf2yBJ_uysgjWFq6OZp41K6_j2NlPbsSldrcv1DVNDk_WGxcFC1o0BSS_uarvSYlfOjet12BDRrcyhJiPgd57nN8Xb9DNez-ACV4dYck74ntWF8V7UXkzhx-opOxhNNemTQoe0mp40g"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "test token with corrupted signature",
			fields:  fields{PublicKey: config.Config.GetAuthServerJwtPublicKey()},
			args:    args{tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpYXQiOjE1ODU0ODk5NTgsIm5iZiI6MTU4NTQ4OTk1OCwianRpIjoiNWE4OGNlYTctMzQ5NC00MThiLTlmMWYtMGQ5NjI0NDkwNjU1IiwiZXhwIjoxNTg1NDkwODU4LCJpZGVudGl0eSI6MSwiZnJlc2giOmZhbHNlLCJ0eXBlIjoiYWNjZXNzIiwidXNlcl9jbGFpbXMiOnsicm9sZSI6ImFkbWluIiwidXVpZCI6IjNhZWY5YjVjLTg4M2YtNDEyZC1hZGU3LTQ3YmU0MzgyN2Q2OCIsInJlc291cmNlcyI6W119fQ.VC0_icto48HJC1OxezdzeccxyqqpSh6LB1xiSSPNKaiboY3A7bdPA2Zsv7zmUMfkVIeOSaqHkwE2YgjT_LRjTL1pb4j3T0Q_Zgkg5f2eDOi3M-HCVtE4NuZrT8bQJbBO5Fb8HqpSxWwSW-PJAX8jHHMPYYF5dytNOfRJ12amqPNvLk2qW-8CcxMXoi2i4lEF4NDfQqe1kx3xCQI9n1njBP0mrpHuz6O17M13gsj5oeTDiAIlmEQom8Yg7CHqqBjN6Pnihsrx4oghON5v_UGnMCnegKeisIjuwAbMk1yTzSArf_FSl9VSznV_vcsJkgd0_YacrCoqiJ7sBrZ11111"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &JwtTokenReader{
				PublicKey: tt.fields.PublicKey,
			}
			got, err := v.ReadUserUuid(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUserUuid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUserUuid() got = %v, want %v", got, tt.want)
			}
		})
	}
}
