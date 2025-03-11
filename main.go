package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	NAME_KEY = "name"
	DIR_KEY = "cd"
	PATH_JSON = "./folder.json"
	CREATE_IDENTIFIER = "new"
	COMMIT_IDENTIFIER = "commit"
	PULL_IDENTIFIER = "pull"
	QUANT_MIN_ARGS_DO = 2 
	QUANT_MIN_ARGS_NEW = 3 
	QUANT_MIN_PULL = 2
)

type Args struct{
	Name string
	Branch string
	CommitName string
	CommitDesc string 
	Files string 
}

func Init(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.WriteFile(path, []byte("[]"), 0755)
			if err != nil {
				return errors.New("erro ao escrever o arquivo")
			}
		}
	}

	return nil 
}

func create(args []string, pathJson string) {
	if len(args) < QUANT_MIN_ARGS_NEW {
		fmt.Println("Não passado argumentos suficientes")
		return 
	}

	id := args[1]
	path := args[2]

	err := NewCDCommand(path, id, pathJson)
	if err == nil {
		fmt.Println("Criado com sucesso!!!")
	}
}

func doCommit(args []string, description string, files string, branch string, pathJson string) {
	if len(args) < QUANT_MIN_ARGS_DO {
		fmt.Println("Não passado argumentos suficientes")
		return 
	}

	id := args[1]
	nameCommit := args[2]

	data := Args{
		Name: id,
		CommitName: nameCommit,
		CommitDesc: description,
		Files: files,
		Branch: branch,
	}

	err := data.RealizeFile(pathJson)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Comitado com sucesso!!!")
}

func doPull(args []string, pathAbs string) {
	if len(args) < QUANT_MIN_PULL {
		fmt.Println("Não passado argumentos suficientes")
		return 
	}

	id := args[1]
	path, err := SearchDirectory(id, pathAbs, NAME_KEY, DIR_KEY)
	if err != nil {
		return
	}
	err = Pull(path)
	if err != nil {
		return
	}
	fmt.Println("Feito o pull")
}

func getBasePath(pathLocal string) (string, error){
	exePath, err := os.Executable()
	if err != nil {
		err = errors.New("erro ao obter o caminho do executável")
		fmt.Println(err.Error())
		return "", err 
	}

	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, pathLocal), nil 
}

func main() {
	pathAbs, err := getBasePath(PATH_JSON)
	if err != nil {
		return 
	}

	err = Init(pathAbs)
	if err != nil {
		return 
	}
	description := flag.String("desc", "", "Define a descrição do arquivo")
	branch := flag.String("branch", "", "Define o branch em que será comitado")
	files := flag.String("files", "", "Define os arquivos que serão comitados")

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Não passado argumentos suficientes")
		return 
	}

	whatIDo := args[0]
	switch(whatIDo) {
	case CREATE_IDENTIFIER:
		create(args, pathAbs)
	case COMMIT_IDENTIFIER:
		doCommit(args, *description, *files, *branch, pathAbs)
	case PULL_IDENTIFIER:
		doPull(args, pathAbs)
	default:
		fmt.Println("O comando não existe")
	}
}

func (obj Args) RealizeFile(pathJson string) error{
	path, err := SearchDirectory(obj.Name, pathJson, NAME_KEY, DIR_KEY)
	if err != nil {
		return err 
	}

	return CommitAndPush(obj.CommitName, obj.CommitDesc, obj.Branch, path, obj.Files)
}