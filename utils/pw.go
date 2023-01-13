package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/TranQuocToan1996/redislearn/config"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

const (
	constantTimeMatch = 1
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
	ErrNotMatch            = errors.New("notmatch of argon2")
)

var (
	Pw passworder
)

// func init() {
// 	Pw = &bcryptImpl{
// 		cost: bcrypt.DefaultCost,
// 	}
// }

type passworder interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword string, candidatePassword string) error
}

type bcryptImpl struct {
	cost int
}

func (b *bcryptImpl) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func (b *bcryptImpl) VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

// How to set params
// Set the parallelism and memory parameters to the largest amount you are willing to afford, bearing in mind that you probably don't want to max these out completely unless your machine is dedicated to password hashing.
// Increase the number of iterations until you reach your maximum runtime limit (for example, 500ms).
// If you're already exceeding the your maximum runtime limit with the number of iterations = 1, then you should reduce the memory parameter.
type argon2id struct {
	memory      uint32 // The amount of memory used by the algorithm (in kibibytes)
	iterations  uint32 // The number of iterations (or passes) over the memory
	parallelism uint8  // The number of threads (or lanes) used by the algorithm. Change this one will change the hashing output
	saltLength  uint32 // Length of the random salt. 16 bytes is recommended for password hashing
	keyLength   uint32 // Length of the generated key (or password hash). 16 bytes or more is recommended.
}

func NewArgon(cfg config.Config) *argon2id {
	return &argon2id{
		memory:      cfg.ARGON2IDMemory,
		iterations:  cfg.ARGON2IDIteration,
		parallelism: cfg.ARGON2IDParallelsism,
		saltLength:  cfg.ARGON2IDSaltLength,
		keyLength:   cfg.ARGON2IDKeyLength,
	}
}

func (a *argon2id) HashPassword(password string) (string, error) {
	salt, err := a.generateRandomBytes(a.saltLength)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, a.iterations, a.memory, a.parallelism, a.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	// $argon2id — the variant of Argon2 being used.
	// $v=19 — the version of Argon2 being used.
	// $m=65536,t=3,p=2 — the memory (m), iterations (t) and parallelism (p) parameters being used.
	// $c29tZXNhbHQ — the base64-encoded salt, using standard base64-encoding and no padding.
	// $RdescudvJCsgt3ub+b+dWRWJTmaaJObG — the base64-encoded hashed password (derived key), using standard base64-encoding and no padding.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.memory, a.iterations, a.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (a *argon2id) VerifyPassword(hashedPassword string, candidatePassword string) error {
	match, err := a.comparePasswordAndHash(hashedPassword, candidatePassword)
	if match {
		return nil
	}
	return err
}

func (a *argon2id) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (a *argon2id) comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := a.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == constantTimeMatch {
		return true, nil
	}
	return false, nil
}

func (a *argon2id) decodeHash(encodedHash string) (p *argon2id, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &argon2id{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
