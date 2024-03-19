package Comandos

import (
	"MIA_P1_202111849/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"
)

func ValidarDatosREP(texto []string) {
	fmt.Println("=============== FUNCIÓN VALIDAR DATOS REP ===============")
	fmt.Println("Cadena de tokens: " + fmt.Sprint(texto))

	name := ""
	path := ""
	id := ""
	//ruta := ""

	for i := 0; i < len(texto); i++ {
		token := texto[i]
		comando := strings.Split(token, "=")

		if Comparar(comando[0], "name") {
			name = comando[1]
		} else if Comparar(comando[0], "path") {
			path = comando[1]
		} else if Comparar(comando[0], "id") {
			id = comando[1]
		} else {
			Error("REP", "Parámetro no reconocido: "+comando[0])
			return
		}
	}

	if name == "" || path == "" || id == "" {
		Error("REP", "Los parámetros NAME, PATH, ID y RUTA son obligatorios.")
		return
	}

	if Comparar(name, "mbr") {
		reporteMBR(path, id)
	} else if Comparar(name, "disk") {
		reporteDISK(path, id)
	}
}

func reporteMBR(path string, id string) {
	fmt.Println("=============== FUNCIÓN REPORTE MBR ===============")
	fmt.Println("Path: " + path)
	fmt.Println("ID: " + id)

	mbr, _, _ := LeerDatosParticionMontada(id)
	if mbr == nil {
		Error("REP", "Error al leer el MBR del disco montado.")
		return
	}

	// Verificar y crear el directorio si no existe
	if err := verificarDirectorio(path); err != nil {
		Error("REP", "Error al verificar o crear el directorio: "+err.Error())
		return
	}

	// Obtener el nombre del archivo del path proporcionado
	dir, fileNameWithExt := filepath.Split(path)
	fileName := strings.TrimSuffix(fileNameWithExt, filepath.Ext(fileNameWithExt))

	// Crear el archivo DOT para el reporte del MBR
	dotContent := "digraph MBR_Report {\n"
	dotContent += "\tlabelloc=top\n"
	dotContent += "\trankdir=TB\n"
	dotContent += "\tnode [shape=plaintext]\n"
	dotContent += "\tedge [style=invis]\n"
	dotContent += "\ttable [\n"
	dotContent += "\t\tlabel=<<table border=\"1\" cellborder=\"1\" cellspacing=\"0\">\n"
	dotContent += "\t\t\t<tr><td colspan=\"2\"> Reporte MBR </td></tr>\n"
	dotContent += fmt.Sprintf("\t\t\t<tr><td>tamano</td><td>%d</td></tr>\n", mbr.Mbr_tamano)
	dotContent += fmt.Sprintf("\t\t\t<tr><td>fecha_creacion</td><td>%s</td></tr>\n", string(mbr.Mbr_fecha_creacion[:]))
	dotContent += fmt.Sprintf("\t\t\t<tr><td>disk_signature</td><td>%d</td></tr>\n", mbr.Mbr_dsk_signature)

	// Agregar la información de las particiones
	dotContent += "\t\t\t<tr><td colspan=\"2\"> Particiones </td></tr>\n"

	particiones := []Structs.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	for _, particion := range particiones {
		// Limpiar el nombre de la partición de caracteres nulos
		particionNombre := strings.TrimRight(string(particion.Part_name[:]), "\x00")
		// Construir la tabla de datos de la partición
		particionTable := fmt.Sprintf("\t\t\t<tr><td colspan=\"2\"> Particion </td></tr>\n"+
			"\t\t\t<tr><td> Status: </td><td>%c</td></tr>\n"+
			"\t\t\t<tr><td> Tipo: </td><td>%c</td></tr>\n"+
			"\t\t\t<tr><td> Fit: </td><td>%c</td></tr>\n"+
			"\t\t\t<tr><td> Start: </td><td>%d</td></tr>\n"+
			"\t\t\t<tr><td> Size: </td><td>%d</td></tr>\n"+
			"\t\t\t<tr><td> Nombre: </td><td>%s</td></tr>\n", particion.Part_status, particion.Part_type, particion.Part_fit, particion.Part_start, particion.Part_size, particionNombre)
		// Agregar la tabla de partición al contenido DOT
		dotContent += particionTable
	}

	dotContent += "\t\t</table>>\n"
	dotContent += "\t]\n"
	dotContent += "}\n"

	dotFilePath := filepath.Join(dir, fileName+".dot")
	err := guardarArchivo(dotFilePath, []byte(dotContent))
	if err != nil {
		Error("REP", "Error al guardar el archivo DOT del reporte del MBR: "+err.Error())
		return
	}

	// Generar la imagen JPG utilizando Graphviz
	jpgFilePath := filepath.Join(dir, fileName)
	err = generarImagenDOT(dotFilePath, jpgFilePath)
	if err != nil {
		Error("REP", "Error al generar la imagen JPG del reporte del MBR: "+err.Error())
		return
	}

	// Reporte generado con éxito
	Mensaje("REP", "Reporte generado con éxito y guardado en: "+jpgFilePath)
}

