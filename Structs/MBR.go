package Structs

type MBR struct {
	Mbr_tamano         int64    // Tamaño total en bytes
	Mbr_fecha_creacion [16]byte // Fecha y hora de creación del disco
	Mbr_dsk_signature  int64    // Número aleatorio que identifica cada disco
	Dsk_fit            [1]byte  // Tipo de ajusta de la partición (F, B o W)
	Mbr_partition_1    Particion
	Mbr_partition_2    Particion
	Mbr_partition_3    Particion
	Mbr_partition_4    Particion
}

func NewMBR() MBR {
	var mb MBR
	return mb
}
