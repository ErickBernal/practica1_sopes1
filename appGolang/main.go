package main

//https://www.youtube.com/watch?v=tk61nYJzCA8&t=37s&ab_channel=BorjaScript
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type idCar struct {
	Id primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}
type strLog struct {
	Funcion string `json:"funcion,omitempty"`
	Time    string `json:"time,omitempty"`
}

type strCarro struct {
	Placa  string `json:"placa,omitempty"`
	Marca  string `json:"marca,omitempty"`
	Modelo string `json:"modelo,omitempty"`
	Serie  string `json:"serie,omitempty"`
	Color  string `json:"color,omitempty"`
	Time   string `json:"time,omitempty"`
}

var client *mongo.Client

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// LOG
func newLog(accion string) strLog {
	loc, _ := time.LoadLocation("America/Guatemala")
	t := time.Now().In(loc).Format("2006-01-02 15:04:05")
	var log strLog
	log.Funcion = accion
	log.Time = t
	return log
}
func postLog(nuevoLog strLog) {
	collectionLog := client.Database("goDB").Collection("log")
	ctxLog, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collectionLog.InsertOne(ctxLog, nuevoLog)
}

//----- LOG

// POST
func optPostCors(toAnsw *http.ResponseWriter) {
	(*toAnsw).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*toAnsw).Header().Set("Access-Control-Allow-Methods", "POST")
	(*toAnsw).Header().Set("Access-Control-Allow-Origin", "*")
}
func createCar_options(response http.ResponseWriter, request *http.Request) { // POST
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	optPostCors(&response)
	json.NewEncoder(response).Encode("<repuesta POST>")
}
func createCar(response http.ResponseWriter, request *http.Request) { // POST
	fmt.Println("<createCar>")
	enableCors(&response)
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	var nuevoCarro strCarro
	json.NewDecoder(request.Body).Decode(&nuevoCarro)

	//creando el LOG
	nuevoLog := newLog("Crear")
	//insertando el nuevo Carro
	nuevoCarro.Time = nuevoLog.Time
	collection := client.Database("goDB").Collection("carro")           //creamos la base de datos y la coleccion
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //canal
	result, _ := collection.InsertOne(ctx, nuevoCarro)

	//insertando el nuevo Log
	postLog(nuevoLog)

	json.NewEncoder(response).Encode(result)
}

// -----

// UPDATE
func optUpdateCors(toAnsw *http.ResponseWriter) {
	(*toAnsw).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*toAnsw).Header().Set("Access-Control-Allow-Methods", "PUT")
	(*toAnsw).Header().Set("Access-Control-Allow-Origin", "*")
}
func updateCar_options(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	optUpdateCors(&response)
	json.NewEncoder(response).Encode("<repuesta UPDATE>")
}
func updateCar(response http.ResponseWriter, request *http.Request) {
	fmt.Print("<UpdateCar>")
	enableCors(&response)
	response.Header().Add("content-type", "application/json")
	var car strCarro
	json.NewDecoder(request.Body).Decode(&car)

	//log
	log := newLog("Update")

	//insertar Update
	filter := bson.M{"placa": car.Placa}
	update := bson.M{"$set": bson.M{"marca": car.Marca, "modelo": car.Modelo, "serie": car.Serie, "color": car.Color, "time": log.Time}}

	collection := client.Database("goDB").Collection("carro")
	result, _ := collection.UpdateOne(context.Background(), filter, update)

	//insertar log
	postLog(log)
	//
	json.NewEncoder(response).Encode(result)
}

// GET
func getCars(response http.ResponseWriter, request *http.Request) { // GET
	enableCors(&response)
	fmt.Println("	<getCars>")
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	var listCar []strCarro
	collection := client.Database("goDB").Collection("carro")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)
	//lectura de datos
	for cursor.Next(ctx) {
		var aux_carro strCarro
		cursor.Decode(&aux_carro)
		listCar = append(listCar, aux_carro)
	}

	postLog(newLog("Get"))

	json.NewEncoder(response).Encode(listCar)
}
func getCar(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var aux_carro2 strCarro
	json.NewDecoder(r.Body).Decode(&aux_carro2)
	fmt.Println("<getCar> (placa: ", aux_carro2.Placa, ")")

	w.Header().Add("content-type", "application/json") //aceptar json en el request
	collection := client.Database("goDB").Collection("carro")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var cursor_car strCarro
		cursor.Decode(&cursor_car)
		if cursor_car.Placa == aux_carro2.Placa {
			json.NewEncoder(w).Encode(cursor_car)
			return
		}
	}
}

// DELETE
func optDeleteCors(toAnsw *http.ResponseWriter) {
	(*toAnsw).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*toAnsw).Header().Set("Access-Control-Allow-Methods", "POST")
	(*toAnsw).Header().Set("Access-Control-Allow-Origin", "*")
}
func DeleteCar_options(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	optDeleteCors(&response)
	json.NewEncoder(response).Encode("<repuesta DELETE>")
}
func deleteCar(response http.ResponseWriter, request *http.Request) {
	enableCors(&response)
	response.Header().Add("content-type", "application/json")
	var car strCarro
	json.NewDecoder(request.Body).Decode(&car)
	fmt.Println("<deleteCar: ", car.Placa, " >")
	filter := bson.M{"placa": car.Placa}
	collection := client.Database("goDB").Collection("carro")
	result, _ := collection.DeleteOne(context.Background(), filter)

	postLog(newLog("Delete"))

	json.NewEncoder(response).Encode(result)
}

