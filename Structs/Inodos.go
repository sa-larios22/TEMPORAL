package Structs

type Inodos struct {
	I_uid   int64     // UID del usuario propietario del archivo o carpeta
	I_gid   int64     // GID del grupo al que pertenece el archivo o carpeta.
	I_size  int64     // Tamaño del archivo en bytes
	I_atime [16]byte  // Última fecha en que se leyó el inodo sin modificarlo
	I_ctime [16]byte  // Fecha en la que se creó el inodo
	I_mtime [16]byte  // Última fecha en la que se modifica el inodo
	I_block [16]int64 // Array en los que los primeros 12 registros son bloques directos.
	/* El 13 será el número del bloque simple indirecto.
	   El 14 será el número del bloque doble indirecto.
	   El 15 será el número del bloque triple indirecto.
	   Si no son utilizados tendrá el valor: -1. 			*/
	I_type int64 // Indica si es archivo o carpeta. Tendrá los siguientes valores:
	/* 1 = Archivo
	   0 = Carpeta */
	I_perm int64 // Guardará los permisos del archivo o carpeta, Se trabajarán usando los permisos UGO (User Group Other) en su forma octal.
	// https://pbs.twimg.com/tweet_video_thumb/GC0ozizXIAACuBZ.jpg
}

func NewInodos() Inodos {
	var inode Inodos
	inode.I_uid = -1
	inode.I_gid = -1
	inode.I_size = -1
	for i := 0; i < 16; i++ {
		inode.I_block[i] = -1
	}
	inode.I_type = -1
	inode.I_perm = -1
	return inode
}
