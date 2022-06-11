package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//Modifique la estructura como publica
type Data struct {
	Rows [][]string
}

//Agregue un delimitador porque en algunos casos los datasets estan delimitados por comas
func Parser(pathArchivo string, delimiter string) (p *Data, err error) {
	p = new(Data)
	archivo, err := os.Open(pathArchivo)
	fileScanner := bufio.NewScanner(archivo)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		p.Rows = append(p.Rows, strings.Split(fileScanner.Text(), delimiter))
	}
	defer archivo.Close()
	return p, err
}

func Print(p *Data) {
	for _, row := range p.Rows {
		fmt.Println("Fila:", row)
	}
}
