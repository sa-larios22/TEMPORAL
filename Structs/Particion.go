package Structs

type Particion struct {
	Part_status byte     // Indica si la partición está montada o no
	Part_type   byte     // Indica el tipo de partición, primaria o extendida
	Part_fit    byte     // Tipo de ajuste de la partición. (B, F o W)
	Part_start  int64    // Indica en qué byte del disco inicia la partición
	Part_size   int64    // Tamaño total de la partición en bytes
	Part_name   [16]byte // Nombre de la partición
}

func NewParticion() Particion {
	var Part Particion
	Part.Part_status = '0'
	Part.Part_type = 'P'
	Part.Part_fit = 'F'
	Part.Part_start = -1
	Part.Part_size = 0
	Part.Part_name = [16]byte{}
	return Part
}
