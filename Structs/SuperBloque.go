package Structs

import "unsafe"

type SuperBloque struct {
	S_filesystem_type   int64    // Número que identifica el sistema de archivos utilizado
	S_inodes_count      int64    // Número total de inodos
	S_blocks_count      int64    // Número total de bloques
	S_free_blocks_count int64    // Número de bloques libres
	S_free_inodes_count int64    // Número de inodos libres
	S_mtime             [16]byte // Última fecha en que el sistema fue montado
	S_umtime            [16]byte // Última fecha en que el sistema fue desmontado
	S_mnt_count         int64    // Indica cuántas veces se ha montado el sistema
	S_magic             int64    // Valor que identifica al sistema de archivos, tendrá el valor 0xEF53
	S_inode_size        int64    // Tamaño del inodo
	S_block_size        int64    // Tamaño del bloque
	S_firts_ino         int64    // Primer inodo libre
	S_first_blo         int64    // Primer bloque libre
	S_bm_inode_start    int64    // Guardará el inicio del bitmap de inodos
	S_bm_block_start    int64    // Guardará el inicio del bitmap de bloques
	S_inode_start       int64    // Guardará el inicio de la tabla de inodos
	S_block_start       int64    // Guardará el inicio de la tabla de bloques
}

func NewSuperBloque() SuperBloque {
	var spr SuperBloque
	spr.S_magic = 0xEF53
	spr.S_inode_size = int64(unsafe.Sizeof(Inodos{}))
	spr.S_block_size = int64(unsafe.Sizeof(BloquesCarpetas{}))
	spr.S_firts_ino = 0
	spr.S_first_blo = 0
	return spr
}
