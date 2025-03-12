package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Notification struct {
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Subject struct {
		Title string `json:"title"`
		Type  string `json:"type"`
		URL   string `json:"url"`
	} `json:"subject"`
	Reason string `json:"reason"`
	Link   string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN not set in .env")
	}

	log.Println("GitHub token loaded successfully!")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/notifications", notificationsHandler)

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Couldn't start the server:", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/notifications.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	notifications, err := fetchGithubNotifications()
	if err != nil {
		log.Println("Failed to fetch notifications:", err)
		notifications = []Notification{}
	}

	if err := tmpl.Execute(w, notifications); err != nil {
		log.Println("Template execution error:", err)
	}
}

func notificationsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/notifications.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	notifications, err := fetchGithubNotifications()
	if err != nil {
		log.Println("Failed to fetch notifications:", err)
		notifications = []Notification{}
	}

	if err := tmpl.ExecuteTemplate(w, "notifications.html", notifications); err != nil {
		log.Println("Template execution error:", err)
	}
}

func fetchGithubNotifications() ([]Notification, error) {
	token := os.Getenv("GITHUB_TOKEN")

	req, err := http.NewRequest("GET", "https://api.github.com/notifications", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s\n%s", resp.Status, string(body))
	}

	var notifications []Notification
	if err := json.NewDecoder(resp.Body).Decode(&notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}
