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

type Transition struct {
	partition int
	start     int
	end       int
	before    int
	after     int
}

var startValue int

func ValidarDatosFDISK(tokens []string) {

	fmt.Println("=============== COMANDOS/FDISK - FUNCIÓN VALIDAR DATOS FDISK - ENTRADA ===============")
	fmt.Println(tokens)

	size := ""
	driveletter := ""
	name := ""
	unit := "k"
	tipo := "P"
	fit := "WF"
	delete := ""
	add := ""

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Comparar(tk[0], "size") {
			size = tk[1]
		} else if Comparar(tk[0], "driveletter") {
			driveletter = tk[1]
		} else if Comparar(tk[0], "name") {
			name = tk[1]
		} else if Comparar(tk[0], "unit") {
			unit = tk[1]
		} else if Comparar(tk[0], "type") {
			tipo = tk[1]
		} else if Comparar(tk[0], "fit") {
			fit = tk[1]
		} else if Comparar(tk[0], "delete") {
			delete = tk[1]
		} else if Comparar(tk[0], "add") {
			add = tk[1]
		}
	}

	if delete != "" && add != "" {
		Error("FDISK", "No se puede agregar y eliminar particiones al mismo tiempo.")
		return
	}

	if delete != "" {
		eliminarParticion(driveletter, name, delete)
		return

	}

	if add != "" {
		agregarEspacioParticion(driveletter, name, unit, add)
		return
	}

	if size == "" || driveletter == "" || name == "" {
		Error("FDISK", "El comando FDISK necesita parametros obligatorios: SIZE, DRIVELETTER y/o NAME")
		return
	} else {
		prevPath := directorioActual()
		if prevPath == "" {
			Error("FDISK", "No se ha encontrado el directorio actual.")
			return
		}

		path := prevPath + "/MIA/P1/" + driveletter + ".dsk"

		generarParticion(size, unit, path, tipo, fit, name)
	}
}

