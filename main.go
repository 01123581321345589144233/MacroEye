package main

import (
	"bigeye/valida"
	//"time"
)

//watcher "bigeye/descobre"

func main() {
	// Só estou a verificar se os servidores que nunca foram vistos
	// online ja estão online
	//filtro := bson.D{{"Valido", false}}
	//filtro_todos := bson.M{}

	//go Valida()
	var pais string = "PT"
	for {

		//fmt.Println("Colocar um rate (2800):")
		//var p1 int
		//fmt.Scan(&p1)

		//watcher.Watcher(p1, watcher.PT, pais)
		//fmt.Println("Guardando...\n")
		valida.Online(pais, 600)
		//guarda.Guarda(400, pais)
		//fmt.Println("Validando...")
		//valida.Valida(pais, 600)
		//time.Sleep(10 * time.Minute)
	}
}
