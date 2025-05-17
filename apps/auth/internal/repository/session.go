package repository

import (
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

const (
	// SessionTTL represents the time to live for the session in seconds. By default, it is set to 1 year.
	SessionTTL = 60 * 60 * 24 * 365 // 1 second * Minutes * Hours * Days * Years
)

// Session represents a session in the database
type Session struct {
	SessionKey string    `json:"session_key" gorm:"primaryKey" gorm:"column:session_key"`
	LastUsed   time.Time `json:"last_used" gorm:"column:last_used"`
	ExpiresAt  time.Time `json:"expires_at" gorm:"column:expires_at"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at" gorm:"autoUpdateTime"`
	UserID     int64     `json:"-" gorm:"column:user_id"`
	User       *Auth     `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

func (Session) TableName() string {
	return "sessions"
}

func NewSession(user *Auth) *Session {

	return &Session{
		SessionKey: GenerateRandomString(64),
		LastUsed:   time.Unix(0, 0),
		ExpiresAt:  time.Now().Add(time.Second * SessionTTL),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		UserID:     user.ID,
	}
}

// SessionRepository represents the repository for the session
type SessionRepository interface {
	Create(session *Session) error
	GetAll() ([]*Session, error)
	Get(sessionKey string) (*Session, error)
	UpdateLastUsed(sessionKey string) error
	Delete(sessionKey string) error
	HardDelete(sessionKey string) error
	HardDeleteAllExpired() error
	HardDeleteAllInactive() error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (s sessionRepository) GetAll() ([]*Session, error) {
	logging.Logger.Info("Getting all sessions")
	var sessions []*Session
	err := s.db.Find(&sessions).Error
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to get all sessions")
		return nil, err
	}
	logging.Logger.Debug("Found ", len(sessions), " sessions")
	return sessions, nil
}

func (s sessionRepository) Create(session *Session) error {
	logging.Logger.Info("Creating session with key: ", session.SessionKey[:5], "...")
	err := s.db.Create(session).Error
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to create session with key: ", session.SessionKey[:5])
		return err
	}
	return nil
}

func (s sessionRepository) Get(sessionKey string) (*Session, error) {
	logging.Logger.Info("Getting session with key: ", sessionKey[:5], "...")
	var session Session
	err := s.db.Preload("User").Preload("User.Roles").Where("session_key = ?", sessionKey).Where("expires_at > ?", time.Now()).First(&session).Error
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to get session with key: ", sessionKey[:5])
		return nil, err
	}
	logging.Logger.Debug("Session found with key: ", sessionKey[:5], "...")
	return &session, nil
}

func (s sessionRepository) UpdateLastUsed(session string) error {
	logging.Logger.Info("Updating last used time for session with key: ", session[:5], "...")
	err := s.db.Model(&Session{}).Where("session_key = ?", session).Update("last_used", time.Now()).Error
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to update last used time for session with key: ", session[:5])
		return err
	}
	return nil
}

func (s sessionRepository) Delete(sessionKey string) error {
	logging.Logger.Info("Expiring session with key: ", sessionKey[:5], "...")
	err := s.db.Model(&Session{}).Where("session_key = ?", sessionKey).Update("expires_at", time.Now()).Error
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to expire session: ", sessionKey[:5])
		return err
	}
	return nil
}

func (s sessionRepository) HardDelete(sessionKey string) error {
	logging.Logger.Warn("Deleting session with key: ", sessionKey[:5], "...")
	err := s.db.Where("session_key = ?", sessionKey).Delete(&Session{}).Error
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to delete session with key: ", sessionKey[:5])
		return err
	}
	return nil
}

func (s sessionRepository) HardDeleteAllExpired() error {
	logging.Logger.Warn("Deleting all expired sessions...")
	err := s.db.Delete(&Session{}, "expires_at < ?", time.Now()).Error
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to delete all expired sessions")
		return err
	}
	return nil
}

func (s sessionRepository) HardDeleteAllInactive() error {
	logging.Logger.Warn("Deleting all inactive sessions...")
	err := s.db.Delete(&Session{}, "last_used < ?", time.Now().Add(-time.Second*SessionTTL)).Error
	if err != nil {
		logging.Logger.WithError(err).Error("Failed to delete all inactive sessions")
		return err
	}
	return nil
}

// TODO: Move this to a separate package

func GenerateRandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
