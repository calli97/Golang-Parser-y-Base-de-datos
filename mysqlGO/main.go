package main

import (
	"capudo/parser"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

//Las estructuras pueden variar segun que parametros decidamos incluir

//Desestime muchos de los parametros que estan en los dataset de los recorridos porque estan repetidos en los
//datasets de usuarios y de las estaciones. Deje solo los Ids
type Trip struct {
	IdUser    string
	IdOrigin  string
	IdDestiny string
	StartDate string
	EndDate   string
	Duration  int
}

type User struct {
	IdUser     int
	Gender     string
	Age        int
	SingUpDate string
}

//Aca desestime para hacerlo mas simple de igual forma se pueden incluir otros parametros
//Un tema es que la hubicacion puede tener 2 formatos:
//Calle 1 y Calle 2
//Calle 1 y Numero
type Station struct {
	IdStation int
	Name      string
	Latitude  float64
	Longitude float64
}

//Toma la matris devuelta por el parser y retorna un arreglo de viajes
func arraysToTrips(data *parser.Data) []Trip {
	trips := make([]Trip, 0)
	var p Trip

	keys := [6]string{"id_usuario", "id_estacion_origen", "id_estacion_destino", "fecha_origen_recorrido", "fecha_destino_recorrido", "duracion_recorrido"}
	for i := 1; i < len(data.Rows); i++ {
		//fmt.Println(len(data.Rows[i]))
		//Esta harcodeado porque en varios casos se hay mas campos de los que deberia
		for j := 0; j < 16; /*len(data.Rows[i])*/ j++ {
			switch data.Rows[0][j] {
			case keys[0]:
				p.IdUser = data.Rows[i][j]
			case keys[1]:
				p.IdOrigin = data.Rows[i][j]
			case keys[2]:
				p.IdDestiny = data.Rows[i][j]
			case keys[3]:
				p.StartDate = data.Rows[i][j]
			case keys[4]:
				p.EndDate = data.Rows[i][j]
			case keys[5]:
				p.Duration, _ = strconv.Atoi(data.Rows[i][j])
			}
		}
		trips = append(trips, p)
	}
	return trips
}

//Toma la matris devuelta por el parser y retorna un arreglo de estaciones
func arraysToStation(data *parser.Data) []Station {
	stations := make([]Station, 0)
	var p Station
	//Desestime muchos campos igual se cambia rapido si es necesario
	keys := [4]string{"id", "nombre", "long", "lat"}
	for i := 1; i < len(data.Rows); i++ {
		for j := 0; j < len(data.Rows[i]); j++ {
			switch data.Rows[0][j] {
			case keys[0]:
				p.IdStation, _ = strconv.Atoi(data.Rows[i][j])
			case keys[1]:
				p.Name = data.Rows[i][j]
			case keys[2]:
				p.Longitude, _ = strconv.ParseFloat(data.Rows[i][j], 64)
			case keys[3]:
				p.Latitude, _ = strconv.ParseFloat(data.Rows[i][j], 64)
			}
		}
		stations = append(stations, p)
	}
	return stations
}

//Toma la matrix devuelta por el parser y retorna un arreglo de usuarios
func arrayToUser(data *parser.Data) []User {
	users := make([]User, 0)
	var p User
	//Dependiendo el anio se cambio el nombre de los campos
	//Desestime el horario del alta
	//Anios 2020 y 2019
	keys := [4]string{"id_usuario", "genero_usuario", "edad_usuario", "fecha_alta"}
	//Anios 2015-2018
	//keys := [4]string{"usuario_id", "usuario_sexo", "usuario_edad", "fecha_alta"}

	for i := 1; i < len(data.Rows); i++ {
		for j := 0; j < len(data.Rows[i]); j++ {
			switch data.Rows[0][j] {
			case keys[0]:
				p.IdUser, _ = strconv.Atoi(data.Rows[i][j])
			case keys[1]:
				p.Gender = data.Rows[i][j]
			case keys[2]:
				p.Age, _ = strconv.Atoi(data.Rows[i][j])
			case keys[3]:
				p.SingUpDate = data.Rows[i][j]
			default:

			}
		}
		users = append(users, p)
	}
	return users
}

func main() {
	///////////////////////SE DEBE MODIFICAR LOS PARAMETROS PARA ACCEDER A LA BASE DE DATOS//////////////////
	//Conexion con la base de datos
	//                           "user:password@..."
	db, err := sql.Open("mysql", "root:pass@tcp(127.0.0.1:3306)/bicis_caba")
	defer db.Close()

	//Rutas de los datasets reducidos
	pathUsers := string("data/users2019.csv")
	pathStations := string("data/bicicleteros.csv")
	pathTrips := string("data/recorridos2020.csv")

	if err != nil {
		log.Fatal(err)
	}

	//Consulta de prueba
	res, err := db.Query("SELECT * FROM users")

	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		var user User
		err := res.Scan(&user.IdUser, &user.Gender, &user.Age, &user.SingUpDate)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", user)
	}
	//Se generan los arreglos de los distintos tipos de datos en funcion de los datasets
	users, _ := parser.Parser(pathUsers, ",")
	a := arrayToUser(users)
	json.NewEncoder(os.Stdout).Encode(a)

	stations, _ := parser.Parser(pathStations, ";")
	b := arraysToStation(stations)
	json.NewEncoder(os.Stdout).Encode(b)

	trips, _ := parser.Parser(pathTrips, ",")
	c := arraysToTrips(trips)
	json.NewEncoder(os.Stdout).Encode(c)
	//Insercion en la base de datos
	/*
		for _, user := range a {

			db.Exec("INSERT INTO users(user_id,gender,age,sign_up_date) VALUES (?,?,?,?)", user.IdUser, user.Gender, user.Age, user.SingUpDate)
		}
	*/
}
