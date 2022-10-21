package handlers

import (
	"Go-server/db"
	"Go-server/server/sessions"
	"Go-server/view"
	"fmt"
	"net/http"
)

func RenderAccountPage(w http.ResponseWriter, r *http.Request,
	renderer *view.Renderer, sm *sessions.Manager) {
	if !sm.SessionExists(w, r) {
		http.Redirect(w, r, "/login", 301)
	}

	if err := renderer.RenderTemplate(w, "account"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request, pdb *db.PstgDB,
	sm *sessions.Manager, renderer *view.Renderer) {

	if r.Method == "GET" {
		if sm.SessionExists(w, r) {
			http.Redirect(w, r, "/account", 301)
		}

		if err := renderer.RenderTemplate(w, "login"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Println(err.Error())
			return
		}

		login := r.FormValue("login")
		psw := r.FormValue("psw")

		if login == "" || psw == "" {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		account, err := pdb.SelectAccount(login, psw)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		sid := sm.CreateSession(account.Id)
		sm.SetCookie(w, "session", sid)
		http.Redirect(w, r, "/account", 301)
	} else {
		http.Error(w, http.StatusText(405), 405)
	}
}

func Logout(w http.ResponseWriter, r *http.Request, pdb *db.PstgDB, sm *sessions.Manager) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	cookie, _ := r.Cookie("session")
	sm.DeleteSession(cookie.Value)
	sm.UnsetCookie(w, "session")
	http.Redirect(w, r, "/login", 301)
}

func Register(w http.ResponseWriter, r *http.Request, pdb *db.PstgDB,
	renderer *view.Renderer) {
	if r.Method == "GET" {
		if err := renderer.RenderTemplate(w, "register"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Println(err.Error())
		}

		login := r.FormValue("login")
		psw := r.FormValue("psw")
		pswRepeat := r.FormValue("psw-repeat")

		if login == "" || psw == "" || pswRepeat == "" {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		if acc, _ := pdb.SelectAccountLogin(login); acc != nil {
			fmt.Fprintln(w, "Account already exist")
			return
		}

		if psw != pswRepeat {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		if err := pdb.CreateAccount(login, psw); err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		} else {
			http.Redirect(w, r, "/login", 301)
		}
	} else {
		http.Error(w, http.StatusText(405), 405)
	}
}
