package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conectionDB() (conexion *sql.DB) {
	driver := "mysql"
	user := "root"
	pass := "P@ssw0rd"
	nombre := "sistema"
	conexion, err := sql.Open(driver, user+":"+pass+"@tcp(127.0.0.1)/"+nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/crear", crear)
	http.HandleFunc("/insertar", insertar)
	http.HandleFunc("/borrar", borrar)
	http.HandleFunc("/editar", editar)
	http.HandleFunc("/actualizar", actualizar)
	log.Println("Listen on Port 8080")
	http.ListenAndServe(":8080", nil)
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func index(w http.ResponseWriter, r *http.Request) {
	conexionEstablecida := conectionDB()
	Registros, err := conexionEstablecida.Query("Select * From empleados")
	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}
	for Registros.Next() {
		var id int
		var nombre, correo string
		err = Registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	plantillas.ExecuteTemplate(w, "index", arregloEmpleado)
}

func crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		conexionEstablecida := conectionDB()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	conexionEstablecida := conectionDB()
	borrarRegistro, err := conexionEstablecida.Prepare("Delete From empleados where id=?")
	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idEmpleado)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	conexionEstablecida := conectionDB()
	Registro, err := conexionEstablecida.Query("Select * From empleados where id=?", idEmpleado)
	empleado := Empleado{}
	for Registro.Next() {
		var id int
		var nombre, correo string
		err = Registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	plantillas.ExecuteTemplate(w, "editar", empleado)
}

func actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		conexionEstablecida := conectionDB()
		actualizarRegistro, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?,correo=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		actualizarRegistro.Exec(nombre, correo, id)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}
