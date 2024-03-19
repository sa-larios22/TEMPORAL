package Comandos

import (
	"MIA_P1_202111849/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

var DiscMont [99]DiscoMontado

type DiscoMontado struct {
	Path        [150]byte
	Estado      byte
	Particiones [26]ParticionMontada
}

type ParticionMontada struct {
	Letra  byte
	Estado byte
	Nombre [20]byte
}

var alfabeto = []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func ValidarDatosMOUNT(context []string) {

	fmt.Println("=============== FUNCIÓN VALIDAR DATOS MOUNT - ENTRADA ===============")
	fmt.Println(context)

	driveletter := ""
	name := ""

	for i := 0; i < len(context); i++ {
		current := context[i]
		comando := strings.Split(current, "=")
		if Comparar(comando[0], "name") {
			name = comando[1]
		} else if Comparar(comando[0], "driveletter") {
			driveletter = comando[1]
		}
	}

	prevPath := directorioActual()

	if prevPath == "" {
		Error("MOUNT", "No se ha encontrado el directorio actual.")
		return

	}

	path := prevPath + "/MIA/P1/" + driveletter + ".dsk"

	if path == "" || name == "" {
		Error("MOUNT", "El comando MOUNT requiere parámetros obligatorios")
		return
	}
	mount(path, name)
	listaMount()
}

func mount(p string, n string) {

	fmt.Println("=============== FUNCIÓN MOUNT ===============")
	fmt.Println("Path: " + p)
	fmt.Println("Name: " + n)

	file, error_ := os.Open(p)
	if error_ != nil {
		Error("MOUNT", "No se ha podido abrir el archivo.")
		return
	}

	disk := Structs.NewMBR()
	file.Seek(0, 0)

	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &disk)
	if err_ != nil {
		Error("FDSIK", "Error al leer el archivo")
		return
	}
	file.Close()

	particion := BuscarParticiones(disk, n, p)
	if particion == nil {
		Error("MOUNT", "No se encontró la partición "+n+" en el disco "+p)
		return

	}

	if particion.Part_type == 'E' || particion.Part_type == 'L' {
		var nombre [16]byte
		copy(nombre[:], n)
		if particion.Part_name == nombre && particion.Part_type == 'E' {
			Error("MOUNT", "No se puede montar una partición extendida.")
			return
		} else {
			ebrs := GetLogicas(*particion, p)
			encontrada := false
			if len(ebrs) != 0 {
				for i := 0; i < len(ebrs); i++ {
					ebr := ebrs[i]
					nombreebr := ""
					for j := 0; j < len(ebr.Part_name); j++ {
						if ebr.Part_name[j] != 0 {
							nombreebr += string(ebr.Part_name[j])
						}
					}

					if Comparar(nombreebr, n) && ebr.Part_status == '1' {
						encontrada = true
						n = nombreebr
						break
					} else if nombreebr == n && ebr.Part_status == '0' {
						Error("MOUNT", "No se puede montar una partición Lógica eliminada.")
						return
					}
				}
				if !encontrada {
					Error("MOUNT", "No se encontró la partición Lógica.")
					return
				}
			}
		}
	}

	//prevDriveLetter := path.Base(p)
	//driveletter := strings.TrimSuffix(prevDriveLetter, path.Ext(prevDriveLetter))

	for i := 0; i < 99; i++ {
		var ruta [150]byte
		copy(ruta[:], p)
		if DiscMont[i].Path == ruta {
			for j := 0; j < 26; j++ {
				var nombre [20]byte
				copy(nombre[:], n)
				if DiscMont[i].Particiones[j].Nombre == nombre {
					Error("MOUNT", "La partición "+n+" ya está montada en el disco "+p)
					return
				}
				if DiscMont[i].Particiones[j].Estado == 0 {
					DiscMont[i].Particiones[j].Estado = 1
					DiscMont[i].Particiones[j].Letra = alfabeto[j]
					copy(DiscMont[i].Particiones[j].Nombre[:], n)
					re := string(alfabeto[j]) + strconv.Itoa(j+1) + strconv.Itoa(49)
					Mensaje("MOUNT", "Se ha realizado correctamente el MOUNT en el disco"+p+" con el ID = "+re)
					return
				}
			}
		}
	}

	for i := 0; i < 99; i++ {
		if DiscMont[i].Estado == 0 {
			DiscMont[i].Estado = 1
			copy(DiscMont[i].Path[:], p)
			for j := 0; j < 26; j++ {
				if DiscMont[i].Particiones[j].Estado == 0 {
					DiscMont[i].Particiones[j].Estado = 1
					DiscMont[i].Particiones[j].Letra = alfabeto[j]
					copy(DiscMont[i].Particiones[j].Nombre[:], n)
					re := string(alfabeto[j]) + strconv.Itoa(j+1) + strconv.Itoa(49)
					Mensaje("MOUNT", "Se ha realizado correctamente el MOUNT en el disco "+p+" con -id = "+re)
					return
				}
			}
		}
	}
}

