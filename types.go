package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type PasswdLine struct {
	Name     string   `json:"name"`
	Password string   `json:"password"`
	UID      int64    `json:"uid"`
	GID      int64    `json:"gid"`
	Fullname string   `json:"fullname"`
	Home     string   `json:"home"`
	Shell    string   `json:"shell"`
	Groups   []string `json:"groups"`
}

func (p PasswdLine) Marshal() ([]byte, error) {
	return []byte(fmt.Sprintf(
		"%s:%s:%d:%d:%s:%s:%s",
		p.Name, p.Password, p.UID, p.GID,
		p.Fullname, p.Home, p.Shell,
	)), nil
}

func (p *PasswdLine) UnmarshalText(text []byte) (err error) {
	items := strings.Split(string(text), ":")

	if len(items) < 7 {
		return fmt.Errorf("invalid passwd entry string, expected 7 values colon delimited, got %d", len(items))
	}

	p.Name, p.Password, p.Fullname, p.Home, p.Shell = items[0], items[1], items[4], items[5], items[6]

	if p.UID, err = strconv.ParseInt(items[2], 10, 32); err != nil {
		return err
	}

	if p.GID, err = strconv.ParseInt(items[2], 10, 32); err != nil {
		return err
	}

	return nil
}

func ParsePasswd(r io.Reader) map[string]PasswdLine {
	rd := bufio.NewReader(r)
	passwd := make(map[string]PasswdLine)

	for {
		line, err := rd.ReadBytes('\n')
		line = []byte(strings.TrimRight(string(line), "\n"))
		if err != nil {
			break
		}

		p := PasswdLine{}

		if err := p.UnmarshalText(line); err != nil {
			log.Println(err.Error())
			return passwd
		}
		passwd[p.Name] = p
	}

	return passwd
}
