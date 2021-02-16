package HttpAction

import (
	"Database"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type PostView struct {
	Posts []Database.Post
}

var mainPage *template.Template
var boardMain *template.Template
var boardPost *template.Template
var newPost *template.Template
var signUp *template.Template
var signIn *template.Template

func HandleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	// 여기서 각 페이지 핸들러들 등록
	myRouter.HandleFunc("/", handleHomePage)
	myRouter.HandleFunc("/posts", handleGetPosts)
	myRouter.HandleFunc("/posts/{id}", handleGetPostByID)
	myRouter.HandleFunc("/new_post", handleNewPost)
	myRouter.HandleFunc("/sign_up", handleSignUp)
	myRouter.HandleFunc("/sign_in", handleSingIn)
	myRouter.HandleFunc("/delete_post/{id}", handleDeletePost).Methods("POST")

	//myRouter.HandleFunc("/new_post", handleNewPostGet)

	mainPage = template.Must(template.ParseFiles("/home/jjun/workspace/toy_project/src/static/index.html"))
	boardMain = template.Must(template.ParseFiles("/home/jjun/workspace/toy_project/src/static/board/index.html"))
	boardPost = template.Must(template.ParseFiles("/home/jjun/workspace/toy_project/src/static/board/post.html"))
	newPost = template.Must(template.ParseFiles("/home/jjun/workspace/toy_project/src/static/board/newPost.html"))
	signUp = template.Must(template.ParseFiles("/home/jjun/workspace/toy_project/src/static/signUp.html"))
	signIn = template.Must(template.ParseFiles("/home/jjun/workspace/toy_project/src/static/signIn.html"))

	log.Fatal(http.ListenAndServe(":9090", myRouter))
}

func handleDeletePost(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: handleDeletePost")
	vars := mux.Vars(r)
	id := vars["id"]

	err := Database.DeletePost(id)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "http://localhost:9090/posts", http.StatusSeeOther)
}

func handleSingIn(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: handleSingIn", r.Method)
	if r.Method == "GET" {
		signIn.Execute(w, nil)
	} else if r.Method == "POST" {
		var user Database.User
		user.UserId = r.FormValue("Id")
		user.UserPwd = r.FormValue("Password")
		fmt.Println(user)
		id, err := Database.GetUser(user)
		if err != nil {
			panic(err)
		} else if id == -1 {
			http.Redirect(w, r, "http://localhost:9090/sign_in", http.StatusNotFound)
		}

		http.Redirect(w, r, "http://localhost:9090/posts", http.StatusSeeOther)
	}
}

func handleSignUp(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: handleSignUp", r.Method)
	if r.Method == "GET" {
		signUp.Execute(w, nil)
	} else if r.Method == "POST" {
		fmt.Print(r.Body)
		var newUser Database.User
		newUser.UserId = r.FormValue("Id")
		newUser.UserPwd = r.FormValue("Password")
		fmt.Println(newUser)
		err := Database.InsertUser(newUser)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "http://localhost:9090/", http.StatusSeeOther)
	}
}


func handleNewPost(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: handleNewPost", r.Method)
	if r.Method == "GET" {
		newPost.Execute(w, nil)
	} else if r.Method == "POST" {
		fmt.Println(r.Body)
		var post Database.Post
		post.Title = r.FormValue("Title")
		post.Text = r.FormValue("Text")
		fmt.Println(post)
		Database.InsertPost(post)

		http.Redirect(w, r, "http://localhost:9090/posts", http.StatusSeeOther)
	}
}

func handleGetPostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("Endpoint Hit: handleGetPostByID")

	post, err := Database.SelectByID(id)
	if err != nil {
		log.Println("Error Occurred", err)
	}
	//json.NewEncoder(w).Encode(post)

	fmt.Println(post)
	boardPost.Execute(w, post)
}

func handleGetPosts(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: handleGetPosts")
	// 이제 여기서 DB 를 긁어와서 주는 형식으로 만들어야 되지 않을까?

	var postView PostView
	var err error

	postView.Posts, err = Database.SelectAll()
	if err != nil {
		log.Println("Error Occurred", err)
	}

	//json.NewEncoder(w).Encode(posts)
	fmt.Println(postView)
	boardMain.Execute(w, postView)
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Welcome to the HomePage!")
	log.Println("Endpoint Hit: homePage")
	mainPage.Execute(w, nil)
}
