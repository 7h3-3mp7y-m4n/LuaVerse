package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	lua "github.com/yuin/gopher-lua"
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
	Link   string // Browser-friendly link
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
	http.HandleFunc("/notifications", notificationHandler)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Couldn't start the server:", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/notifications.html")
	if err != nil {
		http.Error(w, "Template loading error", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	filtered, err := getFilteredNotifications()
	if err != nil {
		log.Println("Failed to fetch notifications:", err)
		filtered = []Notification{}
	}

	if err := tmpl.ExecuteTemplate(w, "index", filtered); err != nil {
		log.Println("Template execution error:", err)
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}
}

func notificationHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/notifications.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Println("Template load error:", err)
		return
	}

	filtered, err := getFilteredNotifications()
	if err != nil {
		log.Println("Failed to fetch notifications:", err)
		filtered = []Notification{}
	}

	if err := tmpl.ExecuteTemplate(w, "notifications", filtered); err != nil {
		log.Println("Notification render error:", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
	}
}

func getFilteredNotifications() ([]Notification, error) {
	raw, err := fetchGithubNotifications(os.Getenv("GITHUB_TOKEN"))
	if err != nil {
		return nil, err
	}

	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile("task/filter.lua"); err != nil {
		log.Println("Lua script error:", err)
		return []Notification{}, nil
	}

	var filtered []Notification
	for _, n := range raw {
		n.Link = convertAPIToLinks(n.Subject.URL)
		if shouldIncludeViaLua(L, n) {
			filtered = append(filtered, n)
		}
	}

	return filtered, nil
}

func fetchGithubNotifications(token string) ([]Notification, error) {
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

func convertAPIToLinks(apiURL string) string {
	if apiURL == "" {
		return "#"
	}
	parsed, err := url.Parse(apiURL)
	if err != nil {
		return "#"
	}
	parts := strings.Split(parsed.Path, "/")
	if len(parts) < 6 {
		return "#"
	}
	owner, repo, kind, id := parts[2], parts[3], parts[4], parts[5]

	switch kind {
	case "pulls":
		return fmt.Sprintf("https://github.com/%s/%s/pull/%s", owner, repo, id)
	case "issues":
		return fmt.Sprintf("https://github.com/%s/%s/issues/%s", owner, repo, id)
	default:
		return fmt.Sprintf("https://github.com/%s/%s", owner, repo)
	}
}

func shouldIncludeViaLua(L *lua.LState, notif Notification) bool {
	luaNotif := L.NewTable()
	luaNotif.RawSetString("title", lua.LString(notif.Subject.Title))
	luaNotif.RawSetString("type", lua.LString(notif.Subject.Type))
	luaNotif.RawSetString("reason", lua.LString(notif.Reason))
	luaNotif.RawSetString("repo", lua.LString(notif.Repository.FullName))

	L.SetGlobal("notification", luaNotif)

	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("allow"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		log.Println("Lua function call error:", err)
		return false
	}

	ret := L.Get(-1)
	L.Pop(1)

	if b, ok := ret.(lua.LBool); ok {
		return bool(b)
	}
	return false
}
