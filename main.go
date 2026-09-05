package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: create-phantomdb-app <project-name>")
		os.Exit(1)
	}
	name := os.Args[1]

	if _, err := os.Stat(name); err == nil {
		fmt.Printf("Error: a folder named %q already exists\n", name)
		os.Exit(1)
	}

	os.MkdirAll(name+"/internal/store", 0755)
	os.MkdirAll(name+"/cmd/server", 0755)

	writeFile(name+"/go.mod", fmt.Sprintf("module %s\n\ngo 1.22\n\nrequire github.com/TobiGrant-byte/phantomdb latest\n", name))
	writeFile(name+"/internal/store/store.go", storeTemplate)
	writeFile(name+"/cmd/server/main.go", fmt.Sprintf(mainTemplate, name))
	writeFile(name+"/.gitignore", gitignoreTemplate)

	fmt.Printf("\n  Created %s\n\n", name)
	fmt.Println("  Next steps:")
	fmt.Printf("    cd %s\n", name)
	fmt.Println("    go mod tidy")
	fmt.Println("    go run cmd/server/main.go")
	fmt.Println()

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = name
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("  Note: go mod tidy failed, run it manually:")
		fmt.Println(string(out))
	}
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Println("error writing", path, ":", err)
	}
}

const storeTemplate = `package store

import "github.com/TobiGrant-byte/phantomdb"

type Store struct {
	DB *phantomdb.DB
}

func Open(path string) (*Store, error) {
	db, err := phantomdb.Open(path)
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) Close() error {
	return s.DB.Close()
}
`

const mainTemplate = `package main

import (
	"fmt"
	"log"

	"%s/internal/store"
)

func main() {
	s, err := store.Open("data")
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	fmt.Println("PhantomDB app started — data.pdb ready")

	if err := s.DB.Put([]byte("hello"), []byte("world")); err != nil {
		log.Fatal(err)
	}
	val, err := s.DB.Get([]byte("hello"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stored and retrieved: %%s\n", val)
}
`

const gitignoreTemplate = `*.pdb
*.wal
*.exe
`