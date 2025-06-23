package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	const maxLine = 1 << 20 // 1 МБ
	scanner.Buffer(make([]byte, 64*1024), maxLine)

	var i int
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		// Parse only an email field - the rest we do not need here
		var tmp struct {
			Email string `json:"email"`
		}
		if err = json.Unmarshal(line, &tmp); err != nil {
			return
		}
		result[i].Email = tmp.Email
		i++
	}
	err = scanner.Err()
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	domain = "." + strings.ToLower(domain)
	stat := make(DomainStat)

	for _, user := range u {
		email := strings.ToLower(user.Email)
		if email == "" {
			continue
		}

		at := strings.LastIndexByte(email, '@')
		if at < 0 {
			continue
		}
		host := email[at+1:]

		if strings.HasSuffix(host, domain) {
			stat[host]++
		}
	}
	return stat, nil
}
