Como mencionas, el código implementa un sistema de archivos personalizado, por defecto, EXT2, siguiendo las reglas:

Para el caso del sistema de archivos EXT2, se deberán implementar las estructuras como se especifican a continuación. La estructura en bloques es la siguiente:

|Superbloque|Bitmap Inodos|Bitmap Bloques|Inodos|Bloques|
|-|-|-|-|-|

El número de bloques será el triple que el número de inodos. El número de inodos y bloques a crear se puede calcular despejando n de la primera ecuación y aplicando la función floor al resultado:
- tamaño_particion = sizeOf(superblock) + n + 3 * n + n * sizeOf(inodos) + 3 * n * sizeOf(block)
- numero_estructuras = floor(n)

Crea un método para la creación del sistema EXT3 en base a las siguientes especificaciones:

Para el caso del sistema de archivos EXT3, se deberán implementar las estructuras como se especifican a continuación. La estructura en bloques es la siguiente:

|Superbloque|Journaling|Bitmap Inodos|Bitmap Bloques|Inodos|Bloques|
|-|-|-|-|-|-|

El número de bloques será el triple que el número de inodos. El número de Journaling, inodos y bloques a crear se puede calcular despejando n de la primera ecuación y aplicando la función floor al resultado:
- tamaño_particion = sizeOf(superblock) + n + n * sizeOf(Journaling) + 3 * n + n * sizeOf(inodos) + 3 * n * sizeOf(block)
- numero_estructuras = floor(n)

Toma como base la función EXT2, y que la función MKFS verifique si "fs" es "2fs" para ejecutar el sistema EXT2 o si es "3fs" para ejecutar el sistema EXT3, en caso de ser diferente mostrar un mensaje de error.