func generarParticion(s string, u string, p string, t string, f string, n string) {
	startValue = 0
	i, error_ := strconv.Atoi(s)
	if error_ != nil {
		Error("FDISK", "Size debe ser un número entero")
		return
	}
	if i <= 0 {
		Error("FDISK", "Size debe ser mayor que 0")
		return
	}

	if Comparar(u, "b") || Comparar(u, "k") || Comparar(u, "m") {
		if Comparar(u, "k") {
			i = i * 1024
		} else if Comparar(u, "m") {
			i = i * 1024 * 1024
		}
	} else {
		Error("FDISK", "Unit no contiene los valores esperados.")
		return
	}
	if !(Comparar(t, "p") || Comparar(t, "e") || Comparar(t, "l")) {
		Error("FDISK", "Type no contiene los valores esperados.")
		return
	}
	if !(Comparar(f, "bf") || Comparar(f, "ff") || Comparar(f, "wf")) {
		Error("FDISK", "Fit no contiene los valores esperados.")
		return
	}

	if !ArchivoExiste(p) {
		Error("FDISK", "El disco no existe.")
		return

	}

	mbr := leerDisco(p)

	if int64(i) > mbr.Mbr_tamano {
		Error("FDISK", "El tamaño de la partición es mayor que el tamaño del disco.")
		return
	}

	particiones := GetParticiones(*mbr)
	var between []Transition

	usado := 0
	ext := 0
	c := 0
	base := int(unsafe.Sizeof(Structs.MBR{}))
	extended := Structs.NewParticion()

	for j := 0; j < len(particiones); j++ {
		prttn := particiones[j]
		if prttn.Part_status == '1' {
			var trn Transition
			trn.partition = c
			trn.start = int(prttn.Part_start)
			trn.end = int(prttn.Part_start + prttn.Part_size)
			trn.before = trn.start - base
			base = trn.end
			if usado != 0 {
				between[usado-1].after = trn.start - (between[usado-1].end)
			}
			between = append(between, trn)
			usado++

			if prttn.Part_type == "e"[0] || prttn.Part_type == "E"[0] {
				ext++
				extended = prttn
			}
		}
		if usado == 4 && !Comparar(t, "l") {
			Error("FDISK", "Limite de particiones alcanzado")
			return
		} else if ext == 1 && Comparar(t, "e") {
			Error("FDISK", "Solo se puede crear una partición extendida")
			return
		}
		c++
	}
	if ext == 0 && Comparar(t, "l") {
		Error("FDISK", "Aún no se han creado particiones extendidas, no se puede agregar una lógica.")
		return
	}
	if usado != 0 {
		between[len(between)-1].after = int(mbr.Mbr_tamano) - between[len(between)-1].end
	}
	regresa := BuscarParticiones(*mbr, n, p)
	if regresa != nil {
		Error("FDISK", "El nombre: "+n+", ya está en uso.")
		return
	}
	temporal := Structs.NewParticion()
	temporal.Part_status = '1'
	temporal.Part_size = int64(i)
	temporal.Part_type = strings.ToUpper(t)[0]
	temporal.Part_fit = strings.ToUpper(f)[0]
	copy(temporal.Part_name[:], n)

	if Comparar(t, "l") {
		Logica(temporal, extended, p)
		return
	}
	mbr = ajustar(*mbr, temporal, between, particiones, usado)
	if mbr == nil {
		return
	}
	file, err := os.OpenFile(strings.ReplaceAll(p, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
	}
	file.Seek(0, 0)
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, mbr)
	EscribirBytes(file, binario2.Bytes())
	if Comparar(t, "E") {
		ebr := Structs.NewEBR()
		ebr.Part_status = '0'
		ebr.Part_start = int64(startValue)
		ebr.Part_size = 0
		ebr.Part_next = -1

		file.Seek(int64(startValue), 0) //5200
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, ebr)
		EscribirBytes(file, binario3.Bytes())
		Mensaje("FDISK", "Partición Extendida: "+n+", creada correctamente.")
		return
	}
	file.Close()
	Mensaje("FDISK", "Partición Primaria: "+n+", creada correctamente.")
}

func eliminarParticion(driveletter string, name string, delete string) {
	path := directorioActual() + "/MIA/P1/" + driveletter + ".dsk"
	mbr := leerDisco(path)
	if mbr == nil {
		Error("FDISK - eliminarParticion", "No se ha encontrado el disco.")
		return
	}

	particion := BuscarParticiones(*mbr, name, path)
	if particion == nil {
		Error("FDISK - eliminarParticion", "No se ha encontrado la partición: "+name)
		return
	}

	if delete != "full" {
		Error("FDISK - eliminarParticion", "Parámetro DELETE únicamente puede ser FULL")
		return
	}

	if delete == "full" {
		particion.Part_status = '0'

		file, err := os.OpenFile(path, os.O_WRONLY, os.ModeAppend)
		if err != nil {
			Error("FDISK", "Error al abrir el archivo del disco.")
			return
		}

		defer file.Close()

		file.Seek(particion.Part_start, 0)

		nullBytes := make([]byte, int(particion.Part_size))

		for i := range nullBytes {
			nullBytes[i] = 0
		}

		_, err = file.Write(nullBytes)
		if err != nil {
			Error("FDISK", "Error al eliminar la partición.")
			return
		}

		Mensaje("FDISK", "Partición: "+name+", eliminada correctamente.")
		return
	}
}

