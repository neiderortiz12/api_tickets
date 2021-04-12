package main

import(
	"fmt"
	"net/http"
	"log"
	"time"
	"database/sql"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"github.com/gorilla/mux"
	_"github.com/go-sql-driver/mysql"
	"text/template"
)
var plantillas = template.Must(template.ParseFiles("platillas/index.html"))

type ticket struct {
	ID int `json:"id"`
	User string `json:"User"`
	Date_create string `json:"Date_create"`
	Date_update string `json:"Date_update"`
	State string `json:"State"`
}

func conexionDB()(conexion *sql.DB)  {
	Driver:="mysql"
	User:="user_ticket"
	Pass:="12345"
	Name:="db_ticket"

	conexion, err := sql.Open(Driver,User+":"+Pass+"@tcp(127.0.0.1)/"+Name)
	if err != nil{
		panic(err.Error())
	}
	return conexion
}

func index(w http.ResponseWriter, r *http.Request)  {
	plantillas.ExecuteTemplate(w, "index",nil)
}

func createTicket(w http.ResponseWriter, r *http.Request)  {
	t := time.Now()
	var newTicket ticket
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		fmt.Fprintf(w,"datos invalidos")
	}
	json.Unmarshal(reqBody, &newTicket)
	conexionEstablecida := conexionDB()
	insertarTicket, err := conexionEstablecida.Prepare("INSERT INTO tickets(user,date_create, date_update,state) VALUES(?,?,?,?)")
	if err != nil{
		panic(err.Error())
	}
	newTicket.Date_create=fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day() )
	newTicket.Date_update=fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day() )
	insertarTicket.Exec(newTicket.User, newTicket.Date_create, newTicket.Date_update, newTicket.State)
	fmt.Fprintf(w, "registro exitoso")
}

func deleteTicket(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	ticketID, err := strconv.Atoi(vars["id"])
	if err != nil{
		fmt.Fprintf(w, "id invalido")
		return
	}

	conexionEstablecida := conexionDB()
	deleteTicket, err := conexionEstablecida.Prepare("DELETE FROM tickets WHERE id=?")
	if err != nil{
		panic(err.Error())
	}

	deleteTicket.Exec(ticketID)
	fmt.Fprintf(w, "el registro con id %v fue eliminado", ticketID)
}

func updateTicket(w http.ResponseWriter, r *http.Request)  {
	t := time.Now()
	vars := mux.Vars(r)
	ticketID, err := strconv.Atoi(vars["id"])
	if err != nil{
		fmt.Fprintf(w, "id invalido")
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	var updTicket ticket
	if err != nil{
		fmt.Fprintf(w,"datos invalidos")
	}
	json.Unmarshal(reqBody, &updTicket)
	updTicket.Date_update=fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day() )
	conexionEstablecida := conexionDB()
	updatedTicket, err := conexionEstablecida.Prepare("UPDATE tickets SET user=?, date_update=?, state=? WHERE id=?")
	if err != nil{
		panic(err.Error())
	}
	updatedTicket.Exec(updTicket.User, updTicket.Date_update, updTicket.State, ticketID)
	fmt.Fprintf(w, "el registro con el id %v se actualizo con exito", ticketID)
}

func getTickets(w http.ResponseWriter, r *http.Request)  {
	conexionEstablecida := conexionDB()
	obtenerTicket, err := conexionEstablecida.Query("SELECT * FROM tickets")
	if err != nil{
		panic(err.Error())
	}
	tick := ticket{}
	arregloTickets :=[]ticket{}
	for obtenerTicket.Next(){
		var id int
		var user string
		var date_create string
		var date_update string
		var state string
		err = obtenerTicket.Scan(&id,&user,&date_create,&date_update,&state)
		if err!=nil{
			panic(err.Error())
		}
		tick.ID=id
		tick.User=user
		tick.Date_create=date_create
		tick.Date_update=date_update
		tick.State = state
		arregloTickets= append(arregloTickets,tick)

	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(arregloTickets)
}

func getTicket(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	ticketID, err := strconv.Atoi(vars["id"])
	if err != nil{
		fmt.Fprintf(w, "id invalido")
		return
	}

	conexionEstablecida := conexionDB()
	obtenerTicket, err := conexionEstablecida.Query("SELECT * FROM tickets WHERE id=?", ticketID)
	if err != nil{
		panic(err.Error())
	}
	tick := ticket{}
	for obtenerTicket.Next(){
		var id int
		var user string
		var date_create string
		var date_update string
		var state string
		err = obtenerTicket.Scan(&id,&user,&date_create,&date_update,&state)
		if err!=nil{
			panic(err.Error())
		}
		tick.ID=id
		tick.User=user
		tick.Date_create=date_create
		tick.Date_update=date_update
		tick.State = state
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(tick)
	}
}

func main()  {
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/tickets/create", createTicket).Methods("POST")
	router.HandleFunc("/tickets/delete/{id}", deleteTicket).Methods("DELETE")
	router.HandleFunc("/tickets/edit/{id}", updateTicket).Methods("PUT")
	router.HandleFunc("/tickets", getTickets).Methods("GET")
	router.HandleFunc("/tickets/{id}", getTicket).Methods("GET")
	log.Println("servidor corriendo")
	http.ListenAndServe(":3001",router)
	
}