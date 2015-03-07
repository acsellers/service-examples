package main

import (
	"encoding/json"
	"time"

	"git.andrewcsellers.com/acsellers/web2015/store"
	"github.com/acsellers/platform/router"
)

type SuggestCtrl struct {
	*router.BaseController
}

/*
	Routes for SuggestCtrl look like:

	GET    /suggestion?term=elephant  SuggestCtrl.Index
	POST   /suggestion with JSON Body SuggestCtrl.Create
	DELETE /suggestion/:suggestionid  SuggestCtrl.Delete
	GET    /suggestion/:suggestionid  SuggestCtrl.Show
*/
func NewSuggestCtrl() router.Controller {
	return SuggestCtrl{&router.BaseController{}}
}

func (SuggestCtrl) Path() string {
	return "suggestion"
}

func (sc SuggestCtrl) PreFilter() router.Result {
	auth_token := sc.Request.Header.Get("AUTH_TOKEN")
	if ValidateToken(auth_token) {
		sc.Log.Println("Accepted Token:", auth_token)
	} else {
		return &router.NotAllowed{}
	}
	return nil
}

// ValidateToken would make a request to an external service
// to determine whether the token is valid.
func ValidateToken(token string) bool {
	return true
}

func (sc SuggestCtrl) Index() router.Result {
	scope := Conn.Post.Published().Eq(true)
	if term := sc.Request.URL.Query().Get("search"); term != "" {
		scope.Name().Like("%" + term + "%")
	}

	if sc.Request.URL.Query().Get("recent") != "" {
		scope.CreatedAt().Gt(time.Now().AddDate(0, 0, -14))
	}

	if sc.Request.URL.Query().Get("all") != "" {
		scope.Limit(15)
	}

	posts, err := scope.RetrieveAll()
	if err != nil {
		return &router.String{"Internal Error", 500}
	} else {
		return router.JSON(posts)
	}
}

func (sc SuggestCtrl) Create() router.Result {
	var posts struct {
		Posts []store.Post `json:"posts"`
	}
	err := json.NewDecoder(sc.Request.Body).Decode(&posts)
	if err != nil {
		sc.Log.Println("JSON Error:", err)
	}

	var imported int
	for _, post := range posts.Posts {
		dupeid, err := Conn.Post.Permalink().Eq(post.Permalink).ID().PluckInt()
		if err != nil {
			sc.Log.Println("Update Error:", err)
			sc.Log.Println("On Post:", post)
			continue
		}
		if len(dupeid) > 0 {
			post.ID = int(dupeid[0])
		}
		err = post.Save(Conn)
		if err != nil {
			sc.Log.Println("Update Error:", err)
			sc.Log.Println("On Post:", post)
			continue
		}
		imported++
	}

	return router.JSON(imported)
}

func (sc SuggestCtrl) Delete() router.Result {
	err := Conn.Post.Permalink().Eq(sc.ID["suggestion"]).Delete()
	if err != nil {
		sc.Log.Println("Delete Error:", err)
		sc.ResponseWriter.WriteHeader(400)
		return router.JSON(map[string]string{"error": err.Error()})
	}
	return router.JSON("ok")
}

func (sc SuggestCtrl) Show() router.Result {
	exists := Conn.Post.Permalink().Eq(sc.ID["suggestion"]).Count() == 1
	return router.JSON(map[string]bool{"exists": exists})
}
