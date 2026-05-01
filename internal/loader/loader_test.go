package loader


import (
    "fmt"
    "testing"
)

func TestDebug_LoadFile(t *testing.T) {
    records, errs := LoadFile("/home/mobina-mousavi/Go_project/users.json") 

    for {
        select {
        case r, ok := <-records:
            if !ok {
                records = nil
                continue
            }
            fmt.Printf("Record: %+v\n", r)
        case e, ok := <-errs:
            if !ok {
                errs = nil
                continue
            }
            fmt.Printf("Error: %v\n", e)
        }
        if records == nil && errs == nil {
            break
        }
    }
}
