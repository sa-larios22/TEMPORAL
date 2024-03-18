package main

import (
	"MIA_P1_202111849/Comandos"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var logged = false

func main() {
	Inicio()

	for {
		fmt.Println("=============== INGRESE UN COMANDO ===============")
		fmt.Println("===== Puede salir de la aplicación con 'exit' =====")
		fmt.Println("\t")

		reader := bufio.NewReader(os.Stdin)
		entrada, _ := reader.ReadString('\n')
		fmt.Println("ENTRADA: " + entrada)
		eleccion := strings.TrimRight(entrada, "\r\n")
		fmt.Println("ELECCION: " + eleccion)

		if eleccion == "exit" {
			break
		}

		comando := Comando(eleccion)
		fmt.Println("COMANDO: " + comando)

		eleccion = strings.TrimSpace(eleccion)
		fmt.Println("ELECCION - TRIMMED SPACE: " + eleccion)

		eleccion = strings.TrimLeft(eleccion, comando)
		fmt.Println("ELECCION - TRIMMED LEFT: " + eleccion)

		tokens := SepararTokens(eleccion)
		fmt.Println("TOKENS SEPARADOS: ")
		fmt.Println(tokens)
		funciones(comando, tokens)

		fmt.Println("\tPresione enter para continuar")
		fmt.Scanln()
	}
}

func Inicio() {
	fmt.Println("===========================================================================")
	fmt.Println("\tSergio Andrés Larios Fajardo")
	fmt.Println("\tCarné 202111849")
	fmt.Println("\tDPI 2989 25877 0101")
	fmt.Println("= = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = =")
	fmt.Println("\tIngeniería en Ciencias y Sistemas")
	fmt.Println("\tManejo e Implementación de Archivos - Sección B")
	fmt.Println("\tPrimer Semestre 2024")
	fmt.Println("\tIngeniero William Escobar")
	fmt.Println("\tAuxiliar Daniel Chicas")
	fmt.Println("===========================================================================")
}

func Comando(text string) string {

	fmt.Println("FUNCIÓN COMANDO - ENTRADA: " + text)

	var token string
	terminar := false

	for i := 0; i < len(text); i++ {
		if terminar {
			if string(text[i]) == " " || string(text[i]) == "-" {
				break
			}
			token += string(text[i])
		} else if string(text[i]) != " " && !terminar {
			if string(text[i]) == "#" {
				token = text
			} else {
				token += string(text[i])
				terminar = true
			}
		}
	}

	fmt.Println("FUNCIÓN COMANDO - SALIDA: " + token)

	return token
}

func SepararTokens(texto string) []string {

	fmt.Println("===== FUNCIÓN SEPARAR TOKENS - ENTRADA =====")
	fmt.Println(texto)

	var tokens []string

	if texto == "" {
		return tokens
	}

	texto += " "

	var token string

	estado := 0

	for i := 0; i < len(texto); i++ {
		c := string(texto[i])

		if estado == 0 && c == "-" {
			estado = 1
		} else if estado == 0 && c == "#" {
			continue
		} else if estado != 0 {
			if estado == 1 {
				if c == "=" {
					estado = 2
				} else if c == " " {
					continue
				} else if (c == "P" || c == "p") && string(texto[i+1]) == " " && string(texto[i-1]) == "-" {
					estado = 0
					tokens = append(tokens, c)
					token = ""
				} else if (c == "R" || c == "r") && string(texto[i+1]) == " " && string(texto[i-1]) == "-" {
					estado = 0
					tokens = append(tokens, c)
					token = ""
					continue
				}
			} else if estado == 2 {
				if c == " " {
					continue
				}
				if c == "\"" {
					estado = 3
					continue
				} else {
					estado = 4
				}
			} else if estado == 3 {
				if c == "\"" {
					estado = 4
					continue
				}
			} else if estado == 4 && c == "\"" {
				tokens = []string{}
				continue
			} else if estado == 4 && c == " " {
				estado = 0
				tokens = append(tokens, token)
				token = ""
				continue
			}
			token += c
		}
	}
	return tokens
}

func funciones(token string, tks []string) {
	if token != "" {
		if Comandos.Comparar(token, "EXEC") {
			fmt.Println("=============== FUNCIÓN EXEC ===============")
			FuncionExec(tks)
		} else if Comandos.Comparar(token, "MKDISK") {
			fmt.Println("=============== FUNCIÓN MKDISK ===============")
			Comandos.ValidarDatosMKDISK(tks)
		} else if Comandos.Comparar(token, "RMDISK") {
			fmt.Println("=============== FUNCIÓN RMDISK ===============")
			Comandos.RMDISK(tks)
		} else if Comandos.Comparar(token, "FDISK") {
			fmt.Println("=============== FUNCIÓN FDISK ===============")
			Comandos.ValidarDatosFDISK(tks)
		} else if Comandos.Comparar(token, "MOUNT") {
			fmt.Println("=============== FUNCIÓN MOUNT ===============")
			Comandos.ValidarDatosMOUNT(tks)
		} else if Comandos.Comparar(token, "MKFS") {
			fmt.Println("=============== FUNCIÓN MKFS ===============")
			Comandos.ValidarDatosMKFS(tks)
		} else if Comandos.Comparar(token, "LOGIN") {
			fmt.Println("=============== FUNCIÓN LOGIN ===============")
			if logged {
				Comandos.Error("LOGIN", "Ya hay una sesión iniciada")
			} else {
				logged = Comandos.ValidarDatosLOGIN(tks)
			}
		} else if Comandos.Comparar(token, "LOGOUT") {
			fmt.Println("=============== FUNCIÓN LOGOUT ===============")
			if !logged {
				Comandos.Error("LOGOUT", "No hay una sesión iniciada")
				return
			} else {
				logged = Comandos.CerrarSesion()
			}
		} else if Comandos.Comparar(token, "MKGRP") {
			if !logged {
				Comandos.Error("MKGRP", "No hay una sesión iniciada")
				return
			} else {
				fmt.Println("=============== FUNCIÓN MKGRP ===============")
				Comandos.ValidarDatosGrupos(tks, "MK")
			}
		} else if Comandos.Comparar(token, "RMGRP") {
			if !logged {
				Comandos.Error("RMGRP", "No hay una sesión iniciada")
				return
			} else {
				fmt.Println("=============== FUNCIÓN RMGRP ===============")
				Comandos.ValidarDatosGrupos(tks, "RM")
			}
		} else if Comandos.Comparar(token, "MKUSR") {
			if !logged {
				Comandos.Error("MKUSR", "No hay una sesión iniciada")
				return
			} else {
				fmt.Println("=============== FUNCIÓN MKUSER ===============")
				Comandos.ValidarDatosUsers(tks, "MK")
			}
		} else if Comandos.Comparar(token, "RMUSER") {
			if !logged {
				Comandos.Error("RMUSER", "No hay una sesión iniciada")
				return
			} else {
				fmt.Println("=============== FUNCIÓN RMUSER ===============")
				Comandos.ValidarDatosUsers(tks, "RM")
			}
		}
	}
}

func FuncionExec(tokens []string) {

	fmt.Println("===== FUNCIÓN 'FunciónExec' - ENTRADA =====")
	fmt.Println(tokens)

	path := ""
	for i := 0; i < len(tokens); i++ {
		datos := strings.Split(tokens[i], "=")
		if Comandos.Comparar(datos[0], "path") {
			path = datos[1]
		}
	}
	if path == "" {
		Comandos.Error("EXEC", "Se requiere el parámetro PATH para este comando")
		return
	}
	Exec(path)
}

func Exec(path string) {

	fmt.Println("===== FUNCIÓN EXECUTE - ENTRADA =====")
	fmt.Println(path)

	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("Error al abrir el archivo: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		texto := fileScanner.Text()
		texto = strings.TrimSpace(texto)
		tk := Comando(texto)

		if texto != "" {
			if Comandos.Comparar(tk, "pausa") {
				fmt.Println("=============== FUNCIÓN PAUSE ===============")
				var pause string
				Comandos.Mensaje("PAUSE", "Presione ENTER para continuar.")
				fmt.Scanln(&pause)
				continue
			} else if string(texto[0]) == "#" {
				fmt.Println("=============== COMENTARIO ===============")
				Comandos.Mensaje("COMENTARIO", texto)
				continue
			}
			texto = strings.TrimLeft(texto, tk)
			tokens := SepararTokens(texto)
			funciones(tk, tokens)
		}
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error al leer el archivo: %s", err)
	}

	file.Close()
}
