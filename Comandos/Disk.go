package Comandos

import (
	"MIA_P1_202111849/Structs"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

// ValidarDatosMKDISK valida los tokens de entrada para el comando MKDISK
func ValidarDatosMKDISK(tokens []string) {
	// Imprimir los tokens de entrada
	fmt.Println("=============== COMANDOS/DISK - FUNCIÓN VALIDAR DATOS MKDISK - ENTRADA ===============")
	fmt.Println(tokens)

	// Inicializar las variables para los parámetros del comando
	size := ""
	fit := ""
	unit := ""
	error_ := false

	// Recorrer los tokens
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// Ignorar los comentarios
		if strings.HasPrefix(token, "#") {
			continue
		}

		// Dividir el token en el nombre del parámetro y su valor
		tk := strings.Split(token, "=")

		// Verificar el nombre del parámetro y asignar su valor a la variable correspondiente
		if Comparar(tk[0], "fit") {
			if fit == "" {
				fit = tk[1]
			} else {
				Error("MKDISK", "Paŕametro F repetido en el comando: "+tk[0])
				return
			}
		} else if Comparar(tk[0], "size") {
			if size == "" {
				size = tk[1]
			} else {
				Error("MKDISK", "Paŕametro SIZE repetido en el comando: "+tk[0])
				return
			}
		} else if Comparar(tk[0], "unit") {
			if unit == "" {
				unit = tk[1]
			} else {
				Error("MKDISK", "Paŕametro U repetido en el comando: "+tk[0])
				return
			}
		} else {
			Error("MKDISK", "No Se esperaba el Paŕametro "+tk[0])
			error_ = true
			return
		}
	}

	// Asignar valores predeterminados a los parámetros que no se especificaron
	if fit == "" {
		fit = "FF"
	}
	if unit == "" {
		unit = "M"
	}
	if error_ {
		return
	}

	// Verificar que se proporcionaron los parámetros necesarios y que sus valores son válidos
	if size == "" {
		Error("MKDISK", "Se requiere parámetro SIZE para este comando")
		return
	} else if !Comparar(fit, "BF") && !Comparar(fit, "FF") && !Comparar(fit, "WF") {
		Error("MKDISK", "Valores en el parámetro FIT no esperados")
		return
	} else if !Comparar(unit, "k") && !Comparar(unit, "m") {
		Error("MKDISK", "Valores en el parámetro UNIT no esperados")
		return
	} else {
		// Si todos los parámetros son válidos, llamar a la función makeFile para crear el archivo
		makeFile(size, fit, unit)
	}
}

// makeFile crea un nuevo archivo de disco con los parámetros especificados
func makeFile(s string, f string, u string) {
	// Imprimir los parámetros de entrada
	fmt.Println("=============== COMANDOS/DISK - FUNCIÓN MKDISK - ENTRADA ===============")
	fmt.Println("Size: " + s)
	fmt.Println("Fit: " + f)
	fmt.Println("Unit: " + u)

	// Crear una nueva estructura MBR (Master Boot Record)
	var disco = Structs.NewMBR()

	// Convertir el tamaño del disco de string a int
	size, err := strconv.Atoi(s)

	// Si hay un error en la conversión, imprimir un mensaje de error y terminar la función
	if err != nil {
		Error("MKDISK", "SIZE debe ser un número entero")
		return
	}

	// Verificar que el tamaño del disco sea mayor que 0
	if size <= 0 {
		Error("MKDISK", "SIZE debe ser mayor a 0")
		return
	}

	// Convertir el tamaño del disco a bytes dependiendo de la unidad especificada
	if Comparar(u, "M") {
		size = 1024 * 1024 * size // Si la unidad es M (megabytes), multiplicar por 1024*1024
	} else if Comparar(u, "k") {
		size = 1024 * size // Si la unidad es K (kilobytes), multiplicar por 1024
	}

	// Tomar el primer carácter del parámetro fit para usarlo en la creación del disco
	f = string(f[0])

	// Establecer el tamaño del disco
	disco.Mbr_tamano = int64(size)

	// Establecer la fecha de creación
	fecha := time.Now().String()
	copy(disco.Mbr_fecha_creacion[:], fecha)

	// Establecer el número aleatorio para el Disk Signature
	aleatorio, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	entero, _ := strconv.Atoi(aleatorio.String())
	disco.Mbr_dsk_signature = int64(entero)

	// Establecer el ajuste del disco
	copy(disco.Dsk_fit[:], string(f[0]))

	// Establecer las particiones
	disco.Mbr_partition_1 = Structs.NewParticion()
	disco.Mbr_partition_2 = Structs.NewParticion()
	disco.Mbr_partition_3 = Structs.NewParticion()
	disco.Mbr_partition_4 = Structs.NewParticion()

	var filename string
	for i := 0; ; i++ {
		filename = string(rune('A'+i)) + ".dsk"
		path := "/home/sergio/GolandProjects/MIA/P1/" + filename
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Crear el archivo de disco
			file, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}

			// Llenar el archivo con ceros
			zeros := make([]byte, size)
			file.Write(zeros)

			// Escribir la estructura MBR al inicio del archivo
			binary.Write(file, binary.BigEndian, &disco)
			file.Close()

			fmt.Println("Disco creado: " + path)
			break
		}
	}
}

