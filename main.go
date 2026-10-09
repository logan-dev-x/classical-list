package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Music struct {
	id int
	title,
	compositor,
	interpreter string
}

func newMusic(title, compositor, interpreter string) Music {
	return Music{title: title, compositor: compositor, interpreter: interpreter}
}

func main() {
	db := setupDB()
	if len(os.Args) < 2 {
		fmt.Println("Provide a valid option\n")
		fmt.Println("-list\n-add\n-sh (search)\n")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "list":
		listAll(db)
	case "add":
		add(db)
	case "sh":
		printMusic(search())
	default:
		invalidAction()
	}
}

func printMusic(music Music) {
	fmt.Printf("ID: %d\nTitle: %s\nCompositor: %s\nInterpreter: %s\n\n",
		music.id, music.title, music.compositor, music.interpreter)
}

func search() Music {
	prop := os.Args[3]
	switch os.Args[2] {
	case "t":
		return findByTitle(prop)
	case "c":
		return findByCompositor(prop)
	case "i":
		return findByInterpreter(prop)
	default:
		fmt.Println("invalid prop")
		os.Exit(1)
		return Music{}
	}
}

func findByInterpreter(interpreter string) Music {
	panic("unimplemented")
}

func findByCompositor(compositor string) Music {
	panic("unimplemented")
}

func findByTitle(title string) Music {
	panic("unimplemented")
}

func add(db *sql.DB) {
	title := getProp("Title")
	compositor := getProp("Compositor")
	interpreter := getProp("Interpreter")

	err := save(
		newMusic(title, compositor, interpreter),
		db,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func save(music Music, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO musics (title, compositor, interpreter) VALUES (?, ?, ?);",
		music.title, music.compositor, music.interpreter)
	return err
}

func getProp(prop string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s: ", prop)
	for scanner.Scan() {
		return scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return ""
}

func listAll(db *sql.DB) {
	rows, err := db.Query("SELECT id, title, compositor, interpreter FROM musics;")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for rows.Next() {
		var m Music
		err = rows.Scan(&m.id, &m.title, &m.compositor, &m.interpreter)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		printMusic(m)
	}
}

func setupDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./musics.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS musics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		compositor TEXT,
		interpreter TEXT);`)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return db
}

func invalidAction() {
	fmt.Println("invalid action")
	os.Exit(1)
}
