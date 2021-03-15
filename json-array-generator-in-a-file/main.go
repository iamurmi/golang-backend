package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "strconv"
)

type emp struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Age    int64  `json:"age"`
    Salary int64  `json:"salary"`
}

func main() {
    var list []emp
    id := 0
    for i := 0; i < 10; i++ {
        id = id + i
        name := strconv.Itoa(i)
        e1 := emp{
            ID:     strconv.Itoa(id),
            Name:   name + "Urmila Kewat",
            Age:    23,
            Salary: 28000,
        }
        list = append(list, e1)
    }
    e, err := json.MarshalIndent(list, "", "")
    if err != nil {
        fmt.Print("ERROR")
    }
	
 _ = ioutil.WriteFile("test.json", e, 0644)
 fmt.Print(string(e))
	
}
