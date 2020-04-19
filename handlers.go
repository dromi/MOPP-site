package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dromi/MOPP-site/model"
)

func frontpageHandler(w http.ResponseWriter, r *http.Request, m *MetaData) {
	renderTemplate(w, "frontpage", m)
}

type DateAvailable struct {
	Date      string
	Available []string
}

type Calendar struct {
	Performers []model.Performer
	Avails     []*DateAvailable
}

func dateAvailableFromAvailability(avails []model.Availability) []*DateAvailable {
	dateAvails := []*DateAvailable{}
	for _, v := range avails {
		date := v.Date.Format("02/01-06")
		reports := make([]string, len(v.Reports.Map))
		for j, h := range v.Reports.Map {
			// Ignore error here for now
			idx, _ := strconv.Atoi(j)
			reports[idx-1] = h.String
		}
		dateAvails = append(dateAvails, &DateAvailable{Date: date, Available: reports})
	}
	return dateAvails
}

func calendarHandler(w http.ResponseWriter, r *http.Request, m *MetaData) {
	performers := model.ListPerformers()
	avails := dateAvailableFromAvailability(model.ListAvailability())
	cal := &Calendar{Performers: performers, Avails: avails}
	renderTemplate(w, "calendar", cal)
}

func performerHandler(w http.ResponseWriter, r *http.Request, m *MetaData) {
	performers := model.ListPerformers()
	renderTemplate(w, "performers", performers)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		_, err := getJWTCookie(r)
		if err == nil {
			http.Redirect(w, r, "/frontpage", http.StatusFound)
		}
		renderTemplate(w, "signin", nil)
	case "POST":
		name := r.FormValue("uname")
		password := r.FormValue("psw")
		performer := model.GetPerformerByName(name)
		if performer.Name != "" {
			if ArePasswordsMatching(&performer.Password, &password) {
				cookie, err := CreateJWTCookie(name)
				if err != nil {
					panic(err)
				}
				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/frontpage", http.StatusFound)
			} else {
				renderTemplate(w, "signin", "Wrong username or password")
			}
		} else {
			renderTemplate(w, "signin", "Wrong username or password")
		}
	default:
		http.NotFound(w, r)
	}
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	claims, err := getJWTCookie(r)
	if err != nil {
		return
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < jwtRefreshTime {
		cookie, err := updateJWTCookie(claims)
		if err != nil {
			panic(err)
		}
		http.SetCookie(w, cookie)
	}
}
