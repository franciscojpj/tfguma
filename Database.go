package main

import(
  "fmt"
	"database/sql"
  "container/list"
  _"github.com/lib/pq"
)

var db *sql.DB


const (
  host     = "localhost"
  port     = 5432
  user     = "administrador"
  password = "tfg2018uma"
  dbname   = "db_genealogy"
)

func create(){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  fmt.Println("Successfully connected!")
}

func crearUsuario(nombre string,apellidos string,email string,pass string) {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("INSERT INTO person (name,surname,email,password,confirmacion) VALUES ($1,$2,$3,$4,2)")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(nombre,apellidos,email,pass)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error crear usuario")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}

func verifyUser(nameUser,passUser string) int{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  fmt.Println(passUser);
  //var name string
  marcador:=-1
  rows, err := db.Query("SELECT id FROM person where email = $1 and password = $2",nameUser, passUser)
  checkErr(err)
  for rows.Next(){
          rows.Scan(&marcador)
  }
  defer rows.Close()

  return marcador
}

func getUser(id string) Person{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()

  var usuario Person
  fmt.Println(id);
  rows:= db.QueryRow("SELECT name,surname,birth,homepage,orcid FROM person where id=$1",id)
  
  var name,surname,birth,homepage,orcid string
  rows.Scan(&name,&surname,&birth,&homepage,&orcid)
  usuario = Person{name,surname,birth,homepage,orcid,id}

  return usuario
}

func checkOrcid(orcid string) int {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()

  rows:= db.QueryRow("SELECT id FROM person where orcid=$1",orcid)
  
  id := 0
  rows.Scan(&id)


  return id
}

func getUserByTesis(id string) int{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  

  rows:= db.QueryRow("SELECT reader FROM thesis where thesis.id=$1",id)
          var reader int
          rows.Scan(&reader)
  return reader
}

func getTesisUser(id string) Thesis{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  var thesis Thesis
  fmt.Println(id);
  rows:= db.QueryRow("SELECT thesis.id, title, abstract, defensedate, thesis.url, institution.name, thesis.department, institution.id FROM public.thesis LEFT JOIN institution ON thesis.institution=institution.id WHERE reader=$1;",id)
          var tid,title,abstract,defensedate,institution,url,departamento string
          var idInst int
          rows.Scan(&tid,&title,&abstract,&defensedate,&url,&institution,&departamento,&idInst)
          thesis = Thesis{tid,title,abstract,defensedate,institution,url,departamento,idInst}
  return thesis
}

func getTesisById(id string) Thesis{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  var thesis Thesis
  fmt.Println(id);
  rows:= db.QueryRow("SELECT thesis.id, title, abstract, defensedate, thesis.url, institution.name, thesis.department, institution.id FROM public.thesis LEFT JOIN institution ON thesis.institution=institution.id WHERE thesis.id=$1;",id)
          var tid,title,abstract,defensedate,institution,url,departamento string
          var idInst int

          rows.Scan(&tid,&title,&abstract,&defensedate,&url,&institution,&departamento, &idInst)
          thesis = Thesis{tid,title,abstract,defensedate,institution,url,departamento,idInst}
  return thesis
}


func getInstitucion(id string) Institucion{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  var institucion Institucion

  rows:= db.QueryRow("SELECT name, url FROM institution WHERE id=$1;",id)
         
          var name,url string


          rows.Scan(&name,&url)
          institucion = Institucion{id,name,url}
  return institucion
}


func getNewTesis() list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  
  l := list.New()
  defer db.Close()
  var editado TesisEditado

  rows,_ := db.Query("SELECT thesis.id, title, person.name, surname, defensedate, institution.name,institution.id FROM public.thesis LEFT JOIN institution ON thesis.institution=institution.id LEFT JOIN person ON thesis.reader=person.id WHERE thesis.confirmacion=2;")
         if rows.Next(){
            var oId,eId,idInst int
            var title,name,surname,defenseDate,institutionName string
            eId=0
            rows.Scan(&oId,&title,&name,&surname,&defenseDate,&institutionName,&idInst)
            editado = TesisEditado{oId,eId,title,name,surname,defenseDate,institutionName,idInst}
            l.PushBack(editado)
         }
         
  return *l
}

