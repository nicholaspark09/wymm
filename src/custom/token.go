package custom

import (
    "encoding/base64"
    "crypto/rand"
    "fmt"
 )

type Token struct{
	ThirtyTwoBit string
	SixtyFourBit string
}


func (t *Token) CreateThirtyTwoBit(){
	size := 32 // change the length of the generated random string here

   rb := make([]byte,size)
   _, err := rand.Read(rb)


   if err != nil {
      fmt.Println(err)
   }

   t.ThirtyTwoBit = base64.URLEncoding.EncodeToString(rb)
}

func (t *Token) CreateSixtyFourBit(){

}

func (t *Token) CreateToken(size int) string{
	//Create random token with size
   rb := make([]byte,size)
   _, err := rand.Read(rb)


   if err != nil {
      fmt.Println(err)
   }

   rs := base64.URLEncoding.EncodeToString(rb)
   return rs
}


/*

 func main() {
   size := 32 // change the length of the generated random string here

   rb := make([]byte,size)
   _, err := rand.Read(rb)


   if err != nil {
      fmt.Println(err)
   }

   rs := base64.URLEncoding.EncodeToString(rb)

   fmt.Println(rs)
 }
 */