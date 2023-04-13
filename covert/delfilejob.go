package convert

import (
	"fmt"
	"os"
)

func DelTrashFile(fileName []string) {
	for i := 0; i < len(fileName); i++ {
		os.Remove(fileName[i])
		fmt.Println("DelTrashFile :", fileName[i])
	}
}
