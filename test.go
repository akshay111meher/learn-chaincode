
package main

import (
	"encoding/base64"
	"fmt"
)

func main()  {
	S:= "aaaaaaaaaaaaaaa................."
	owner, _ := base64.StdEncoding.DecodeString(S)
	fmt.Println(S)
	fmt.Println(owner)
	// fmt.Println(base64.StdEncoding)
}
