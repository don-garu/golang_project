package Database

import (
	"database/sql"
	"fmt"
	_"github.com/go-sql-driver/mysql"
)

type Post struct {
	Id int
	Title string
	Text string
}

type User struct {
	Id int
	Key int
	UserId string
	UserPwd string
}

func (p Post) ToString() string {
	return fmt.Sprintf("ID: %d\nTitle: %s\nText: %s\n", p.Id, p.Title, p.Text)
}

var db *sql.DB

func Initialize() error {
	var err error
	db, err = sql.Open("mysql", "test_user:9637@tcp(127.0.0.1:3306)/test")

	return err
}

func Close() {
	db.Close()
}

func InsertUser(u User) error {
	// ToDo : Encryption
	hashPwd, err := GetHashPassword(u.UserPwd)
	if err != nil {
		return err
	}
	insert, err := db.Query("INSERT INTO users (user_id, user_pwd) values (?, ?)", u.UserId, hashPwd)
	if err != nil {
		return err
	}
	defer insert.Close()
	return nil
}

func GetUser(u User) (int, error) {
	results, err := db.Query("SELECT _id, user_pwd FROM users WHERE user_id = ?", u.UserId)
	if err != nil {
		return 0, err
	}
	var hashPwd string
	var id int
	results.Next()
	err = results.Scan(&id, &hashPwd)
	if err != nil {
		return 0, err
	}
	if CheckHashPassword(u.UserPwd, hashPwd) {
		return id, nil
	}
	return -1, nil
}

func InsertPost(p Post) error {
	insert, err := db.Query("INSERT INTO posts (title, text) VALUES(?, ?)", p.Title, p.Text)
	if err != nil {
		return err
	}
	defer insert.Close()
	return nil
}

func DeletePost(id string) error {
	_, err := db.Query("UPDATE posts SET blocked = 1 WHERE _id = ?", id)
	return err
}

func SelectByID(id string) (Post, error) {
	var post Post
	results, err := db.Query("SELECT _id, title, text from posts where _id = ?", id)
	if err != nil {
		return Post{}, err
	}
	for results.Next() {
		err = results.Scan(&post.Id, &post.Title, &post.Text)
		if err != nil {
			return Post{}, err
		}
	}
	return post, nil
}

func SelectAll() ([]Post, error) {
	var posts []Post
	posts = []Post{}
	results, err := db.Query("SELECT _id, title, text from posts where blocked = 0")
	if err != nil {
		return []Post{}, err
	}

	for results.Next() {
		var post Post

		err = results.Scan(&post.Id, &post.Title, &post.Text)
		if err != nil {
			return []Post{}, err
		}

		posts = append(posts, post)
		println(post.Id, post.Title, post.Text)
		// log.Printf(post.Title, post.Text)
	}
	return posts, nil
}