func getDelTesis() list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  
  l := list.New()
  defer db.Close()
  var editado TesisEditado

  rows,_ := db.Query("SELECT thesis.id, title, person.name, surname, defensedate, institution.name,institution.id FROM public.thesis LEFT JOIN institution ON thesis.institution=institution.id LEFT JOIN person ON thesis.reader=person.id WHERE thesis.del=1;")
         if rows.Next(){
            var oId,eId,idInst int
            var title,name,surname,defenseDate,institutionName string
            eId=0
            rows.Scan(&oId,&title,&name,&surname,&defenseDate,&institutionName,&idInst)
            editado = TesisEditado{oId,eId,title,name,surname,defenseDate,institutionName,idInst}
            l.PushBack(editado)
         }
         
  return *l
}

func getEditedTesis() list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  
  l := list.New()
  defer db.Close()
  var editado TesisEditado

  rows,_ := db.Query("SELECT original, editado, title, person.name, surname, defensedate, institution.name,institution.id FROM public.reledittesis LEFT JOIN thesis ON thesis.id=editado LEFT JOIN person ON thesis.reader=person.id LEFT JOIN institution ON thesis.institution= institution.id;")
         if rows.Next(){
            var oId,eId,idInst int
            var title,name,surname,defenseDate,institutionName string
            rows.Scan(&oId,&eId,&title,&name,&surname,&defenseDate,&institutionName,&idInst)
            editado = TesisEditado{oId,eId,title,name,surname,defenseDate,institutionName,idInst}
            l.PushBack(editado)
         }
         
  return *l
}

func getEditedUsuarios() list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  
  l := list.New()
  defer db.Close()
  var editado UsuarioEditado

  rows,_ := db.Query("SELECT original, editado, person.name, surname, defensedate, institution.name FROM public.releditusers LEFT JOIN person ON person.id=editado LEFT JOIN thesis ON thesis.reader=person.id LEFT JOIN institution ON thesis.institution= institution.id;")
         if rows.Next(){
            var oId,eId int
            var name,surname,defenseDate,institutionName string
            rows.Scan(&oId,&eId,&name,&surname,&defenseDate,&institutionName)
            editado = UsuarioEditado{oId,eId,name,surname,defenseDate,institutionName}

            l.PushBack(editado)
         }
         
  return *l
}

func getTesis(id int) Thesis{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  var thesis Thesis
  fmt.Println(id,id);
  rows:= db.QueryRow("SELECT thesis.id, title, abstract, defensedate, thesis.url, institution.name, thesis.department,institution.id FROM public.thesis LEFT JOIN institution ON thesis.institution=institution.id WHERE reader=$1;",id)
          var tid,title,abstract,defensedate,institution,url,departamento string
          var idInst int
          rows.Scan(&tid,&title,&abstract,&defensedate,&url,&institution,&departamento,&idInst)
          thesis = Thesis{tid,title,abstract,defensedate,institution,url,departamento,idInst}
  return thesis
}

func getModifiedTesis(id int) (Thesis,Thesis) {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  var thesis1,thesis2 Thesis
  fmt.Println(id,id);
  rows,_ := db.Query("SELECT thesis.id, title, abstract, defensedate, thesis.url, institution.name, thesis.department,institution.id FROM public.thesis LEFT JOIN institution ON thesis.institution=institution.id WHERE reader=$1;",id)
  var tid,title,abstract,defensedate,institution,url,departamento string
  var idInst int
  if rows.Next() {
    rows.Scan(&tid,&title,&abstract,&defensedate,&url,&institution,&departamento,&idInst)
    thesis1 = Thesis{tid,title,abstract,defensedate,institution,url,departamento,idInst}
  }
  if rows.Next() {
    rows.Scan(&tid,&title,&abstract,&defensedate,&url,&institution,&departamento,&idInst)
    thesis2 = Thesis{tid,title,abstract,defensedate,institution,url,departamento,idInst}
  }

  return thesis1,thesis2
}

func findUsuarios(nombre,apellidos string) list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  

  l := list.New()
  rows,err:= db.Query("SELECT id, name, surname, birth , homepage, orcid FROM person where lower(name) LIKE '%' || $1 || '%'  AND lower(surname) LIKE  '%' || $2 || '%' AND confirmacion!=5",nombre,apellidos)
         
         if err != nil {
          panic(err)
         }
        for rows.Next() { 
          var id,name,surname,birth,homepage, orcid string
          rows.Scan(&id,&name,&surname,&birth,&homepage,&orcid)
          usuario := Person{name,surname,birth,homepage,orcid,id}
          l.PushBack(usuario)
          fmt.Println(id)
        }
        return *l
}

