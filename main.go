package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Redirect defines a redirect, obviously
type Redirect struct {
	gorm.Model
	RedirectTo string `gorm:"not_null"`
}

// RedirectHit records each open on a Redirect
type RedirectHit struct {
	gorm.Model
	Redirect   Redirect
	RedirectID int
	IP         string
	Browser    Browser
	BrowserID  int
	Tracking   string
}

// Browser defines a browser
type Browser struct {
	gorm.Model
	UserAgent string `gorm:"unique"`
}

var server *http.Server
var db *gorm.DB

func main() {
	db, _ = gorm.Open("mysql", "root:@/shorturls?charset=utf8&parseTime=True&loc=Local")
	db.AutoMigrate(&Redirect{}, &RedirectHit{}, &Browser{})
	defer db.Close()

	establishHTTPServer()

	log.Fatal(server.ListenAndServe())
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "" {
		if r.Method == "POST" {
			// Handle creating a new short URL
		} else {
			// Show index page
		}
	} else {
		handleShortURL(w, r)
	}
}

func handleShortURL(w http.ResponseWriter, r *http.Request) {
	// Comes through as `/key/other/stuff` -- we just want `key`
	splits := strings.SplitN(r.URL.Path[1:], "/", 2)
	key := splits[0]

	var tracking string
	if len(splits) == 2 {
		tracking = splits[1]
	} else {
		tracking = ""
	}

	keyInt, _ := strconv.ParseInt(key, 36, 64)

	var redirect Redirect
	if err := db.Where("id = ?", keyInt).First(&redirect).Error; err != nil {
		fmt.Fprintf(w, "Could not find %s", key)
		fmt.Println("Finished with error response")
		return
	}

	browser := Browser{
		UserAgent: r.UserAgent(),
	}

	db.FirstOrCreate(&browser, browser)

	redirectHit := RedirectHit{
		Redirect: redirect,
		IP:       r.RemoteAddr,
		Browser:  browser,
		Tracking: tracking,
	}

	db.Save(&redirectHit)

	//responseString := fmt.Sprintf("Did find %s -- redirecting to %s", key, redirect.RedirectTo)
	//io.WriteString(w, responseString)
	w.Header().Set("Location", redirect.RedirectTo)
	w.WriteHeader(302)
	fmt.Println("Finished with success response")
}

func establishHTTPServer() {
	server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: http.HandlerFunc(httpHandler),
	}
}
