package encryption

import "golang.org/x/crypto/bcrypt"


type BcrypterOpt func(*Bcrypter)

func WithCost(cost int) BcrypterOpt {
	return func(b *Bcrypter) {
		b.cost = cost
	}
}

type Bcrypter struct {
	cost int
}

func (b *Bcrypter) Encrypt(in string) (string, error) {
	val, err := bcrypt.GenerateFromPassword([]byte(in), b.cost)

	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (b *Bcrypter) HashMatches(hash string, val string) (bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(val))

	return err == nil
}

func NewBcrypter(opts... BcrypterOpt) *Bcrypter {
	b := &Bcrypter{
		cost: 12,
	}

	for _, opt := range(opts) {
		opt(b)
	}

	return b
}