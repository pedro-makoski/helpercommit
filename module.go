package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	EMPTY_STRING = ""
	CONTAINED_IN_PULL_ERROR = "Updates were rejected"
)

func WriteFile(path string, text []byte) error {
	err := os.WriteFile(path, text, 0600)
	if err != nil {
		return errors.New("erro ao escrever o arquivo")
	}

	return err
}

func Dessirealizar(object []map[string]any) ([]byte, error) {
	res, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return nil, errors.New("erro ao dessirealizar o arquivo")
	}

	return res, nil
}

func DessirealizarAndWriteFile(object []map[string]any, path string) error{
	res, err := Dessirealizar(object)
	if err != nil {
		return err 
	}

	err = WriteFile(path, res)
	return err 
}

func ReadFileInOBJ(path string) ([]map[string]any, error) {
	res, err := os.ReadFile(path);

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("erro ao ler o arquivo")
	}

	var resAsObj []map[string]any

	err = json.Unmarshal(res, &resAsObj)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("erro ao ler o arquivo")
	}

	return resAsObj, nil
}

func RunCommand(nameCommand string, args []string, errorCase string, dir string) (error, string) {
	cmd := exec.Command(nameCommand, args...);
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		aditError := errors.New(errorCase)
		fmt.Printf("%v:%v\n", aditError.Error(), string(output))
		return aditError, string(output)
	}

	return nil, string(output)
}

func Add(dir string, paths string) error {
	path := paths
	if path == EMPTY_STRING {
		path = "."
	}
	err, _ := RunCommand("git", []string{"add", path}, "erro ao adicionar", dir)
	return err 
}

func Commit(name string, description string, dir string) error {
	commands := []string{"commit", "-m", name}
	if description != EMPTY_STRING {
		commands = append(commands, "-m", description)
	}
	err, _ := RunCommand("git", commands, "erro ao fazer o commit", dir)
	return err 
}

func MoveToBrench(branch string, dir string) error {
	if branch == EMPTY_STRING {
		return nil 
	}

	err, _ := RunCommand("git", []string{"checkout", branch}, "erro ao trocar de branch", dir)
	return err 
}

func Pull(dir string) error{
	err, _ := RunCommand("git", []string{"pull", "origin"}, "erro ao fazer o pull", dir)
	return err 
}

func PushIfNecessariPull(dir string) error{
	err, output := RunCommand("git", []string{"push", "-u", "origin"}, "erro ao fazer o push", dir)
	if err != nil {
		if(strings.Contains(output, CONTAINED_IN_PULL_ERROR)) {
			err := Pull(dir)
			if err != nil {
				return err 
			}
			err = PushIfNecessariPull(dir)
			if err != nil {
				return err 
			}
		}

		return err 
	}

	return nil 
}

func CommitAndPush(commitName string, commitDescriptionOptional string, branch string, dir string, pathsFile string) error {
	err := Add(dir, pathsFile);
	if err != nil {
		return err
	}
	
	err = MoveToBrench(branch, dir)
	if err != nil {
		return err
	}

	err = Commit(commitName, commitDescriptionOptional, dir);
	if err != nil {
		return err 
	}

	err = PushIfNecessariPull(dir)
	if err != nil {
		return err 
	}

	return nil 
}

func NewCDCommand(absDirectory string, identifier string, pathJSON string) error {
	resp, err := ReadFileInOBJ(pathJSON);
	if err != nil {
		fmt.Println(err)
		return err 
	}

	resp = append(resp, map[string]any{"cd":absDirectory, "name":identifier})
	err = DessirealizarAndWriteFile(resp, pathJSON)
	if err != nil {
		fmt.Println(err)
		return err 
	}

	return nil 
}

func SearchDirectory(name string, path string, nameKey string, dirKey string) (string, error) {
	resp, err := ReadFileInOBJ(path)
	if err != nil {
		return "", err 
	}

	errNotFound := errors.New("erro ao encontar o nome/diretorio")
	for _, directoryObj := range resp {
		nameNow, ok := directoryObj[nameKey]
		if !ok {
			fmt.Println(errNotFound.Error())
			return "", errNotFound 
		}

		if nameNow == name {
			directory, ok := directoryObj[dirKey]
			if !ok {
				fmt.Println(errNotFound.Error())
				return "", errNotFound
			}
			
			return fmt.Sprintf("%v", directory), nil 
		}
	}

	fmt.Println(errNotFound.Error())
	return "", errNotFound
}