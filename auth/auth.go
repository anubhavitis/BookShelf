package auth

import (
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
var SessionStore = sessions.NewFilesystemStore("./store", HashKey)

//CreateSession ..
func CreateSession(uID string, sID string,
	w http.ResponseWriter, r *http.Request) error {
	session, err := SessionStore.Get(r, "Allsessions")
	if err != nil {
		return nil
	}

	session.Values["sessionID"] = sID
	session.Values["userID"] = uID

	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

//ClearSession ..
func ClearSession(w http.ResponseWriter, r *http.Request) error {

	session, err := SessionStore.Get(r, "Allsessions")
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	return nil
}

//CheckSession .
func CheckSession(sID string, r *http.Request) (bool, error) {
	session, err := SessionStore.Get(r, "Allsessions")
	if err != nil {
		return false, err
	}
	if session.Values["sessionID"] == sID {
		return true, nil
	}
	return false, nil
}

//CreateCookie ..
func CreateCookie(uID string, sID string, w http.ResponseWriter) error {
	val := map[string]string{
		"userID":    uID,
		"sessionID": sID,
	}
	if encode, err := CookieHandler.Encode("mycookie", val); err == nil {
		cookie := &http.Cookie{
			Name:     "mycookie",
			Path:     "/",
			Value:    encode,
			MaxAge:   3600,
			Secure:   false,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		return nil
	} else {
		return err
	}
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
func ReadCookie(r *http.Request) (map[string]string, error) {
	val := make(map[string]string)
	if cookie, err := r.Cookie("mycookie"); err == nil {
		if err = CookieHandler.Decode("mycookie", cookie.Value, &val); err == nil {
			return val, err
		}
		return nil, err
	}
	return nil, nil
}

//LoginForm ..
// func LoginForm(w)