func agregarEspacioParticion(driveletter string, name string, unit string, add string) {
	path := directorioActual() + "/MIA/P1/" + driveletter + ".dsk"
	mbr := leerDisco(path)
	if mbr == nil {
		Error("FDISK - agregarEspacioParticion", "No se ha encontrado el disco.")
		return
	}

	particion := BuscarParticiones(*mbr, name, path)
	if particion == nil {
		Error("FDISK - agregarEspacioParticion", "No se ha encontrado la partición: "+name)
		return
	}

	newSize, error_ := strconv.Atoi(add)
	if error_ != nil {
		Error("FDISK - agregarEspacioParticion", "Add debe ser un número entero")
		return
	}
	if newSize == 0 {
		Error("FDISK - agregarEspacioParticion", "No se puede agregar o restar 0 espacio a la partición.")
		return
	}
	if newSize > 0 {
		Mensaje("FDISK - agregarEspacioParticion", "Agregando espacio a la partición: "+name)
	} else {
		Mensaje("FDISK - agregarEspacioParticion", "Restando espacio a la partición: "+name)
	}

	if Comparar(unit, "b") || Comparar(unit, "k") || Comparar(unit, "m") {
		if Comparar(unit, "k") {
			newSize = newSize * 1024
		} else if Comparar(unit, "m") {
			newSize = newSize * 1024 * 1024
		}
	} else {
		Error("FDISK - agregarEspacioParticion", "Unit no contiene los valores esperados.")
		return

	}

	prevSize := int(particion.Part_size)

	finalSize := int(particion.Part_size) + newSize
	if finalSize > int(mbr.Mbr_tamano) {
		Error("FDISK - agregarEspacioParticion", "No hay suficiente espacio en la partición.")
		return

	}

	if newSize < 0 {
		if prevSize+newSize < 0 {
			Error("FDISK - agregarEspacioParticion", "No se puede reducir el tamaño de la partición a un valor negativo.")
			return
		} else {
			particion.Part_size = particion.Part_size + int64(newSize)
		}
	}

	if newSize > 0 {
		particion.Part_size = particion.Part_size + int64(newSize)
	}

	file, err := os.OpenFile(path, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("FDISK - agregarEspacioParticion", "Error al abrir el archivo del disco.")
		return
	}

	defer file.Close()

	file.Seek(particion.Part_start, 0)

	var binario2 bytes.Buffer

	binary.Write(&binario2, binary.BigEndian, particion)

	EscribirBytes(file, binario2.Bytes())

	Mensaje("FDISK - agregarEspacioParticion", "Partición: "+name+", modificada correctamente. Tamaño anterior: "+strconv.Itoa(prevSize)+"KB, Tamaño nuevo: "+strconv.Itoa(int(particion.Part_size))+"KB")
}

func GetParticiones(disco Structs.MBR) []Structs.Particion {
	var v []Structs.Particion
	v = append(v, disco.Mbr_partition_1)
	v = append(v, disco.Mbr_partition_2)
	v = append(v, disco.Mbr_partition_3)
	v = append(v, disco.Mbr_partition_4)
	return v
}

func BuscarParticiones(mbr Structs.MBR, name string, path string) *Structs.Particion {
	var particiones [4]Structs.Particion
	particiones[0] = mbr.Mbr_partition_1
	particiones[1] = mbr.Mbr_partition_2
	particiones[2] = mbr.Mbr_partition_3
	particiones[3] = mbr.Mbr_partition_4

	ext := false
	extended := Structs.NewParticion()
	for i := 0; i < len(particiones); i++ {
		particion := particiones[i]
		if particion.Part_status == "1"[0] {
			nombre := ""
			for j := 0; j < len(particion.Part_name); j++ {
				if particion.Part_name[j] != 0 {
					nombre += string(particion.Part_name[j])
				}
			}
			if Comparar(nombre, name) {
				return &particion
			} else if particion.Part_type == "E"[0] || particion.Part_type == "e"[0] {
				ext = true
				extended = particion
			}
		}
	}

	if ext {
		ebrs := GetLogicas(extended, path)
		for i := 0; i < len(ebrs); i++ {
			ebr := ebrs[i]
			if ebr.Part_status == '1' {
				nombre := ""
				for j := 0; j < len(ebr.Part_name); j++ {
					if ebr.Part_name[j] != 0 {
						nombre += string(ebr.Part_name[j])
					}
				}
				if Comparar(nombre, name) {
					tmp := Structs.NewParticion()
					tmp.Part_status = '1'
					tmp.Part_type = 'L'
					tmp.Part_fit = ebr.Part_fit
					tmp.Part_start = ebr.Part_start
					tmp.Part_size = ebr.Part_size
					copy(tmp.Part_name[:], ebr.Part_name[:])
					return &tmp
				}
			}
		}
	}
	return nil
}

