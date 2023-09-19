// Online Go compiler to run Golang program online
// Print "Hello World!" message

package main
import (
    "fmt"
    "time"
    "google.golang.org/protobuf/types/known/timestamppb"
)



//2023-06-07T22:09:49-03:00

func main() {

    t, err := time.Parse(time.RFC3339, "2023-06-07T22:09:49-03:00")
    if err != nil {
        panic(err)
    }
    
    pb := timestamppb.New(t).AsTime()
    fmt.Println(pb)
}