// RMDISK es una función que elimina un archivo de disco basado en la letra de la unidad proporcionada
func RMDISK(tokens []string) {
	// Imprimir los tokens de entrada
	fmt.Println("=============== COMANDOS/DISK - FUNCIÓN RMDISK - ENTRADA ===============")
	fmt.Println(tokens)

	// Verificar que solo se proporcionó un token
	if len(tokens) > 1 {
		Error("RMDISK", "Únicamente se acepta el parámetro DRIVELETTER.")
		return
	}

	// Inicializar la variable para la letra de la unidad
	driveLetter := ""

	// Recorrer los tokens
	for i := 0; i < len(tokens); i++ {
		// Dividir el token en el nombre del parámetro y su valor
		token := tokens[i]
		tk := strings.Split(token, "=")

		// Verificar que el nombre del parámetro es "driveletter"
		if Comparar(tk[0], "driveletter") {
			// Si la letra de la unidad aún no se ha asignado, asignarla
			if driveLetter == "" {
				driveLetter = tk[1]
			} else {
				// Si la letra de la unidad ya se ha asignado, imprimir un error y terminar la función
				Error("RMDISK", "Parámetro DRIVELETTER repetido en el comando: "+tk[0])
				return
			}
		} else {
			// Si el nombre del parámetro no es "driveletter", imprimir un error y terminar la función
			Error("RMDISK", "Parámetro: "+tk[0]+" no esperado.")
			return
		}
	}

	// Verificar que se proporcionó la letra de la unidad
	if driveLetter == "" {
		Error("RMDISK", "Se requiere parámetro DRIVELETTER para este comando")
		return
	} else {
		// Construir la ruta al archivo basado en la letra de la unidad
		path := "/home/sergio/GolandProjects/MIA/P1/" + driveLetter + ".dsk"

		// Verificar que el archivo existe
		if !ArchivoExiste(path) {
			Error("RMDISK", "No se encontró el disco en la ruta indicada: "+path)
			return
		}

		// Pedir confirmación al usuario antes de eliminar el archivo
		if Confirmar("¿Está seguro de que desea eliminar el disco: " + path + " ?") {
			// Intentar eliminar el archivo
			err := os.Remove(path)

			// Si hay un error al eliminar el archivo, imprimir un error y terminar la función
			if err != nil {
				Error("RMDISK", "Error eliminar el archivo.")
				return
			}

			// Si el archivo se eliminó con éxito, imprimir un mensaje de éxito
			Mensaje("RMDISK", "Disco ubicado en "+path+", eliminado exitosamente.")
			return
		} else {
			// Si el usuario canceló la operación, imprimir un mensaje de cancelación
			Mensaje("RMDISK", "Operación: Eliminación del disco "+path+", cancelada.")
			return
		}
	}
}
