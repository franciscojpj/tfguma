package main

import(
	"github.com/julienschmidt/httprouter"
)

func Router(r *httprouter.Router){
  //LOGIN
  r.POST("/insertUserOrcid",InsertUserOrcid)
  r.POST("/guardarCookie",WriteCookie)

  //Comunes
  r.GET("/getIdCookie",GetIdCookie)

  //Interfaz login
  r.POST("/crearUsuario",CrearUsuario)

  //Interfaz doctor
  r.GET("/getUsuario",GetUserInfo)
  r.GET("/askDeleteUser",AskDeleteUser)

  //Interfaz Tesis
  r.GET("/askDeleteTesis",AskDeleteTesis)

  //Interfaz edicion
  r.POST("/editUser", EditarUsuario)

  //Interfaz busqueda
  r.GET("/searchResults",FindResults)

  //Interfaz Submit
  r.GET("/searchUsuarios",FindUsers)
  r.POST("/insertUser",insertarUsuario)
  r.POST("/insertTesis",insertarTesis)

  //Interfaz Institucion
  r.GET("/getInstitucion",GetInstitucion)


  r.GET("/login",LoginHandler)
  r.GET("/orcid",OrcidHandler)
  r.POST("/syslogin",InternLoginHandler)


  
  r.GET("/visualizacion",VisualizacionGrafo)// Modificada


  
  r.GET("/getSupervisores",GetSupervisors)
  //r.GET("/getAlumnos",GetAlumns)
  r.POST("/pruebalogin",ReadPost)//Va fuera 
  r.GET("/getTesis", GetTesisUser)

  r.GET("/getDireccion", GetDireccion)
  r.GET("/getJurado", GetJurado)
  r.GET("/getKeywords", GetKeywords)

  //Interfaces
  r.GET("/",VisualizacionSearch)
  r.GET("/submitData",VisualizacionSubmit)
  r.GET("/editData",VisualizacionEdit)
  r.GET("/searchData",VisualizacionSearch)
  r.GET("/admin",VisualizacionAdmin)
  r.GET("/thesis",VisualizacionThesis)
  r.GET("/doctor",VisualizacionDoctor)
  r.GET("/institucion",VisualizacionInstitucion)
  r.GET("/faq",VisualizacionFAQ)
  
 
  r.POST("/insertInstitucion",insertarInstitucion)
  
  r.POST("/editarTesis",editarTesis)
  r.POST("/insertDoctores",insertarDoctores)
  r.POST("/insertJurado",insertarJurado)
  
  r.GET("/searchInstituciones",FindInstituciones)
  r.GET("/getInstituciones",GetAllInstitutions)




  r.GET("/getFathers",GetFathers)
  r.GET("/getSons",GetSons)



  r.GET("/getNewTesis",GetNewTesis)
  r.GET("/getDelTesis",GetDelTesis)
  r.GET("/getEditedTesis",GetEditedTesis)
  r.GET("/getNewUsuarios",GetNewUsuarios)
  r.GET("/getDelUsuarios",GetDelUsuarios)
  r.GET("/getEditedUsuarios",GetEditedUsuarios)
  r.GET("/getTesisID",GetTesisId)
  r.GET("/getUserByTesis",GetUserByTesis)
  r.POST("/updateUsuario",UpdateUsuario)
  r.POST("/updateTesis",UpdateTesis)

  r.POST("/deleteTesis",DeleteTesis)
  r.POST("/deleteUsuario",DeleteUsuario)
  r.POST("/deleteJuryTesis",DeleteJury)
  r.POST("/deleteDirectorTesis",DeleteDirector)

  r.GET("/closeSession",CloseSession)
  
  
}