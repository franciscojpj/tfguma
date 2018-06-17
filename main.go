package main

import (
  "fmt"
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/rs/cors" 
  "container/list"
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "io"
  "encoding/base64"
  _"database/sql"
  _"github.com/lib/pq"
)

var key []byte

type User struct{
  Name string
}

type Person struct{
  Name string
  Surname string
  Birth string
  Webpage string
  Orcid string
  Id string
}

type Relacion struct{
  Name string
  Surname string
  Id int
}

type GenealogyRelationship struct{
  Name string
  Surname string
  IdPadre int
  IdHijo int
}

type Keyws struct{
  Word string
  Thesis int
}

type Alumns struct{
  RelAlumns list.List
}

type Thesis struct{
  Id string
  Title string
  Abstract string
  DefenseDate string
  Institution string
  Url string
  Department string
  IdInstitucion int
}

type Institucion struct{
  Id string
  Nombre string
  Url string
}

type ResultadoLista struct{
  IdDoctor string
  Nombre string
  Apellidos string
  IdTesis string
  Titulo string
}

type TesisEditado struct{
  IdOrig int
  IdEdit int
  Titulo string
  Nombre string
  Apellidos string
  Fecha string
  Institucion string
  IdInstitucion int
}

type UsuarioEditado struct{
  IdOrig int
  IdEdit int
  Nombre string
  Apellidos string
  Fecha string
  Institucion string
}


func main(){
  key = []byte("KeyDeMiTFG123456")
  originalText := "encrypt this golang"
  fmt.Println(originalText)

  // encrypt value to base64
  cryptoText := Encrypt(originalText)
  fmt.Println(cryptoText)

   // encrypt base64 crypto to original value
  text := Decrypt(cryptoText)
  fmt.Println(text)

  create()
  r := httprouter.New()
  Router(r)

  fmt.Println("Starting server on :8000")
  
  handler := cors.Default().Handler(r)
  http.ListenAndServe(":8000", handler)
}

// encrypt string to base64 crypto using AES
func Encrypt(text string) string {
  // key := []byte(keyText)
  plaintext := []byte(text)

  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }

  // The IV needs to be unique, but not secure. Therefore it's common to
  // include it at the beginning of the ciphertext.
  ciphertext := make([]byte, aes.BlockSize+len(plaintext))
  iv := ciphertext[:aes.BlockSize]
  if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    panic(err)
  }

  stream := cipher.NewCFBEncrypter(block, iv)
  stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

  // convert to base64
  return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func Decrypt(cryptoText string) string {
  ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }

  // The IV needs to be unique, but not secure. Therefore it's common to
  // include it at the beginning of the ciphertext.
  if len(ciphertext) < aes.BlockSize {
    panic("ciphertext too short")
  }
  iv := ciphertext[:aes.BlockSize]
  ciphertext = ciphertext[aes.BlockSize:]

  stream := cipher.NewCFBDecrypter(block, iv)

  // XORKeyStream can work in-place if the two arguments are the same.
  stream.XORKeyStream(ciphertext, ciphertext)

  return fmt.Sprintf("%s", ciphertext)
  }