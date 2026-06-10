package service

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"sespima_api/internal/repository"
	"sespima_api/models"
	"sespima_api/pkg/hash"
	pkgjwt "sespima_api/pkg/jwt"
)

type AuthService struct {
	users  repository.UserRepository
	secret string
}

func NewAuthService(users repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{users: users, secret: jwtSecret}
}

type LoginResult struct {
	User         *models.User
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) Login(nrpNip, password, fcmToken string) (*LoginResult, error) {
	user, err := s.users.FindByNrpOrNip(nrpNip)
	if err != nil {
		return nil, fmt.Errorf("NRP/NIP atau password salah")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("akun tidak aktif")
	}
	if !hash.CheckPassword(password, user.Password) {
		return nil, fmt.Errorf("NRP/NIP atau password salah")
	}

	nrp, nip := "", ""
	if user.NRP != nil {
		nrp = *user.NRP
	}
	if user.NIP != nil {
		nip = *user.NIP
	}

	tokens, err := pkgjwt.GenerateTokenPair(s.secret, user.ID, user.Email, user.Role, nrp, nip)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token")
	}

	user.CurrentToken = &tokens.AccessToken
	user.RefreshToken = &tokens.RefreshToken
	if fcmToken != "" {
		user.FcmToken = &fcmToken
	}
	if err := s.users.Save(user); err != nil {
		return nil, fmt.Errorf("gagal menyimpan sesi: %v", err)
	}

	return &LoginResult{User: user, AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *AuthService) GenerateResetToken(nrpNip string) (string, error) {
	user, err := s.users.FindByNrpOrNip(nrpNip)
	if err != nil {
		return "", fmt.Errorf("user tidak ditemukan")
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	token := fmt.Sprintf("%06d", n.Int64()+100000)
	user.ResetToken = &token
	return token, s.users.Save(user)
}

func (s *AuthService) VerifyResetToken(nrpNip, token string) error {
	user, err := s.users.FindByNrpOrNip(nrpNip)
	if err != nil {
		return fmt.Errorf("NRP tidak ditemukan")
	}
	if user.ResetToken == nil || *user.ResetToken != token {
		return fmt.Errorf("token tidak valid")
	}
	return nil
}

func validatePasswordStrength(pw string) error {
	if len(pw) < 8 {
		return fmt.Errorf("password minimal 8 karakter")
	}
	hasUpper, hasDigit := false, false
	for _, r := range pw {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}
	if !hasUpper {
		return fmt.Errorf("password harus mengandung huruf kapital")
	}
	if !hasDigit {
		return fmt.Errorf("password harus mengandung angka")
	}
	return nil
}

func (s *AuthService) ResetPassword(nrpNip, token, newPassword string) error {
	if err := validatePasswordStrength(newPassword); err != nil {
		return err
	}
	user, err := s.users.FindByNrpOrNip(nrpNip)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}
	if user.ResetToken == nil || *user.ResetToken != token {
		return fmt.Errorf("token tidak valid")
	}
	hashed, err := hash.Password(newPassword)
	if err != nil {
		return fmt.Errorf("gagal mengenkripsi password")
	}
	user.Password = hashed
	user.ResetToken = nil
	user.IsFirstLogin = false
	return s.users.Save(user)
}

func (s *AuthService) UpdatePassword(userID uint, newPassword string) error {
	if err := validatePasswordStrength(newPassword); err != nil {
		return err
	}
	user, err := s.users.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}
	hashed, err := hash.Password(newPassword)
	if err != nil {
		return fmt.Errorf("gagal mengenkripsi password")
	}
	user.Password = hashed
	user.IsFirstLogin = false
	return s.users.Save(user)
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	return s.users.FindByID(userID)
}

func (s *AuthService) Logout(userID uint) error {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}
	user.CurrentToken = nil
	user.RefreshToken = nil
	return s.users.Save(user)
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (*LoginResult, error) {
	userID, err := pkgjwt.ParseRefreshToken(s.secret, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh token tidak valid")
	}
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user tidak ditemukan")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("akun tidak aktif")
	}
	if user.RefreshToken == nil || *user.RefreshToken != refreshToken {
		return nil, fmt.Errorf("refresh token tidak valid atau sudah kadaluarsa")
	}

	nrp, nip := "", ""
	if user.NRP != nil {
		nrp = *user.NRP
	}
	if user.NIP != nil {
		nip = *user.NIP
	}

	tokens, err := pkgjwt.GenerateTokenPair(s.secret, user.ID, user.Email, user.Role, nrp, nip)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token")
	}

	user.CurrentToken = &tokens.AccessToken
	user.RefreshToken = &tokens.RefreshToken
	if err := s.users.Save(user); err != nil {
		return nil, fmt.Errorf("gagal menyimpan sesi: %v", err)
	}

	return &LoginResult{User: user, AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}