func GetLogicas(particion Structs.Particion, path string) []Structs.EBR {
	var ebrs []Structs.EBR
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
		return nil
	}
	file.Seek(0, 0)
	tmp := Structs.NewEBR()
	file.Seek(particion.Part_start, 0)

	data := leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &tmp)
	if err_ != nil {
		Error("FDSIK", "Error al leer el archivo")
		return nil
	}
	for {
		if int(tmp.Part_next) != -1 && int(tmp.Part_status) != 0 {
			ebrs = append(ebrs, tmp)
			file.Seek(tmp.Part_next, 0)

			data = leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &tmp)
			if err_ != nil {
				Error("FDSIK", "Error al leer el archivo")
				return nil
			}
		} else {
			file.Close()
			break
		}
	}

	return ebrs
}

func Logica(particion Structs.Particion, ep Structs.Particion, path string) {
	logic := Structs.NewEBR()
	logic.Part_status = '1'
	logic.Part_fit = particion.Part_fit
	logic.Part_size = particion.Part_size
	logic.Part_next = -1
	copy(logic.Part_name[:], particion.Part_name[:])

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("FDISK", "Error al abrir el archivo del disco.")
		return
	}
	file.Seek(0, 0)

	tmp := Structs.NewEBR()
	tmp.Part_status = 0
	tmp.Part_size = 0
	tmp.Part_next = -1
	file.Seek(ep.Part_start, 0) //0

	data := leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &tmp)

	if err_ != nil {
		Error("FDSIK", "Error al leer el archivo")
		return
	}
	if err != nil {
		Error("FDISK", "Error al abrir el archivo del disco.")
		return
	}
	var size int64 = 0
	file.Close()
	for {
		size += int64(unsafe.Sizeof(Structs.EBR{})) + tmp.Part_size
		if (tmp.Part_size == 0 && tmp.Part_next == -1) || (tmp.Part_size == 0 && tmp.Part_next == 0) {
			file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
			logic.Part_start = tmp.Part_start
			logic.Part_next = logic.Part_start + logic.Part_size + int64(unsafe.Sizeof(Structs.EBR{}))
			if (ep.Part_size - size) <= logic.Part_size {
				Error("FDISK", "No queda más espacio para crear más particiones lógicas")
				return
			}
			file.Seek(logic.Part_start, 0)

			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, logic)
			EscribirBytes(file, binario2.Bytes())
			nombre := ""
			for j := 0; j < len(particion.Part_name); j++ {
				nombre += string(particion.Part_name[j])
			}
			file.Seek(logic.Part_next, 0)
			addLogic := Structs.NewEBR()
			addLogic.Part_status = '0'
			addLogic.Part_next = -1
			addLogic.Part_start = logic.Part_next

			file.Seek(addLogic.Part_start, 0)

			var binarioLogico bytes.Buffer
			binary.Write(&binarioLogico, binary.BigEndian, addLogic)
			EscribirBytes(file, binarioLogico.Bytes())

			Mensaje("FDISK", "Partición Lógica: "+nombre+", creada correctamente.")
			file.Close()
			return
		}
		file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
		if err != nil {
			Error("FDISK", "Error al abrir el archivo del disco.")
			return
		}
		file.Seek(tmp.Part_next, 0)
		data = leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &tmp)

		if err_ != nil {
			Error("FDSIK", "Error al leer el archivo")
			return
		}
	}
}