func findNewUsuarios() list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  
  l := list.New()
  rows,err:= db.Query("SELECT person.id, person.name, surname, thesis.defensedate , institution.name FROM person LEFT JOIN thesis ON thesis.reader = person.id LEFT JOIN institution ON institution.id=thesis.institution where person.confirmacion=2 ")
         
         if err != nil {
          panic(err)
         }
        for rows.Next() { 
          var idO, idE int
          idE =0
          var name,surname,fecha,institucion string
          rows.Scan(&idO,&name,&surname,&fecha,&institucion)
          usuario := UsuarioEditado{idO,idE,name,surname,fecha,institucion}
          l.PushBack(usuario)
        }
        return *l
}

func findDelUsuarios() list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  
  l := list.New()
  rows,err:= db.Query("SELECT person.id, person.name, surname, thesis.defensedate , institution.name FROM person LEFT JOIN thesis ON thesis.reader = person.id LEFT JOIN institution ON institution.id=thesis.institution where person.del=1 ")
         
         if err != nil {
          panic(err)
         }
        for rows.Next() { 
          var idO, idE int
          idE =0
          var name,surname,fecha,institucion string
          rows.Scan(&idO,&name,&surname,&fecha,&institucion)
          usuario := UsuarioEditado{idO,idE,name,surname,fecha,institucion}
          l.PushBack(usuario)
        }
        return *l
}

func findResultados(titulo string, nombre string,apellidos string,orcid string, institution int) list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  l := list.New()

  fmt.Println("Datos de la busqueda")
  fmt.Println(titulo,nombre,apellidos,orcid,institution)
  if institution==0 {
    rows,err:= db.Query("SELECT person.id, person.name, person.surname, thesis.id, thesis.title FROM person LEFT JOIN thesis on (person.id = thesis.reader) where (lower(thesis.title) LIKE '%' || $1 || '%') AND thesis.confirmacion=0 AND (lower(person.name) LIKE '%' || $2 || '%'  AND lower(surname) LIKE  '%' || $3 || '%') AND (orcid LIKE '%' || $4 || '%') AND person.confirmacion=0",titulo,nombre,apellidos,orcid)
    
    if err != nil {
          panic(err)
         }
        for rows.Next() { 
          var persona,name,surname,tesis,titulo string
          rows.Scan(&persona,&name,&surname,&tesis,&titulo)
          resultado := ResultadoLista{persona,name,surname,tesis,titulo}
          l.PushBack(resultado)
        }

  }else{
    rows,err:= db.Query("SELECT person.id, person.name, person.surname, thesis.id, thesis.title FROM person LEFT JOIN thesis on (person.id = thesis.reader) where (lower(thesis.title) LIKE '%' || $1 || '%') AND thesis.confirmacion=0 AND (thesis.institution = $5)  AND (lower(person.name) LIKE '%' || $2 || '%'  AND lower(surname) LIKE  '%' || $3 || '%') AND (orcid LIKE '%' || $4 || '%') AND person.confirmacion=0",titulo,nombre,apellidos,orcid,institution)
    
    if err != nil {
          panic(err)
         }
        for rows.Next() { 
          var persona,name,surname,tesis,titulo string
          rows.Scan(&persona,&name,&surname,&tesis,&titulo)
          resultado := ResultadoLista{persona,name,surname,tesis,titulo}
          l.PushBack(resultado)
        }
  }
         
        return *l
}


func findInst(nombre string) list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  
  l := list.New()
  rows,err:= db.Query("SELECT id, name, url from institution where lower(name) LIKE '%' || $1 || '%'",nombre)
         
         if err != nil {
          panic(err)
         }
        for rows.Next() { 
          var id,name,url string
          rows.Scan(&id,&name,&url)
          institucion := Institucion{id,name,url}
          l.PushBack(institucion)
        }
        return *l
}