func GetMount(comando string, id string, p *string) Structs.Particion {

	fmt.Println("=============== FUNCIÓN GETMOUNT ===============")
	fmt.Println("Comando: " + comando)
	fmt.Println("Id: " + id)
	fmt.Println("Path: " + *p)

	if !(id[2] == '4' && id[3] == '9') {
		Error(comando, "El primer identificador no es válido.")
		return Structs.Particion{}
	}
	letra := id[len(id)-1]
	id = strings.ReplaceAll(id, "49", "")
	i, _ := strconv.Atoi(string(id[0] - 1))
	if i < 0 {
		Error(comando, "El primer identificador no es válido.")
		return Structs.Particion{}
	}
	for j := 0; j < 26; j++ {
		if DiscMont[i].Particiones[j].Estado == 1 {
			if DiscMont[i].Particiones[j].Letra == letra {

				path := ""
				for k := 0; k < len(DiscMont[i].Path); k++ {
					if DiscMont[i].Path[k] != 0 {
						path += string(DiscMont[i].Path[k])
					}
				}

				file, error := os.Open(strings.ReplaceAll(path, "\"", ""))
				if error != nil {
					Error(comando, "No se ha encontrado el disco")
					return Structs.Particion{}
				}
				disk := Structs.NewMBR()
				file.Seek(0, 0)

				data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
				buffer := bytes.NewBuffer(data)
				err_ := binary.Read(buffer, binary.BigEndian, &disk)

				if err_ != nil {
					Error("FDSIK", "Error al leer el archivo")
					return Structs.Particion{}
				}
				file.Close()

				nombreParticion := ""
				for k := 0; k < len(DiscMont[i].Particiones[j].Nombre); k++ {
					if DiscMont[i].Particiones[j].Nombre[k] != 0 {
						nombreParticion += string(DiscMont[i].Particiones[j].Nombre[k])
					}
				}
				*p = path
				return *BuscarParticiones(disk, nombreParticion, path)
			}
		}
	}
	return Structs.Particion{}
}

func listaMount() {
	fmt.Println("\n========== LISTADO DE MOUNTS ==========")
	for i := 0; i < 99; i++ {
		for j := 0; j < 26; j++ {
			if DiscMont[i].Particiones[j].Estado == 1 {
				nombre := ""
				for k := 0; k < len(DiscMont[i].Particiones[j].Nombre); k++ {
					if DiscMont[i].Particiones[j].Nombre[k] != 0 {
						nombre += string(DiscMont[i].Particiones[j].Nombre[k])
					}
				}
				fmt.Println("\tID = " + strings.ToUpper(string(alfabeto[j])) + strconv.Itoa(i+1) + strconv.Itoa(49) + ", Nombre = " + nombre)
			}
		}
	}
}
