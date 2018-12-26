package server

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"regexp"
	"time"
)

func GenerateSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func DetectLanguage(acceptLanguage string) string {
	regExpEn := regexp.MustCompile(".*en.*")
	regExpTr := regexp.MustCompile(".*tr.*")

	if regExpEn.MatchString(acceptLanguage) {
		return "en"
	}
	if regExpTr.MatchString(acceptLanguage) {
		return "tr"
	}
	//Default
	return "en"
}

func DeleteCookieHandler(rw http.ResponseWriter, cookieName string) {
	c := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(rw, c)
}

func ValidateAge(age int) bool {
	return age >= 0 && age <= 99
}
