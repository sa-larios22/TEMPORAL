package Comandos

import (
	"MIA_P1_202111849/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"
)

func Comparar(a string, b string) bool {
	if strings.ToUpper(a) == strings.ToUpper(b) {
		return true
	}
	return false
}

func Error(op string, mensaje string) {
	fmt.Println("\tERROR: " + op + "\n\tDESCRIPCIÓN: " + mensaje)
}

func Mensaje(op string, mensaje string) {
	fmt.Println("\tCOMANDO: " + op + "\n\tTIPO: " + mensaje)
}

func Confirmar(mensaje string) bool {
	fmt.Println(mensaje + " (y/n)")
	var respuesta string
	fmt.Scanln(&respuesta)
	return Comparar(respuesta, "y")
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func EscribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func leerDisco(path string) *Structs.MBR {

	fmt.Println("===== FUNCIÓN LEER DISCO - ENTRADA =====")
	fmt.Println("PATH: " + path)

	m := Structs.MBR{}
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		Error("FDISK - leerDisco (Additional Commands)", "Error al abrir el archivo")
		return nil
	}

	file.Seek(0, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &m)

	if err_ != nil {
		Error("FDSIK - leerDisco (Additional Commands)", "Error al leer el archivo")
		return nil
	}

	var mDir *Structs.MBR = &m
	return mDir
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes

	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func directorioActual() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}
