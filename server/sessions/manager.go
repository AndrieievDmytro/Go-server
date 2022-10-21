package sessions

import (
	"Go-server/db"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Manager struct {
	pdb *db.PstgDB
}

func NewManager(pdb *db.PstgDB) *Manager {
	return &Manager{
		pdb: pdb,
	}
}

func (o *Manager) CreateSession(uid string) string {
	sid := o.gen_SID32()
	tStamp := time.Now()
	if err := o.pdb.InsertSession(sid, uid, tStamp); err != nil {
		fmt.Println(err.Error())
		return ""
	}

	fmt.Printf("Session created. SID: %s\n", sid)
	return sid
}

func (o *Manager) DeleteSession(sid string) {
	if err := o.pdb.DeleteSessionSID(sid); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Session deleted. SID: %s\n", sid)
}

func (o *Manager) SetCookie(w http.ResponseWriter, name string, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (o *Manager) UnsetCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (o *Manager) SessionExists(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	s, _ := o.pdb.SelectSessionSID(cookie.Value)
	return s != nil
}

func (o *Manager) gen_SID32() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("SessionsManager: id creation fault")
	}

	return base64.URLEncoding.EncodeToString(b)
}
