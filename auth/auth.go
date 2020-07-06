package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

//HashKey ..
var HashKey = securecookie.GenerateRandomKey(32)

//BlockKey ..
var BlockKey = securecookie.GenerateRandomKey(32)

//CookieHandler ...
var CookieHandler = securecookie.New(HashKey, BlockKey)

//SessionStore ...
var SessionStore = sessions.NewFilesystemStore("/tmp", HashKey)

//Identity ..
type Identity struct {
	userID    int
	sessionID string
}

//CreateSession ..
func CreateSession(uID int, sID string,
	w http.ResponseWriter, r *http.Request) error {
	session, err := SessionStore.Get(r, "Allsessions")
	if err != nil {
		fmt.Println(err)
		log.Fatal()
	}
	session.Values["sessionID"] = sID
	session.Values["userID"] = uID

	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err)
		log.Fatal()
	}
	return nil
}

//ClearSession ..
func ClearSession(w http.ResponseWriter, r *http.Request) {

	session, err := SessionStore.Get(r, "Allsessions")
	if err != nil {
		fmt.Println(err)
		log.Fatal()
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
}

//CheckSession .
func CheckSession(sID string, r *http.Request) bool {
	session, err := SessionStore.Get(r, "Allsessions")
	if err != nil {
		log.Fatalln(err)
	}
	if session.Values["sessionID"] == sID {
		return true
	}
	return false
}

//CreateCookie ..
func CreateCookie(uID int, sID string, w http.ResponseWriter) error {
	var val Identity
	val.userID = uID
	val.sessionID = sID

	if encode, err := CookieHandler.Encode("mycookie", val); err == nil {
		cookie := &http.Cookie{
			Name:   "mycookie",
			Path:   "/",
			Value:  encode,
			MaxAge: 3600,
		}
		http.SetCookie(w, cookie)
	} else {
		return err
	}
	return nil
}

//DeleteCookie ..
func DeleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "mycookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

//ReadCookie ..
func ReadCookie(r *http.Request) (Identity, error) {
	var val Identity
	if cookie, err := r.Cookie("mycookie"); err == nil {
		if err = CookieHandler.Decode("mycookie", cookie.Value, &val); err == nil {
			return val, err
		}
		return val, err
	}
	return val, nil
}

//LoginForm ..
// func LoginForm(w)
