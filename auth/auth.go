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

//CreateSession ..
func CreateSession(name string, sID string, w http.ResponseWriter, r *http.Request) error {
	session, err := SessionStore.Get(r, "Allsessions")
	if err != nil {
		fmt.Println(err)
		log.Fatal()
	}
	session.Values["sessionID"] = sID
	session.Values["username"] = name

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

//CreateCookie ..
func CreateCookie(name string, sID string, w http.ResponseWriter, r *http.Request) error {
	val := map[string]string{
		"username":  name,
		"sessionId": sID,
	}

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
func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "mycookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

//ReadCookie ..
func ReadCookie(w http.ResponseWriter, r *http.Request) (map[string]string, error) {
	if cookie, err := r.Cookie("mycookie"); err == nil {
		val := make(map[string]string)
		if err = CookieHandler.Decode("mycookie", cookie.Value, &val); err == nil {
			return val, err
		}
		return nil, err
	}
	return nil, nil
}

//LoginForm ..
// func LoginForm(w)
