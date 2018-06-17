package main

import(
	"github.com/julienschmidt/httprouter"
  "html/template"
  "path"
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
  "time"
  "encoding/json"
  "strconv"
  "strings"
)

type UserP struct {
  Name    string
  Email string
}

type Respuesta struct{
  Valor int
}

func CrearUsuario(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
    nombre := r.FormValue("nombre")
    apellidos := r.FormValue("apellidos")
    email := r.FormValue("email")
    pass := r.FormValue("pass")

    crearUsuario(nombre,apellidos,email,pass)
}

func LoginHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fp := path.Join("templates","index.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"jiji"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}


func OrcidHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
    keys, ok := r.URL.Query()["code"]
    
    if !ok || len(keys) < 1 {
        fmt.Println("Url Param 'key' is missing")
        fmt.Println(keys)
        return
    }
    //urlPag := "https://sandbox.orcid.org/oauth/token"

    key := keys[0]
    
    form := url.Values{}
    form.Add("client_id", "APP-DI48OJOVPURZQ7N2")
    form.Add("client_secret", "d7f4b0b1-2a95-41f0-a381-61a424054093")
    form.Add("grant_type", "authorization_code")
    form.Add("code", key)
    form.Add("redirect_uri", "http://localhost:8000/orcid")
    //req, err := http.NewRequest("POST", urlPag, strings.NewReader(form.Encode()))
    //req, err := http.NewRequest("POST", urlPag, bytes.NewBufferString(form.Encode()))
    req, err := http.PostForm("https://sandbox.orcid.org/oauth/token",url.Values{"client_id": {"APP-DI48OJOVPURZQ7N2"}, "client_secret": {"d7f4b0b1-2a95-41f0-a381-61a424054093"},
      "grant_type": {"authorization_code"}, "code": {key}, "redirect_uri": {"http://localhost:8000/orcid"}})

    //client := &http.Client{}
    //req.Header.Add("Accept","application/json")
    fmt.Println("Request: ",req)
    fmt.Println("Request: ",req.Body)
    bodyBytes,_ := ioutil.ReadAll(req.Body)
    bodyString := string(bodyBytes)
    fmt.Println("Response: ",bodyString)
    //resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer req.Body.Close()
    var msg User
    msg= User{bodyString}


      fp := path.Join("templates","search2.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,msg); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }

}

func Orcid2Handler(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
    fmt.Println("hace OrcidHandler")
    
    fp := path.Join("templates","confirm.html")
    tmpl, err := template.ParseFiles(fp)
    if err != nil{
      http.Error(rw , err.Error(), http.StatusInternalServerError)
      return
    }
    if err := tmpl.Execute(rw,"jiji"); err != nil {
      http.Error(rw, err.Error(), http.StatusInternalServerError)
    }
}

func InternLoginHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fmt.Println("Entra")
  name := r.FormValue("name")
  pass := r.FormValue("password")

  fmt.Println(name)
  fmt.Println(pass)

  sig := verifyUser(name,pass)
  fmt.Println(sig)

  fp := path.Join("templates","visual.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"jiji"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}

func GetIdCookie(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fmt.Println("Get idsfsfs")
  var cookie,err = r.Cookie("__prueba")
  fmt.Println("Get id cookie", cookie.Value)
  value := Decrypt(cookie.Value)
  js, err := json.Marshal(value)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
  }

func VisualizacionGrafo(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fp := path.Join("templates","visual.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"abrir"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}

func VisualizacionThesis(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fp := path.Join("templates","thesis.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"abrir"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}

func VisualizacionDoctor(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fp := path.Join("templates","doctor.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"abrir"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}

func VisualizacionSubmit(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fp := path.Join("templates","submit.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"abrir"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}

func VisualizacionEdit(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fp := path.Join("templates","edit.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"abrir"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}

func VisualizacionInstitucion(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fp := path.Join("templates","institucion.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"abrir"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}

func VisualizacionFAQ(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fp := path.Join("templates","faq.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"abrir"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}


func VisualizacionSearch(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fp := path.Join("templates","search.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"abrir"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}

func VisualizacionAdmin(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  fp := path.Join("templates","panel.html")
  tmpl, err := template.ParseFiles(fp)
  if err != nil{
    http.Error(rw , err.Error(), http.StatusInternalServerError)
    return
  }
  if err := tmpl.Execute(rw,"abrir"); err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
  }
}


func GetUserInfo(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]
  usuario := getUser(id)

  fmt.Println(usuario)
  js, err := json.Marshal(usuario)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetUserByTesis(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]
  usuario := getUserByTesis(id)

  js, err := json.Marshal(usuario)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func FindUsers(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  params,_ := r.URL.Query()["nombre"]
  nombre := params[0]

  params,_ = r.URL.Query()["apellidos"]
  apellidos := params[0]

  encuentros := findUsuarios(nombre,apellidos)
  longi := encuentros.Len()
  left := make([]Person, longi)

  sum := 0
  e := encuentros.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(Person)  
    sum++
    e = e.Next()
  }

  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetNewUsuarios(rw http.ResponseWriter, r *http.Request, p httprouter.Params){

  encuentros := findNewUsuarios()
  longi := encuentros.Len()
  left := make([]UsuarioEditado, longi)

  sum := 0
  e := encuentros.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(UsuarioEditado)  
    sum++
    e = e.Next()
  }

  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetDelUsuarios(rw http.ResponseWriter, r *http.Request, p httprouter.Params){

  encuentros := findDelUsuarios()
  longi := encuentros.Len()
  left := make([]UsuarioEditado, longi)

  sum := 0
  e := encuentros.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(UsuarioEditado)  
    sum++
    e = e.Next()
  }

  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func FindResults(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  params,_ := r.URL.Query()["titulo"]
  titulo := params[0]

  params,_ = r.URL.Query()["nombre"]
  nombre := params[0]

  params,_ = r.URL.Query()["apellidos"]
  apellidos := params[0]

  params,_ = r.URL.Query()["orcid"]
  orcid := params[0]

  params,_ = r.URL.Query()["institucion"]
  institucion := params[0]
  idIns,_ := strconv.Atoi(institucion)

  resultados := findResultados(titulo,nombre,apellidos,orcid,idIns)
  longi := resultados.Len()
  left := make([]ResultadoLista, longi)

  sum := 0
  e := resultados.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(ResultadoLista)  
    sum++
    e = e.Next()
  }

  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func FindInstituciones(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  
  params,_ := r.URL.Query()["nombre"]
  nombre := params[0]

  encuentros := findInst(nombre)
  longi := encuentros.Len()
  left := make([]Institucion, longi)

  sum := 0
  e := encuentros.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(Institucion)  
    sum++
    e = e.Next()
  }

  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetAllInstitutions(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  encuentros := getInsts()
  longi := encuentros.Len()
  left := make([]Institucion, longi)

  sum := 0
  e := encuentros.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(Institucion)  
    sum++
    e = e.Next()
  }

  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetSupervisors(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]

  relaciones := getSups(id)
  longi := relaciones.Len()
  left := make([]Relacion, longi)

  sum := 0
  e := relaciones.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(Relacion)  
    sum++
    e = e.Next()
  }
  
  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetDireccion(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]

  relaciones := getDirec(id)
  longi := relaciones.Len()
  left := make([]Relacion, longi)

  sum := 0
  e := relaciones.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(Relacion)  
    sum++
    e = e.Next()
  }
  
  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetJurado(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]

  relaciones := getJur(id)
  longi := relaciones.Len()
  left := make([]Relacion, longi)

  sum := 0
  e := relaciones.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(Relacion)  
    sum++
    e = e.Next()
  }
  
  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetKeywords(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]

  palabras := getKeys(id)
  longi := palabras.Len()
  left := make([]Keyws, longi)

  sum := 0
  e := palabras.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(Keyws)  
    sum++
    e = e.Next()
  }
  
  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}


func GetSons(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]

  keys,_ = r.URL.Query()["depth"]
  depth := keys[0]

  idR,_ := strconv.Atoi(id)
  ndepth,_ := strconv.Atoi(depth)

  relaciones := getSons(idR,ndepth)
  longi := relaciones.Len()
  left := make([]GenealogyRelationship, longi)

  sum := 0
  e := relaciones.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(GenealogyRelationship)  
    sum++
    e = e.Next()
  }
  
  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetFathers(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]

  keys,_ = r.URL.Query()["depth"]
  depth := keys[0]

  idR,_ := strconv.Atoi(id)
  ndepth,_ := strconv.Atoi(depth)

  relaciones := getFathers(idR,ndepth)
  longi := relaciones.Len()
  left := make([]GenealogyRelationship, longi)

  sum := 0
  e := relaciones.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(GenealogyRelationship)  
    sum++
    e = e.Next()
  }
  
  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func ReadPost(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  name := r.FormValue("user")
  pass := r.FormValue("pass")
  sig := verifyUser(name,pass)
  resultado := Respuesta{sig}

  expiration := time.Now().Add(365 * 24 * time.Hour)

  valor := Encrypt(strconv.Itoa(sig))
  //cookie := http.Cookie{Name: "__prueba", Value: strconv.Itoa(sig), Expires: expiration}
  cookie := http.Cookie{Name: "__prueba", Value: valor, Expires: expiration}
  http.SetCookie(rw, &cookie)

  js, err := json.Marshal(resultado)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func WriteCookie(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  id := r.FormValue("id")

  expiration := time.Now().Add(365 * 24 * time.Hour)

  valor := Encrypt(id)
  //cookie := http.Cookie{Name: "__prueba", Value: strconv.Itoa(sig), Expires: expiration}
  cookie := http.Cookie{Name: "__prueba", Value: valor, Expires: expiration}
  http.SetCookie(rw, &cookie)

  fmt.Println("Guarda cookie con valor", valor)
}

func CloseSession(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  resultado := Respuesta{0}

  expiration := time.Now().Add(365 * 24 * time.Hour)
  cookie := http.Cookie{Name: "__prueba", Value: strconv.Itoa(0), Expires: expiration}
  http.SetCookie(rw, &cookie)

  js, err := json.Marshal(resultado)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}



func insertarUsuario(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  nombre := r.FormValue("nombre")
  apellidos := r.FormValue("apellidos")
  orcid := r.FormValue("orcid")
  pagina := r.FormValue("pagina")


  //Cambiar el metodo para que tmb guarde quien hace la peticion
  insertUser(nombre,apellidos,orcid,pagina)
}


func InsertUserOrcid(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  nombre := r.FormValue("nombre")
  orcid := r.FormValue("orcid")


  //mira si ya está creado
  id := checkOrcid(orcid)

  if id==0 {
    id = insertUserOrcid(nombre,orcid)
  }

    js, err := json.Marshal(id)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func UpdateUsuario(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  nombre := r.FormValue("nombre")
  apellidos := r.FormValue("apellidos")
  orcid := r.FormValue("orcid")
  pagina := r.FormValue("pagina")
  fecha := r.FormValue("fecha")
  usuario := r.FormValue("id")

  idUsuario,_ := strconv.Atoi(usuario)

  //Cambiar el metodo para que tmb guarde quien hace la peticion
  updateUser(nombre,apellidos,orcid,pagina,fecha,idUsuario)
}

func AskDeleteUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]

  idUsuario,_ := strconv.Atoi(id)

  //Cambiar el metodo para que tmb guarde quien hace la peticion
  askDeleteUser(idUsuario)
}

func AskDeleteTesis(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]

  idTesis,_ := strconv.Atoi(id)

  //Cambiar el metodo para que tmb guarde quien hace la peticion
  askDeleteTesis(idTesis)
}


func UpdateTesis(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  titulo := r.FormValue("titulo")
  fecha := r.FormValue("fecha")
  url := r.FormValue("url")
  abstract := r.FormValue("abstract")
  departamento := r.FormValue("departamento")
  institucion := r.FormValue("institucion")
  id := r.FormValue("id")

  idInstitucion,_ := strconv.Atoi(institucion)
  idTesis,_ := strconv.Atoi(id)

  fmt.Println("Handler",idInstitucion)

  updateTesis(titulo,fecha,url,abstract,departamento,idInstitucion,idTesis)
}

func DeleteUsuario(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  usuario := r.FormValue("id")

  idUsuario,_ := strconv.Atoi(usuario)

  //Cambiar el metodo para que tmb guarde quien hace la peticion
 deleteUser(idUsuario)
}

func DeleteTesis(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  tesis := r.FormValue("id")

  idTesis,_ := strconv.Atoi(tesis)

  //Cambiar el metodo para que tmb guarde quien hace la peticion
 deleteTesis(idTesis)
}

func DeleteJury(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  tesis := r.FormValue("id")

  idTesis,_ := strconv.Atoi(tesis)

  
 deleteJury(idTesis)
}

func DeleteDirector(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  tesis := r.FormValue("id")

  idTesis,_ := strconv.Atoi(tesis)


  
 deleteDirector(idTesis)
}

func EditarUsuario(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  nombre := r.FormValue("nombre")
  apellidos := r.FormValue("apellidos")
  orcid := r.FormValue("orcid")
  pagina := r.FormValue("pagina")
  doctor := r.FormValue("doctor")
  nacimiento := r.FormValue("fecha")

  idDoctor,_ := strconv.Atoi(doctor)
  //Cambiar dejar el insert  user y ademas una tabla nueva que relacione usuario antiguo, el nuevo creado,  y el que ha creado la peticion
  id := insertEditedUser(nombre,apellidos,orcid,pagina,idDoctor,nacimiento)
  insertRelationEditUsers(idDoctor,id)
  // getUsersEdited(nombre,apellidos)
}

func insertarInstitucion(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  nombre := r.FormValue("nombre")
  url := r.FormValue("url")

  insertInstitucion(nombre,url)
}

func insertarTesis(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  titulo := r.FormValue("titulo")
  fecha := r.FormValue("fecha")
  url := r.FormValue("url")
  abstract := r.FormValue("abstract")
  lector := r.FormValue("lector")
  departamento := r.FormValue("departamento")
  institucion := r.FormValue("institucion")
  keywords := r.FormValue("keywords")

  palabra := strings.Split(keywords, ",")

  idLector,_ := strconv.Atoi(lector)
  idInstitucion,_ := strconv.Atoi(institucion)

  insertTesis(titulo,fecha,url,abstract,idLector,departamento,idInstitucion)

  //tesis := getTesisUser(lector)
  tesis := getTesis(idLector)

  for i := 0; i < len(palabra); i++ {
    insertKeyword(palabra[i],tesis.Id)
  }

  fmt.Println(tesis)
  js, err := json.Marshal(tesis)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func editarTesis(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  
  titulo := r.FormValue("titulo")
  fecha := r.FormValue("fecha")
  url := r.FormValue("url")
  abstract := r.FormValue("abstract")
  lector := r.FormValue("lector")
  departamento := r.FormValue("departamento")
  institucion := r.FormValue("institucion")

  idLector,_ := strconv.Atoi(lector)
  idInstitucion,_ := strconv.Atoi(institucion)

  insertTesis(titulo,fecha,url,abstract,idLector,departamento,idInstitucion)

  //tesis := getTesisUser(lector)
  tesis1, tesis2 := getModifiedTesis(idLector)
  idt1,_ := strconv.Atoi(tesis1.Id)
  idt2,_ := strconv.Atoi(tesis2.Id)

  if idt2<idt1 {
    aux := idt1
    idt1 = idt2
    idt2 = aux
  }

  insertRelationEditTesis(idt1,idt2)

  js1, err := json.Marshal(idt2)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js1)
}

func insertarDoctores(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  director := r.FormValue("director")
  codirector := r.FormValue("codirector")
  codirector2 := r.FormValue("codirector2")
  tesis := r.FormValue("tesis")

  idDir,_ :=  strconv.Atoi(director)
  idCod,_ := strconv.Atoi(codirector)
  idCod2,_ := strconv.Atoi(codirector2)
  idTesis,_ := strconv.Atoi(tesis)


  insertDirector(idDir,idTesis)
  insertDirector(idCod,idTesis)
  insertDirector(idCod2,idTesis)
}

func insertarJurado(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  j1 := r.FormValue("j1")
  j2 := r.FormValue("j2")
  j3 := r.FormValue("j3")
  j4 := r.FormValue("j4")
  j5 := r.FormValue("j5")
  tesis := r.FormValue("tesis")

  idj1,_ := strconv.Atoi(j1)
  idj2,_ := strconv.Atoi(j2)
  idj3,_ := strconv.Atoi(j3)
  idj4,_ := strconv.Atoi(j4)
  idj5,_ := strconv.Atoi(j5)
  idTesis,_ := strconv.Atoi(tesis)

  insertJurado(idj1,idTesis)
  insertJurado(idj2,idTesis)
  insertJurado(idj3,idTesis)
  insertJurado(idj4,idTesis)
  insertJurado(idj5,idTesis)
}

func GetTesisUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]
  tesis := getTesisUser(id)

  fmt.Println(tesis)
  js, err := json.Marshal(tesis)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetTesisId(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]
  tesis := getTesisById(id)

  js, err := json.Marshal(tesis)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)

}


func GetInstitucion(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  keys,_ := r.URL.Query()["id"]
  id := keys[0]
  tesis := getInstitucion(id)

  js, err := json.Marshal(tesis)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)

}


func GetNewTesis(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  editados := getNewTesis()
  longi := editados.Len()
  left := make([]TesisEditado, longi)

  sum := 0
  e := editados.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(TesisEditado)  
    sum++
    e = e.Next()
  }
  
  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetDelTesis(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  editados := getDelTesis()
  longi := editados.Len()
  left := make([]TesisEditado, longi)

  sum := 0
  e := editados.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(TesisEditado)  
    sum++
    e = e.Next()
  }
  
  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}


func GetEditedTesis(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  editados := getEditedTesis()
  longi := editados.Len()
  left := make([]TesisEditado, longi)

  sum := 0
  e := editados.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(TesisEditado)  
    sum++
    e = e.Next()
  }
  
  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}

func GetEditedUsuarios(rw http.ResponseWriter, r *http.Request, p httprouter.Params){
  editados := getEditedUsuarios()
  longi := editados.Len()
  left := make([]UsuarioEditado, longi)

  sum := 0
  e := editados.Front()
  for sum < longi && e!=nil{
    left[sum] = e.Value.(UsuarioEditado)  
    sum++
    e = e.Next()
  }
  
  js, err := json.Marshal(left)
  if err != nil {
    http.Error(rw, err.Error(), http.StatusInternalServerError)
    return
  }
  rw.Header().Set("Content-Type", "application/json")
  rw.Write(js)
}