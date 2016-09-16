package main

import (
  "os"
  "os/user"
  "bufio"
  "fmt"
  "log"
  "io"
  "net/http"
  "strings"
  "sync"
)

func check (e error){
  if e != nil {
      log.Fatal(e)
  }
}

func download(url string, wg *sync.WaitGroup){
  defer wg.Done()
  usr, err := user.Current()
  check(err)
  folder := usr.HomeDir + "/pennyArcadeComics/"
  // Creates folder if its not there
  if _, err := os.Stat(folder); os.IsNotExist(err) {
     os.Mkdir(folder, 0777)
  }

  response, err := http.Get(url)
    check(err)
    defer response.Body.Close()
    urlSplit := strings.Split(url, "/")
    title := urlSplit[len(urlSplit) -1]
    file, err := os.Create(folder + "/" + title)
    check(err)
    _, err = io.Copy(file, response.Body)
    check(err)
    file.Close()
    fmt.Println(title + " downloaded")
}

func readLine(scanner *bufio.Scanner) string{

  success := scanner.Scan()
    if success == false {
        // False on error or EOF. Check error
        err := scanner.Err()
        if err == nil {
            log.Println("Scan completed and reached EOF")
        } else {
            log.Fatal(err)
        }
    }

    url := scanner.Text()
    return url
}

func getComicLinks (scanner *bufio.Scanner) []string{
  var list []string
  for scanner.Scan() {
    list = append(list, scanner.Text())
  }
  return list
}

func createScanner (path string) *bufio.Scanner{
  file, err := os.Open(path)
  check(err)
  scanner := bufio.NewScanner(file)
  return scanner
}

func runTaks (list []string){
  var wg sync.WaitGroup
  maxConcurrency := 50
  if (len(list) < maxConcurrency) {
    maxConcurrency = len(list)
  }
  chunk := append(list[:0], list[len(list) - maxConcurrency:]...)
  list = append(list[:0], list[maxConcurrency:]...)
  wg.Add(len(chunk));
  for _, url := range chunk {
    go download(url, &wg)
  }
  wg.Wait()
  fmt.Println("Done waiting")

  if(len(list) > 0){
    fmt.Println("Number left")
    fmt.Println("%+v", len(list))
    runTaks(list)
  }

}

func main() {
  fmt.Printf("penny arcade downloader\n")
  scanner := createScanner("comicLinks1.txt")
  list := getComicLinks(scanner)
  runTaks(list)
  fmt.Println("done")
}
