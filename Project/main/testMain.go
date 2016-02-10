// testMain
package main
import "../test"
import "fmt"

func main(){
  fmt.Printf("Test: %d",test.Test())
  //fmt.Printf("test: %d",test.test())
  //fmt.Printf("private: %d",test.aName)
  fmt.Printf("public: %d",test.BigBro)
}