func ajustar(mbr Structs.MBR, p Structs.Particion, t []Transition, ps []Structs.Particion, u int) *Structs.MBR {
	if u == 0 {
		p.Part_start = int64(unsafe.Sizeof(mbr))
		startValue = int(p.Part_start)
		mbr.Mbr_partition_1 = p
		return &mbr
	} else {
		var usar Transition
		c := 0
		for i := 0; i < len(t); i++ {
			tr := t[i]
			if c == 0 {
				usar = tr
				c++
				continue
			}

			if Comparar(string(mbr.Dsk_fit[0]), "F") {
				if int64(usar.before) >= p.Part_size || int64(usar.after) >= p.Part_size {
					break
				}
				usar = tr
			} else if Comparar(string(mbr.Dsk_fit[0]), "B") {
				if int64(tr.before) >= p.Part_size || int64(usar.after) < p.Part_size {
					usar = tr
				} else {
					if int64(tr.before) >= p.Part_size || int64(tr.after) >= p.Part_size {
						b1 := usar.before - int(p.Part_size)
						a1 := usar.after - int(p.Part_size)
						b2 := tr.before - int(p.Part_size)
						a2 := tr.after - int(p.Part_size)

						if (b1 < b2 && b1 < a2) || (a1 < b2 && a1 < a2) {
							c++
							continue
						}
						usar = tr
					}
				}
			} else if Comparar(string(mbr.Dsk_fit[0]), "W") {
				if int64(usar.before) >= p.Part_size || int64(usar.after) < p.Part_size {
					usar = tr
				} else {
					if int64(tr.before) >= p.Part_size || int64(tr.after) >= p.Part_size {
						b1 := usar.before - int(p.Part_size)
						a1 := usar.after - int(p.Part_size)
						b2 := tr.before - int(p.Part_size)
						a2 := tr.after - int(p.Part_size)

						if (b1 > b2 && b1 > a2) || (a1 > b2 && a1 > a2) {
							c++
							continue
						}
						usar = tr
					}
				}
			}
			c++
		}
		if usar.before >= int(p.Part_size) || usar.after >= int(p.Part_size) {
			if Comparar(string(mbr.Dsk_fit[0]), "F") {
				if usar.before >= int(p.Part_size) {
					p.Part_start = int64(usar.start - usar.before)
					startValue = int(p.Part_start)
				} else {
					p.Part_start = int64(usar.end)
					startValue = int(p.Part_start)
				}
			} else if Comparar(string(mbr.Dsk_fit[0]), "B") {
				b1 := usar.before - int(p.Part_size)
				a1 := usar.after - int(p.Part_size)

				if (usar.before >= int(p.Part_size) && b1 < a1) || usar.after < int(p.Part_start) {
					p.Part_start = int64(usar.start - usar.before)
					startValue = int(p.Part_start)
				} else {
					p.Part_start = int64(usar.end)
					startValue = int(p.Part_start)
				}
			} else if Comparar(string(mbr.Dsk_fit[0]), "W") {
				b1 := usar.before - int(p.Part_size)
				a1 := usar.after - int(p.Part_size)

				if (usar.before >= int(p.Part_size) && b1 > a1) || usar.after < int(p.Part_start) {
					p.Part_start = int64(usar.start - usar.before)
					startValue = int(p.Part_start)
				} else {
					p.Part_start = int64(usar.end)
					startValue = int(p.Part_start)
				}
			}
			var partitions [4]Structs.Particion
			for i := 0; i < len(ps); i++ {
				partitions[i] = ps[i]
			}

			for i := 0; i < len(partitions); i++ {
				partition := partitions[i]
				if partition.Part_status != '1' {
					partitions[i] = p
					break
				}
			}
			mbr.Mbr_partition_1 = partitions[0]
			mbr.Mbr_partition_2 = partitions[1]
			mbr.Mbr_partition_3 = partitions[2]
			mbr.Mbr_partition_4 = partitions[3]
			return &mbr
		} else {
			Error("FDISK", "No hay espacio suficiente.")
			return nil
		}
	}
}