func getInsts() list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  
  l := list.New()
  rows,err:= db.Query("SELECT id, name, url from institution")
         
         if err != nil {
          panic(err)
         }
        for rows.Next() { 
          var id,name,url string
          rows.Scan(&id,&name,&url)
          institucion := Institucion{id,name,url}
          l.PushBack(institucion)
        }
        return *l
}

func getSups(id string) list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  //var name string
  var usuario Relacion
  l := list.New()
  rows,_ := db.Query("SELECT supervisor, p.name as supname, p.surname as supsurname FROM public.reltes INNER JOIN person p ON (supervisor = p.id) WHERE reader=$1",id)

        for rows.Next() {
            var name,surname string
            var sup int
            rows.Scan(&sup,&name,&surname)
            usuario = Relacion{name,surname,sup}
            l.PushBack(usuario)
            fmt.Println(l.Front().Value)
        }
  return *l
}

func getDirec(id string) list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  //var name string
  var usuario Relacion
  l := list.New()
  rows,_ := db.Query("SELECT p.name as name, p.surname as surname, p.id as id FROM public.directorthesis INNER JOIN person p ON (director = p.id) WHERE thesis=$1",id)

        for rows.Next() {
            var name,surname string
            var sup int
            rows.Scan(&name,&surname,&sup)
            usuario = Relacion{name,surname,sup}
            l.PushBack(usuario)
            fmt.Println(l.Front().Value)
        }
  return *l
}

func getJur(id string) list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  //var name string
  var usuario Relacion
  l := list.New()
  rows,_ := db.Query("SELECT p.name as name, p.surname as surname, p.id as id FROM public.jurythesis INNER JOIN person p ON (jury = p.id) WHERE thesis=$1",id)

        for rows.Next() {
            var name,surname string
            var sup int
            rows.Scan(&name,&surname,&sup)
            usuario = Relacion{name,surname,sup}
            l.PushBack(usuario)
            fmt.Println(l.Front().Value)
        }
  return *l
}

func getKeys(id string) list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  //var name string
  var palabra Keyws
  l := list.New()
  rows,_ := db.Query("SELECT word, id FROM public.keyword WHERE id=$1",id)

        for rows.Next() {
            var word string
            var thesis int
            rows.Scan(&word,&thesis)
            palabra = Keyws{word,thesis}
            l.PushBack(palabra)
            fmt.Println(l.Front().Value)
        }
  return *l
}

func getSons(id int,depth int) list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  //var name string
  var relacion GenealogyRelationship
  l := list.New()
  rows,_ := db.Query("WITH RECURSIVE hijos AS (SELECT reader, name,surname, supervisor, 1 as depth FROM reltes INNER JOIN person on (reader = id) WHERE supervisor = $1 UNION SELECT r.reader,p.name,p.surname, r.supervisor, h.depth+1 FROM reltes r INNER JOIN person p on (r.reader = p.id) INNER JOIN hijos h ON h.reader = r.supervisor) SELECT reader, name,surname, supervisor FROM hijos WHERE depth<$2;",id,depth)

        for rows.Next() {
            var name,surname string
            var idS,idR int
            rows.Scan(&idR,&name,&surname,&idS)
            relacion = GenealogyRelationship{name,surname,idS,idR}
            l.PushBack(relacion)
            fmt.Println(l.Front().Value)
        }
  return *l
}


func getFathers(id int,depth int) list.List{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  //var name string
  var relacion GenealogyRelationship
  l := list.New()
  rows,_ := db.Query("WITH RECURSIVE hijos AS (SELECT reader, name,surname, supervisor, 1 as depth FROM reltes INNER JOIN person on (supervisor = id) WHERE reader = $1 UNION SELECT r.reader,p.name,p.surname, r.supervisor, h.depth+1 FROM reltes r INNER JOIN person p on (r.supervisor = p.id) INNER JOIN hijos h ON h.supervisor = r.reader) SELECT reader, name,surname, supervisor FROM hijos WHERE depth<$2;",id,depth)

        for rows.Next() {
            var name,surname string
            var idS,idR int
            rows.Scan(&idR,&name,&surname,&idS)
            relacion = GenealogyRelationship{name,surname,idS,idR}
            l.PushBack(relacion)
            fmt.Println(l.Front().Value)
        }
  return *l
}

