package Structs

type EBR struct {
	Part_status byte     // Indica si la partición está montada o no
	Part_fit    byte     // Tipo de ajuste de la partición. (B, F o W)
	Part_start  int64    // Indica en qué byte del disco inicia la partición
	Part_size   int64    // Tamaño total de la partición en bytes.
	Part_next   int64    // Byte en el que está el próximo EBR. -1 si no hay siguiente
	Part_name   [16]byte // Nombre de la partición
}

func NewEBR() EBR {
	var eb EBR
	eb.Part_status = '0'
	eb.Part_size = 0
	eb.Part_next = -1
	return eb
}
