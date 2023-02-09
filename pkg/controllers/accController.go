package controllers

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"photo-media/pkg/dao"
	"photo-media/pkg/util"
	"strconv"
	"time"
)

//TODO: Drop session after expires > now() IN MYSQL!

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	passwordStr := r.FormValue("password")
	accHash := dao.ReadHashFromUsername(username)

	if util.CompareHashAndPassword(passwordStr, accHash) {

		tokenFound, oldToken := isOldSession(username)
		if tokenFound {
			createSessionCookie(w, r, username, oldToken)
			return
		}

		newSessionCreated, token := isNewSessionCreated(username)
		if newSessionCreated {
			createSessionCookie(w, r, username, token)
			return
		}

	} else {
		getHeaderLinkTemplate(w, "error: wrong username or password", "http://localhost:8080/login.html")
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	passwordStr := r.FormValue("password")
	minUsernameLenght := 5

	if isInvalidMinLenght(minUsernameLenght, username, passwordStr) {
		getHeaderLinkTemplate(w, "We are sorry, username and password must be between 5 and 45 characters", "http://localhost:8080/signup.html")
		return
	}

	hashedPassword, err := util.HashPassword(passwordStr)
	util.CheckError("Error encrypting password", err)

	if isNewAccCreated(username, hashedPassword) {
		getHeaderLinkTemplate(w, "New account created!", "http://localhost:8080/")

	} else {
		getHeaderLinkTemplate(w, "Unable to create new account, try again", "http://localhost:8080/signup.html")
	}
}

func createSessionCookie(w http.ResponseWriter, r *http.Request, usernameStr, tokenStr string) {
	path := "http://localhost:8080/home.html"
	username := &http.Cookie{
		Name:   "username",
		Value:  usernameStr,
		MaxAge: 600,
		Path:   path,
	}
	token := &http.Cookie{
		Name:   "token",
		Value:  tokenStr,
		MaxAge: 600,
		Path:   path,
	}
	http.SetCookie(w, username)
	http.SetCookie(w, token)
	http.Redirect(w, r, path, http.StatusMovedPermanently)
}

func isOldSession(username string) (bool, string) {
	tokenFound, sessionToken := dao.ReadSessionTokenByUsername(username)
	return tokenFound, sessionToken
}

func getHeaderLinkTemplate(w http.ResponseWriter, title, linkRef string) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h3>"+title+"</h3>")
	fmt.Fprintf(w, "<a href=\""+linkRef+"\">Go to main page</a>")
}

func isNewAccCreated(username, hashedPassword string) bool {
	return dao.CreateAcc(username, hashedPassword)
}

func isInvalidMinLenght(minLenght int, username, passwordStr string) bool {
	return len(username) < minLenght || len(passwordStr) < minLenght
}

func isNewSessionCreated(username string) (bool, string) {
	sessionToken := hex.EncodeToString([]byte(username + strconv.Itoa(time.Now().Second())))
	return dao.CreateSession(username, sessionToken), sessionToken
}