// ------------------------------- FILTRADOS
// POST
func optPostCors_ModelFilter(toAnsw *http.ResponseWriter) {
	(*toAnsw).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*toAnsw).Header().Set("Access-Control-Allow-Methods", "POST")
	(*toAnsw).Header().Set("Access-Control-Allow-Origin", "*")
}
func ModelFilter_options(response http.ResponseWriter, request *http.Request) { // POST
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	optPostCors_ModelFilter(&response)
	json.NewEncoder(response).Encode("<repuesta POST>")
}

func ModelFilter(response http.ResponseWriter, request *http.Request) { // POST
	//datos parra buscar el carro
	fmt.Println("<Filter model car>")
	enableCors(&response)
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	var findCar strCarro
	json.NewDecoder(request.Body).Decode(&findCar)

	//revisar en todo el listado de los carros
	var listCar []strCarro
	collection := client.Database("goDB").Collection("carro")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)

	//lectura de todos los datos del get
	for cursor.Next(ctx) {
		var aux_carro strCarro
		cursor.Decode(&aux_carro)
		if aux_carro.Modelo == findCar.Modelo {
			fmt.Println("-> ", aux_carro)
			listCar = append(listCar, aux_carro)
		}
	}
	postLog(newLog("Get"))

	json.NewEncoder(response).Encode(listCar)
}

func marcaFilter(response http.ResponseWriter, request *http.Request) { // POST
	//datos parra buscar el carro
	fmt.Println("<Filter Marca Car>")
	enableCors(&response)
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	var findCar strCarro
	json.NewDecoder(request.Body).Decode(&findCar)

	//revisar en todo el listado de los carros
	var listCar []strCarro
	collection := client.Database("goDB").Collection("carro")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)

	//lectura de todos los datos del get
	for cursor.Next(ctx) {
		var aux_carro strCarro
		cursor.Decode(&aux_carro)
		if aux_carro.Marca == findCar.Modelo { //se usa modelo, por que reutilice el mismo INPUT de "App.js"
			fmt.Println("-> ", aux_carro)
			listCar = append(listCar, aux_carro)
		}
	}
	postLog(newLog("Get"))

	json.NewEncoder(response).Encode(listCar)
}

func colorFilter(response http.ResponseWriter, request *http.Request) { // POST
	//datos parra buscar el carro
	fmt.Println("<Filter Color Car>")
	enableCors(&response)
	response.Header().Add("content-type", "application/json") //aceptar json en el request
	var findCar strCarro
	json.NewDecoder(request.Body).Decode(&findCar)

	//revisar en todo el listado de los carros
	var listCar []strCarro
	collection := client.Database("goDB").Collection("carro")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	defer cursor.Close(ctx)

	//lectura de todos los datos del get
	for cursor.Next(ctx) {
		var aux_carro strCarro
		cursor.Decode(&aux_carro)
		if aux_carro.Color == findCar.Modelo { //se usa modelo, por que reutilice el mismo INPUT de "App.js"
			fmt.Println("-> ", aux_carro)
			listCar = append(listCar, aux_carro)
		}
	}
	postLog(newLog("Get"))

	json.NewEncoder(response).Encode(listCar)
}

// ------------------
func main() {
	fmt.Println("Practica1. BackEnd Golang y MongoDB")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //canal
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://database:27017"))
	router := mux.NewRouter()
	//endPoint
	router.HandleFunc("/create", createCar).Methods("POST")
	router.HandleFunc("/create", createCar_options).Methods("OPTIONS")

	router.HandleFunc("/get", getCars).Methods("GET")
	router.HandleFunc("/getCar", getCar).Methods("GET")

	router.HandleFunc("/updateCar", updateCar).Methods("PUT")
	router.HandleFunc("/updateCar", updateCar_options).Methods("OPTIONS")

	router.HandleFunc("/deleteCar", deleteCar).Methods("POST")
	router.HandleFunc("/deleteCar", DeleteCar_options).Methods("OPTIONS")

	//filters
	router.HandleFunc("/filterCarModel", ModelFilter).Methods("POST")
	router.HandleFunc("/filterCarModel", ModelFilter_options).Methods("OPTIONS")

	router.HandleFunc("/filterCarMarca", marcaFilter).Methods("POST")
	router.HandleFunc("/filterCarMarca", ModelFilter_options).Methods("OPTIONS")

	router.HandleFunc("/filterCarColor", colorFilter).Methods("POST")
	router.HandleFunc("/filterCarColor", ModelFilter_options).Methods("OPTIONS")

	puerto := 8000
	fmt.Println("server on port: ", puerto)
	http.ListenAndServe(":"+strconv.Itoa(puerto), router)

	defer client.Disconnect(ctx)
}
