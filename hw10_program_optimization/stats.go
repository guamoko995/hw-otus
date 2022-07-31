package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

var regDom = regexp.MustCompile("[a-z][a-z]*")

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	line := bufio.NewScanner(r)
	result := make(DomainStat)
	if regDom.FindString(domain) != domain {
		return nil, errors.New("invalid domain")
	}
	regDomain, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}
	for i := 0; line.Scan(); i++ {
		b := line.Bytes()
		var user User
		if err := user.UnmarshalJSON(b); err != nil {
			return nil, err
		}
		if regDomain.MatchString(user.Email) {
			str := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[str]++
		}
	}

	return result, nil
}
