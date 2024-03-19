package Comandos

import (
	"fmt"
	"strconv"
	"strings"
)

func ValidarDatosUNMOUNT(tokens []string) {
	fmt.Println("=============== FUNCIÓN VALIDAR DATOS UNMOUNT ===============")
	fmt.Println("Cadena de tokens: " + fmt.Sprint(tokens))

	if len(tokens) > 1 {
		Error("UNMOUNT", "Únicamente se acepta el parámetro ID.")
		return
	}

	id := ""

	for i := 0; i < len(tokens); i++ {
		current := tokens[i]
		comando := strings.Split(current, "=")
		if Comparar(comando[0], "id") {
			id = comando[1]
		} else {
			Error("UNMOUNT", "Parámetro no reconocido: "+comando[0])
			return
		}
	}

	if id == "" {
		Error("UNMOUNT", "El parámetro ID es obligatorio.")
		return
	} else if len(id) > 4 {
		Error("UNMOUNT", "El parámetro ID no puede tener más de 4 caracteres.")
		return
	} else {
		unmount(id)
	}
}

func unmount(id string) {

	fmt.Println("=============== FUNCIÓN UNMOUNT ===============")
	fmt.Println("ID: " + id)

	letra := id[0]

	i, err := strconv.Atoi(string(id[1]))
	if err != nil {
		Error("UNMOUNT", "El primer identificador: "+string(id[1])+"no es válido.")
		return
	}

	if i < 1 || i > 26 {
		Error("UNMOUNT", "El primer identificador: "+string(id[1])+"no es válido.")
		return
	}

	for j := 0; j < 99; i++ {
		if DiscMont[j].Particiones[i].Estado == 1 {
			if DiscMont[j].Particiones[i].Letra == letra {
				DiscMont[j].Particiones[i].Estado = 0
				fmt.Println("Se ha desmontado la partición " + id)
				return
			} else {
				Error("UNMOUNT", "La partición no se encuentra montada "+id)
				return
			}
		}
	}
}
