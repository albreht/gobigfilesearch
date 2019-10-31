// test project main.go
package main

import (
	//	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var searchWord string
var searchWordBytes []byte
var searchWordBytesLength int
var wg sync.WaitGroup
var outputStrings [256]strings.Builder

func main() {

	searchWord = "ala"
	searchWordBytes = []byte(searchWord)
	searchWordBytesLength = len(searchWordBytes)

	var filePath string = "C:\\Users\\test\\source\\repos\\FastFileReader\\FastFileReader\\bin\\Debug\\test.txt"

	file, _ := os.Open(filePath)
	fileInfo, _ := file.Stat()
	fmt.Printf("Wielkość pliku %db\n", fileInfo.Size())

	partSieze := fileInfo.Size() / 256

	wg.Add(256)

	for i := 0; i < 256; i++ {
		go readFileChunk(filePath, int64(i)*partSieze, int(partSieze), i)
	}

	// reader := bufio.NewReader(os.Stdin)
	// reader.ReadLine()
	wg.Wait()

	f, err := os.Create("result123")
	check(err, 0)

	for _, value := range outputStrings {
		//fmt.Printf("%s", str)
		f.WriteString(value.String())
	}
	f.Close()
	fmt.Printf("Koniec\n")
}

func check(e error, count int64) {
	if e != nil {
		panic(e)
	}
}

func readFile(filePath string, start int64, length int64, filePart int) {

	fmt.Printf("Czytam czesc %d\n", filePart)
	var offset int64 = start

	file, err := os.Open(filePath) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 1024)

	for {

		offset++
		o2, err := file.Seek(offset*1024, 0)
		check(err, o2)
		o2 = 0
		count, err := file.Read(data)
		//check(err, (int64)count)

		//fmt.Printf("read bytes: %q\n", data)

		if count == 0 {
			break
		}
		if offset*1024 > start*1024+length {
			break
		}
	}
	//fmt.Printf("Przeczytano plik plik %s\n", filePath)

}

func readFileChunk(filename string, startPositon int64, length int, filePart int) {

	var sb strings.Builder

	file, err := os.Open(filename) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, length)

	o2, err := file.Seek(startPositon, 0)
	check(err, o2)

	separator := []byte("\r\n")
	var separatorLen int = len(separator)
	var lastSeparatorPosition int = 0

	file.Read(data)

	//fmt.Printf("%d\n", length)
	for i := 0; i < length; i++ {
		if separator[0] == data[i] {
			separatorFounded := separatorLen - 1
			for j := 1; j < separatorLen; j++ {
				separatorFounded--
			}
			if separatorFounded == 0 {
				//fmt.Printf("%d %d\n", lastSeparatorPosition, i)
				var matchedLine = findWord(data[lastSeparatorPosition:i])
				if matchedLine != nil {
					sb.Write(matchedLine)
				}

				lastSeparatorPosition = i + 1

			}
		}
	}

	defer wg.Done()
	outputStrings[filePart] = sb
	fmt.Printf("koniec przetwarzania fragmentu %d \n", filePart)

}

func findWord(data []byte) []byte {
	for i := 0; i < len(data); i++ {
		if data[i] == searchWordBytes[0] {
			var founded = searchWordBytesLength - 1
			for j := 1; j < searchWordBytesLength; j++ {

				if data[i+j] == searchWordBytes[j] {
					founded--

					if founded == 0 {
						//fmt.Printf("%s\n", string(data))
						return data
					}
				}

			}
		}
	}
	return nil
}