type InformacionParticion struct {
	Nombre         string
	Tipo           string
	Tamano         int64
	PosicionInicio int64
	// Puedes agregar más campos según sea necesario
}

func guardarArchivo(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0644)
}

func generarImagenDOT(dotFilePath, jpgFilePath string) error {
	cmd := exec.Command("dot", "-Tjpg", "-o", jpgFilePath, dotFilePath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func LeerDatosParticionMontada(id string) (*Structs.MBR, *Structs.Particion, *Structs.EBR) {
	// Obtener el disco montado usando su ID
	var discoMontado DiscoMontado
	encontrado := false
	for _, disco := range DiscMont {
		for _, particion := range disco.Particiones {
			particionNombre := strings.TrimRight(string(particion.Nombre[:]), "\x00")
			if particionNombre == id {
				discoMontado = disco
				encontrado = true
				break
			}
		}
		if encontrado {
			break
		}
	}

	if !encontrado {
		fmt.Println("No se encontró la partición montada con el ID especificado:", id)
		return nil, nil, nil
	}

	// Leer el MBR del disco montado
	path := fmt.Sprintf(id[:1])
	mbr := leerDisco(path)
	if mbr == nil {
		fmt.Println("Error al leer el MBR del disco montado:", path)
		return nil, nil, nil
	}

	// Obtener la partición correspondiente del MBR
	particiones := []Structs.Particion{mbr.Mbr_partition_1, mbr.Mbr_partition_2, mbr.Mbr_partition_3, mbr.Mbr_partition_4}
	var particion *Structs.Particion
	for i := range particiones {

		// Convertir el nombre de la partición a string y limpiar caracteres nulos
		particionNombre := strings.TrimRight(string(particiones[i].Part_name[:]), "\x00")
		nombreDiscoMontado := strings.TrimRight(string(discoMontado.Particiones[0].Nombre[:]), "\x00")

		// Comparar los nombres después de la conversión y limpieza
		if particionNombre == nombreDiscoMontado {
			particion = &particiones[i]
			break
		}
	}

	// Verificar si la partición es extendida para leer los EBRs
	if particion != nil && particion.Part_type == 'E' {
		ebrs := GetLogicas(*particion, path)
		if len(ebrs) == 0 {
			fmt.Println("No se encontraron EBRs para la partición extendida:", string(particion.Part_name[:]))
			return mbr, particion, nil
		}

		// Encontrar el EBR correspondiente usando el nombre de la partición montada
		var ebr *Structs.EBR
		for _, e := range ebrs {
			if string(e.Part_name[:]) == string(discoMontado.Particiones[0].Nombre[:]) {
				ebr = &e
				break
			}
		}
		if ebr == nil {
			fmt.Println("No se encontró el EBR correspondiente para la partición lógica:", string(discoMontado.Particiones[0].Nombre[:]))
		}

		return mbr, particion, ebr
	}

	return mbr, particion, nil
}

func verificarDirectorio(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func reporteDISK(p string, id string) {
	fmt.Println("=============== FUNCIÓN REPORTE DISK ===============")
	fmt.Println("Path: " + p)
	fmt.Println("ID: " + id)

	var path string
	GetMount("REP", id, &path)

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP - DISK", "No se encontró el disco. Línea 242")
		return
	}

	var disk Structs.MBR
	file.Seek(0, 0)

	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &disk)
	if err_ != nil {
		Error("REP - DISK", "Error al leer el archivo")
		return
	}
	file.Close()

	aux := strings.Split(p, ".")
	if len(aux) > 2 {
		Error("REP - DISK", "El nombre del disco no puede contener puntos (.)")
	}
	pd := aux[0] + ".dot"

	carpeta := ""
	direccion := strings.Split(pd, "/")

	fileaux, _ := os.Open(strings.ReplaceAll(pd, "\"", ""))
	if fileaux == nil {
		for i := 0; i < len(direccion); i++ {
			carpeta += "/" + direccion[i]
			if _, err2 := os.Stat(carpeta); os.IsNotExist(err2) {
				os.Mkdir(carpeta, 0777)
			}
		}
		os.Remove(pd)
	} else {
		fileaux.Close()
	}

	partitions := GetParticiones(disk)

	var extended Structs.Particion

	ext := false

	for i := 0; i < 4; i++ {
		if partitions[i].Part_status == 1 {
			if partitions[i].Part_type == "E"[0] || partitions[i].Part_type == "e"[0] {
				ext = true
				extended = partitions[i]
			}
		}
	}

	content := "digraph G {\n rankdir=TB;\n forcelabels=true;\n graph [ dpi = \"600\" ];\n node [shape = plaintext];\n nodo1 [label = <<table>\n <tr>\n"

	var positions [5]int64
	var positionsii [5]int64

	positions[0] = disk.Mbr_partition_1.Part_start - (1 * int64(unsafe.Sizeof(Structs.MBR{})))
	positions[1] = disk.Mbr_partition_2.Part_start - disk.Mbr_partition_1.Part_start + disk.Mbr_partition_1.Part_size
	positions[2] = disk.Mbr_partition_3.Part_start - disk.Mbr_partition_2.Part_start + disk.Mbr_partition_2.Part_size
	positions[3] = disk.Mbr_partition_4.Part_start - disk.Mbr_partition_3.Part_start + disk.Mbr_partition_3.Part_size
	positions[4] = disk.Mbr_tamano + 1 - disk.Mbr_partition_4.Part_start + disk.Mbr_partition_4.Part_size

	copy(positionsii[:], positions[:])

	logic := 0
	tmplogic := ""

	if ext {
		tmplogic = "<tr>\n"
		auxEbr := Structs.NewEBR()

		file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
		if err != nil {
			Error("REP - DISK", "No se encontró el disco. Línea 317")
			return
		}

		file.Seek(extended.Part_start, 0)
		data = leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &auxEbr)
		if err_ != nil {
			Error("REP - DISK", "Error al leer el archivo")
			return
		}
		file.Close()

		var tamGen int64 = 0
		for auxEbr.Part_next != -1 {
			tamGen += auxEbr.Part_size
			res := float64(auxEbr.Part_size) / float64(disk.Mbr_tamano)
			res = res * 100
			tmplogic += "<td>\"EBR\"</td>"
			s := fmt.Sprintf("%.2f", res)
			tmplogic += "<td>" + s + "%" + " de la partición extendida</td>"

			resta := float64(auxEbr.Part_next) - (float64(auxEbr.Part_start) + float64(auxEbr.Part_size))
			resta = resta / float64(disk.Mbr_tamano)
			resta = resta * 10000
			resta = math.Round(resta) / 100

			if resta != 0 {
				s = fmt.Sprintf("%f", resta)
				tmplogic += "<td>" + s + "%" + " libre de la partición extendida</td>"
				logic++
			}
			logic += 2
			file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
			if err != nil {
				Error("REP - DISK", "No se encontró el disco. Línea 353")
				return
			}

			file.Seek(auxEbr.Part_next, 0)
			data = leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &auxEbr)
			if err_ != nil {
				Error("REP - DISK", "Error al leer el archivo")
				return
			}
			file.Close()
		}
		resta := float64(extended.Part_size) - float64(tamGen)
		resta = resta / float64(disk.Mbr_tamano)
		resta = math.Round(resta * 100)
		if resta != 0 {
			s := fmt.Sprintf("%.2f", resta)
			tmplogic += "<td>\"Libre\n" + s + "%" + " de la partición extendida\"</td>"
			logic++
		}
		tmplogic += "</tr>\n"
		logic += 2
	}

	var tanPrim int64

	for i := 0; i < 4; i++ {
		if partitions[i].Part_type == 'E' {
			tanPrim += partitions[i].Part_size
			res := float64(partitions[i].Part_size) / float64(disk.Mbr_tamano)
			res = math.Round(res*10000) / 100
			s := fmt.Sprintf("%.2f", res)
			content += "<td COLSPAN='" + strconv.Itoa(logic) + "'> Extendida\n" + s + "%" + " del disco</td>\n"
		} else if partitions[i].Part_start != -1 {
			tanPrim += partitions[i].Part_size
			res := float64(partitions[i].Part_size) / float64(disk.Mbr_tamano)
			res = math.Round(res*10000) / 100
			s := fmt.Sprintf("%.2f", res)
			content += "<td COLSPAN='2'> Primaria\n" + s + "%" + " del disco</td>\n"
		}
	}

	if tanPrim != 0 {
		libre := disk.Mbr_tamano - tanPrim
		res := float64(libre) / float64(disk.Mbr_tamano)
		res = math.Round(res * 100)
		s := fmt.Sprintf("%.2f", res)
		content += "<td ROWSPAN='2'> Libre \n" + s + "%" + " del disco</td>\n"
	}

	content += "</tr>\n\n"
	content += tmplogic
	content += "</table>>];}\n"

	fmt.Println(content)

	// Crear imagen
	b := []byte(content)
	err_ = ioutil.WriteFile(pd, b, 0644)
	if err_ != nil {
		Error("REP - DISK", "Error al crear el archivo")
		log.Fatal(err_)
		return
	}

	terminacion := strings.Split(p, ".")

	newPath, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(newPath, "-T"+terminacion[1], pd).Output()
	mode := int(0777)
	ioutil.WriteFile(p, cmd, os.FileMode(mode))
	disco := strings.Split(newPath, "/")
	Mensaje("REP", "Reporte DISK de: "+disco[len(disco)-1]+" generado con éxito")
}
