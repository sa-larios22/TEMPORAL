package Structs

type Content struct {
	B_name  [12]byte // Nombre de la carpeta o archivo
	B_inodo int64    // Apuntador hacia un inodo asociado al archivo o carpeta
}

func NewContent() Content {
	var cont Content
	cont.B_inodo = -1
	return cont
}
