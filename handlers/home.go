package handlers

import (
	"encoding/json"
	"github.com/go-resty/resty"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/immutability-org/vaulty/libhttp"
	"github.com/immutability-org/vaulty/models"
	"html/template"
	"net/http"
)

type SealStatus struct {
	Sealed      bool   `json:"sealed"`
	T           int    `json:"t"`
	N           int    `json:"n"`
	Progress    int    `json:"progress"`
	Version     string `json:"version"`
	ClusterName string `json:"cluster_name"`
	ClusterID   string `json:"cluster_id"`
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := context.Get(r, "sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "vaulty-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}
	resp, err := resty.R().Get("https://127.0.0.1:8200/v1/sys/seal-status")
	var sealStatus SealStatus
	var sealStatus2 SealStatus

	json.Unmarshal(resp.Body(), &sealStatus)
	json.Unmarshal(resp.Body(), &sealStatus2)

	payload := []SealStatus{sealStatus, sealStatus2}

	data := struct {
		CurrentUser *models.UserRow
		Payload     []SealStatus
	}{
		currentUser,
		payload,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl", "templates/messages/menu.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