func getUsersEdited(nombre string,apellidos string) (int,int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  id1:=0
  id2:=0

  defer db.Close()
  //var name string
  rows,_ := db.Query("SELECT id FROM person where name = $1 and surname = $2",nombre,apellidos)

  if rows.Next() {
    rows.Scan(&id1)
  }
  if rows.Next() {
    rows.Scan(&id2)
  }
  if rows.Next() {
    id1=0
    id2=0
  }
        
  return id1, id2
}

func insertUser(nombre string,apellidos string, orcid string, pagina string){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("INSERT INTO person (name,surname,orcid,homepage,confirmacion) VALUES ($1,$2,$3,$4,$5)")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(nombre,apellidos,orcid,pagina,2)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error insercion usuario")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}


func insertUserOrcid(nombre string,orcid string) int{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  id:=0
  apellidos := ""
  defer db.Close()
  //var name string
  err = db.QueryRow("INSERT INTO person (name,surname,orcid,confirmacion) VALUES ($1,$2,$3,0) RETURNING id",nombre,apellidos,orcid).Scan(&id)
  if err != nil {
    fmt.Println(err)
    }
  return id
}

func updateUser(nombre string,apellidos string, orcid string, pagina string,fecha string,usuario int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("UPDATE person SET name=$1 ,surname=$2 ,orcid=$3 ,homepage=$4, birth=$5, confirmacion=0 WHERE id=$6")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(nombre,apellidos,orcid,pagina,fecha,usuario)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error actualizacion usuario")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}


func askDeleteUser(usuario int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("UPDATE person SET del=1 WHERE id=$1")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(usuario)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error eliminar usuario")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}

func askDeleteTesis(tesis int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("UPDATE thesis SET del=1 WHERE reader=$1")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(tesis)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error insercion tesis")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}

func deleteUser(usuario int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("DELETE FROM person WHERE id=$1")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(usuario)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error borrar usuario")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}

func deleteTesis(tesis int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("DELETE FROM thesis WHERE id=$1")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(tesis)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error borrado tesis")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}

func deleteJury(tesis int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("DELETE FROM jurythesis WHERE thesis=$1")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(tesis)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error eliminar tribunal")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}

func deleteDirector(tesis int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("DELETE FROM directorthesis WHERE thesis=$1")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(tesis)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error eliminar direccion director")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}

func insertEditedUser(nombre string,apellidos string, orcid string, pagina string,doctor int, nacimiento string) int{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  id:=0

  defer db.Close()
  //var name string
  err = db.QueryRow("INSERT INTO person (name,surname,orcid,homepage,confirmacion,birth) VALUES ($1,$2,$3,$4,1,$5) RETURNING id",nombre,apellidos,orcid,pagina,nacimiento).Scan(&id)
  if err != nil {
    fmt.Println(err)
    }
  return id
}

func insertRelationEditUsers(doctor int,editado int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("INSERT INTO releditusers (original,editado) VALUES ($1,$2)")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(doctor,editado)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error insercion relacion editado")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}

func insertRelationEditTesis(original int,editado int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("INSERT INTO reledittesis (original,editado) VALUES ($1,$2)")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(original,editado)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error insercion relacion editado")
    fmt.Println(err)
  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}

func insertInstitucion(nombre string,url string){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
      
  }
  stmt, err := tx.Prepare("INSERT INTO institution (name,url) VALUES ($1,$2)")
  if err != nil {
      tx.Rollback()
      fmt.Println("no consigue crear sentencia")
  }

  _, err = stmt.Exec(nombre,url)
  if err != nil {
    tx.Rollback()
    fmt.Println("Error insercion institucion")
    fmt.Println(err)

  }

  err = tx.Commit()
  if err != nil {
      stmt.Close()
      tx.Rollback()
     
  }
  stmt.Close()
}

func insertTesis(titulo string,fecha string, url string, abstract string, lector int,departamento string,institution int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
  }

  if institution == 0 {
    stmt, err := tx.Prepare("INSERT INTO thesis (title,abstract,url,defensedate,reader,department,confirmacion) VALUES ($1,$2,$3,$4,$5,$6,$7)")
      if err != nil {
          tx.Rollback()
          fmt.Println("no consigue crear sentencia")
      }

      _, err = stmt.Exec(titulo,abstract,url,fecha,lector,departamento,2)
      if err != nil {
        tx.Rollback()
        fmt.Println("Error insercion tesis")
        fmt.Println(err)
      }
      err = tx.Commit()
      if err != nil {
          stmt.Close()
          tx.Rollback()
         
      }
      stmt.Close()
  }else{
    stmt, err := tx.Prepare("INSERT INTO thesis (title,abstract,url,defensedate,reader,department,institution,confirmacion) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)")
    if err != nil {
        tx.Rollback()
        fmt.Println("no consigue crear sentencia")
    }

    _, err = stmt.Exec(titulo,abstract,url,fecha,lector,departamento,institution,2)
    if err != nil {
      tx.Rollback()
      fmt.Println("Error insercion tesis")
      fmt.Println(err)
    }
    err = tx.Commit()
    if err != nil {
        stmt.Close()
        tx.Rollback()
       
    }
    stmt.Close()
  }
  fmt.Println("termina insercion")

  
}

//Error institucion pasarlo a int en el handler para q se vuelva 0
func updateTesis(titulo string,fecha string, url string, abstract string,departamento string,idInstitucion int,id int){
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    fmt.Println("no consigue abrirlo")
  }

  tx, err := db.Begin()
  if err != nil {
  }
    fmt.Println(titulo,fecha,url,abstract,departamento,idInstitucion,id)
    stmt, err := tx.Prepare("UPDATE thesis SET title=$1, abstract=$2, url=$3, defensedate=$4, department=$5, institution=$6, confirmacion=$7 WHERE id=$8")
      if err != nil {
          tx.Rollback()
          fmt.Println("no consigue crear sentencia")
      }

      _, err = stmt.Exec(titulo,abstract,url,fecha,departamento,idInstitucion,0,id)
      if err != nil {
        tx.Rollback()
        fmt.Println("Error insercion tesis")
        fmt.Println(err)
      }
      err = tx.Commit()
      if err != nil {
          stmt.Close()
          tx.Rollback()
         
      }
      stmt.Close()
  
}

func insertDirector(dir int,tesis int){
  if dir!=0 {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
      fmt.Println("no consigue abrirlo")
    }

    tx, err := db.Begin()
    if err != nil {
        
    }
    stmt, err := tx.Prepare("INSERT INTO directorthesis (director,thesis) VALUES ($1,$2)")
    if err != nil {
        tx.Rollback()
        fmt.Println("no consigue crear sentencia")
    }

    _, err = stmt.Exec(dir,tesis)
    if err != nil {
      tx.Rollback()
      fmt.Println("Error insercion director")
      fmt.Println(err)
    }

    err = tx.Commit()
    if err != nil {
        stmt.Close()
        tx.Rollback()
       
    }
    stmt.Close()

  }
}

func insertKeyword(palabra string, id string){
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
      fmt.Println("no consigue abrirlo")
    }

    tx, err := db.Begin()
    if err != nil {
        
    }
    stmt, err := tx.Prepare("INSERT INTO keyword (word,id) VALUES ($1,$2)")
    if err != nil {
        tx.Rollback()
        fmt.Println("no consigue crear sentencia")
    }

    _, err = stmt.Exec(palabra,id)
    if err != nil {
      tx.Rollback()
      fmt.Println("Error insercion director")
      fmt.Println(err)
    }

    err = tx.Commit()
    if err != nil {
        stmt.Close()
        tx.Rollback()
       
    }
    stmt.Close()

  }


func insertJurado(jur int,tesis int){
  if jur!=0 {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
      fmt.Println("no consigue abrirlo")
    }

    fmt.Println(jur,tesis,jur,tesis)

    tx, err := db.Begin()
    if err != nil {
        
    }
    stmt, err := tx.Prepare("INSERT INTO jurythesis (jury,thesis) VALUES ($1,$2)")
    if err != nil {
        tx.Rollback()
        fmt.Println("no consigue crear sentencia")
    }

    _, err = stmt.Exec(jur,tesis)
    if err != nil {
      tx.Rollback()
      fmt.Println("Error insercion jurado")
      fmt.Println("jurado",jur)
      fmt.Println("tesis",tesis)
      fmt.Println(err)
    }

    err = tx.Commit()
    if err != nil {
        stmt.Close()
        tx.Rollback()
       
    }
    stmt.Close()

  }
}

func checkErr(err error) {
        if err != nil {
            panic(err)
        }
    }