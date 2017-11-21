package collector

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func TestA(t *testing.T) {

	repo, _ := git.PlainOpen("/home/wookesh/GoProjects/src/github.com/wookesh/distributed")
	iter, _ := repo.CommitObjects()
	iter.ForEach(func(commit *object.Commit) error {
		fmt.Println("##################################")
		fmt.Println(commit.Hash)
		fmt.Println("##################################")
		files, _ := commit.Files()
		files.ForEach(func(f *object.File) error {
			if strings.HasSuffix(f.Name, ".go") {
				fmt.Println("==================================")
				fmt.Println(f.Name)
				fmt.Println("==================================")
				rd, _ := f.Blob.Reader()
				reader := bufio.NewReader(rd)
				for {
					line, _, err := reader.ReadLine()
					if err == io.EOF {
						break
					}
					fmt.Println(string(line))
				}
			}
			return nil
		})
		return nil
	})

}
