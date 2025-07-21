package dtos

// import (
// 	"testing"
// )

// func TestCreatePlayerDto_Validate(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		dto     CreatePlayerDto
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid dto",
// 			dto: CreatePlayerDto{
// 				Name:     "testuser",
// 				Password: "password123",
// 				Email:    "test@example.com",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "empty name",
// 			dto: CreatePlayerDto{
// 				Name:     "",
// 				Password: "password123",
// 				Email:    "test@example.com",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "short name",
// 			dto: CreatePlayerDto{
// 				Name:     "ab",
// 				Password: "password123",
// 				Email:    "test@example.com",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "long name",
// 			dto: CreatePlayerDto{
// 				Name:     "this_is_a_very_long_name_that_exceeds_fifty_characters",
// 				Password: "password123",
// 				Email:    "test@example.com",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "short password",
// 			dto: CreatePlayerDto{
// 				Name:     "testuser",
// 				Password: "12345",
// 				Email:    "test@example.com",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid email",
// 			dto: CreatePlayerDto{
// 				Name:     "testuser",
// 				Password: "password123",
// 				Email:    "invalid-email",
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.dto.Validate()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestPlayerAuthUpdateDto_Validate(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		dto     PlayerAuthUpdateDto
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid dto",
// 			dto: PlayerAuthUpdateDto{
// 				OldPassword: "oldpass",
// 				NewPassword: "newpassword",
// 				Email:       "new@example.com",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "empty old password",
// 			dto: PlayerAuthUpdateDto{
// 				OldPassword: "",
// 				NewPassword: "newpassword",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "short new password",
// 			dto: PlayerAuthUpdateDto{
// 				OldPassword: "oldpass",
// 				NewPassword: "12345",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid email",
// 			dto: PlayerAuthUpdateDto{
// 				OldPassword: "oldpass",
// 				NewPassword: "newpassword",
// 				Email:       "invalid-email",
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.dto.Validate()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestPlayerMatchUpdateDto_Validate(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		dto     PlayerMatchUpdateDto
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid dto",
// 			dto: PlayerMatchUpdateDto{
// 				MatchID: "match123",
// 				Won:     true,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "empty match ID",
// 			dto: PlayerMatchUpdateDto{
// 				MatchID: "",
// 				Won:     false,
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "whitespace match ID",
// 			dto: PlayerMatchUpdateDto{
// 				MatchID: "   ",
// 				Won:     true,
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.dto.Validate()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestPlayerLevelUpdateDto_Validate(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		dto     PlayerLevelUpdateDto
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid level",
// 			dto: PlayerLevelUpdateDto{
// 				Level: 5.5,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "zero level",
// 			dto: PlayerLevelUpdateDto{
// 				Level: 0,
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "negative level",
// 			dto: PlayerLevelUpdateDto{
// 				Level: -1,
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.dto.Validate()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestPlayerStateUpdateDto_Validate(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		dto     PlayerStateUpdateDto
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid state - Online",
// 			dto: PlayerStateUpdateDto{
// 				State: "Online",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "valid state - Offline",
// 			dto: PlayerStateUpdateDto{
// 				State: "Offline",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "valid state - InGame",
// 			dto: PlayerStateUpdateDto{
// 				State: "InGame",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "valid state - InQuery",
// 			dto: PlayerStateUpdateDto{
// 				State: "InQuery",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "invalid state",
// 			dto: PlayerStateUpdateDto{
// 				State: "InvalidState",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "empty state",
// 			dto: PlayerStateUpdateDto{
// 				State: "",
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.dto.Validate()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestIsValidEmail(t *testing.T) {
// 	tests := []struct {
// 		email string
// 		valid bool
// 	}{
// 		{"test@example.com", true},
// 		{"user.name@domain.co.uk", true},
// 		{"user+tag@example.org", true},
// 		{"invalid-email", false},
// 		{"@example.com", false},
// 		{"test@", false},
// 		{"", false},
// 		{"test.example.com", false},
// 	}

// 	for _, test := range tests {
// 		if isValidEmail(test.email) != test.valid {
// 			t.Errorf("Expected %s to be valid: %v", test.email, test.valid)
// 		}
// 	}
// }