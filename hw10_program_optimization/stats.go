package hw10programoptimization

import (
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

var ErrReaderIsNil = fmt.Errorf("reader is nil")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	users, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(users, domain)
}

func getUsers(r io.Reader) (*jsoniter.Decoder, error) {
	if r == nil {
		return nil, ErrReaderIsNil
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(r)

	return decoder, nil
}

func countDomains(usersDecoder *jsoniter.Decoder, domain string) (DomainStat, error) {
	domain = strings.ToLower(domain)
	result := make(DomainStat)

	for usersDecoder.More() {
		var user User
		if err := usersDecoder.Decode(&user); err != nil {
			return nil, err
		}
		email := strings.ToLower(user.Email)
		if strings.HasSuffix(email, "."+domain) {
			domainPart := strings.SplitN(email, "@", 2)[1]
			result[domainPart]++
		}
	}

	return result, nil
}
