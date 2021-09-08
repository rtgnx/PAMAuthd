package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type PasswdEntry struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	UID      uint   `json:"uid"`
	GID      uint   `json:"gid"`
	Fullname string `json:"fullname"`
	Home     string `json:"home"`
	Shell    string `json:"shell"`
}

type Passwd []PasswdEntry

func (p *PasswdEntry) Decode(b []byte) error {

	items := strings.Split(string(b), ":")

	if len(items) < 7 {
		return fmt.Errorf("invalid passwd entry string, expected 7 values colon delimited, got %d", len(items))
	}

	p.Name, p.Password, p.Fullname, p.Home, p.Shell = items[0], items[1], items[4], items[5], items[6]

	var uid, gid uint64
	var uidErr, gidErr error

	uid, uidErr = strconv.ParseUint(items[2], 10, 32)

	gid, gidErr = strconv.ParseUint(items[2], 10, 32)

	if uidErr != nil {
		return uidErr
	} else if gidErr != nil {
		return gidErr
	}

	p.UID, p.GID = uint(uid), uint(gid)

	return nil
}

func (p PasswdEntry) Encode() []byte {
	return []byte(
		fmt.Sprintf(
			"%s:%s:%d:%d:%s:%s:%s",
			p.Name, p.Password, p.UID, p.GID, p.Fullname, p.Home, p.Shell,
		),
	)
}

func (p *Passwd) ReadFrom(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	for _, entry := range strings.Split(string(buf), "\n") {
		v := new(PasswdEntry)
		if err := v.Decode([]byte(entry)); err != nil {
			log.Println(err.Error())
			continue
		}

		*p = append(*p, *v)
	}

	return nil
}

func (p Passwd) FindByName(name string) (PasswdEntry, bool) {
	for _, entry := range p {
		if strings.Compare(entry.Name, name) == 0 {
			return entry, true
		}
	}

	return PasswdEntry{}, false
}

func FetchPasswdFile() (Passwd, error) {
	passwd := make(Passwd, 0)
	fd, err := os.Open("/etc/passwd")

	if err != nil {
		return Passwd{}, nil
	}

	defer fd.Close()

	err = passwd.ReadFrom(fd)
	return passwd, err
}
