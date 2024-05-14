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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type userEmails [100_000]string

func getUsers(r io.Reader) (result userEmails, err error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	decoder := json.NewDecoder(r)
	for i := 0; decoder.More(); i++ {
		var user User
		if err = decoder.Decode(&user); err != nil {
			return
		}
		result[i] = user.Email
	}

	return
}

func countDomains(uE userEmails, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, email := range uE {
		if strings.HasSuffix(strings.ToLower(email), "."+strings.ToLower(domain)) {
			emailKey := strings.ToLower(strings.SplitN(email, "@", 2)[1])

			result[emailKey]++
		}
	}
	return result, nil
}
