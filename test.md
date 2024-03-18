FDISK
Este comando administra las particiones en el archivo que representa al disco duro. Deberá mostrar un error si no se pudo realizar la operación solicitada sobre la partición, especificando por qué razón no pudo crearse (Por espacio, por restricciones de particiones, etc.). No se considerará el caso de que se pongan parámetros incompatibles, por ejemplo, en un mismo comando fdisk llamar a delete y add. Tendrá los siguientes parámetros:

|Parámetro|Categoría|Descripción|
|---------|---------|-----------|
|-size|Obligatorio al crear|Este parámetro recibirá un número que indicará el tamaño de la partición a crear. Debe ser positivo y mayora cero, de lo contrario se mostrará un mensaje de error.|
|-driveletter|Obligatorio|Este parámetro será la letra del disco a buscar. Si el archivo no existe, debe mostrar un mensaje de error.|
|-name	|Obligatorio|Indicará el nombre de la partición. El nombre no debe repetirse dentro de las particiones de cada disco. Si se va a eliminar, la partición ya debe existir, si no existe debe mostrar un mensaje de error.|
|-unit	|Opcional|Este parámetro recibirá una letra que indicará las unidades que utilizará el parámetro s. Podrá tener los siguientes valores: B: indicará que se utilizarán bytes. K: indicará que se utilizarán Kilobytes(1024 bytes) M:indicará que se utilizarán Megabytes(1024 * 1024) Este parámetro es opcional, si no se encuentra se creará una partición en Kilobytes. Si se utiliza un valor diferente mostrará un mensaje de error.|
|-type|Opcional|Indicará que tipo de partición se creará. Ya que es opcional, se tomará como primaria en caso de que no se indique. Podrá tener los siguientes valores: P: Se creará una partición primaria. E: Se creará una partición extendida. L: Se creará una partición lógica.Si se utiliza otro valor diferente a los anteriores deberá mostrar un mensaje de error. Las particiones lógicas sólo pueden estar dentro de la extendida sin sobrepasar su tamaño. Deberá tener en cuenta las restricciones de teoría de particiones: ● La suma de primarias y extendidas debe ser como máximo 4. ● Solo puede haber una partición extendida por disco. ● No se puede crear una partición lógica si no hay una extendida.|
|-fit|Opcional|Indicará el ajuste que utilizará la partición para asignar espacio. Podrá tener los siguientes valores: BF: Indicará el mejor ajuste (Best Fit) FF: Utilizará el primer ajuste (First Fit) WF: Utilizará el peor ajuste (Worst Fit) Ya que es opcional, se tomará el peor ajuste (WF) si no está especificado en el comando. Si se utiliza otro valor que no sea alguno de los anteriores mostrará un mensaje de error.|
|-delete|Opcional|Este parámetro indica que se eliminará una partición. Este parámetro se utiliza junto con -name y -path. Se deberá mostrar un mensaje que permita confirmar la eliminación de dicha partición. Si la partición no existe deberá mostrar error. Si se elimina la partición extendida, deben eliminarse las particiones lógicas que tenga adentro. Recibirá el único siguiente valor: Full: Esta opción además marcar como vació el espacio en la tabla de particiones, rellena el espacio con el carácter \0. Si se utiliza otro valor diferente, mostrará un mensaje de error.|
|-add|Opcional|Este parámetro se utilizará para agregar o quitar espacio de la partición. Puede ser positivo o negativo. Tomará el parámetro unit para las unidades a agregar o eliminar. En el caso de agregar espacio, deberá comprobar que exista espacio libre después de la partición. En el caso de quitar espacio se debe comprobar que quede espacio en la partición (no espacio negativo).|

Ejemplos:
#Crea una partición primaria llamada Particion1 de 300 kb
#con el peor ajuste en el disco A.dsk
fdisk -size=300 -driveletter=A -name=Particion1
#Crea una partición extendida dentro del disco B.dsk de 300 kb
#Tiene el peor ajuste
fdisk -type=E -driveletter=B -unit=K -name=Particion2 -size=300
#Crea una partición lógica con el mejor ajuste, llamada Partición 3,
#de 1 Mb en el disco B
fdisk -size=1 -type=L -unit=M -fit=bf -driveletter=B -name="Particion3"
#Intenta crear una partición extendida dentro del disco B de 200 kb
#Debería mostrar error ya que ya existe una partición extendida
#dentro de Disco2
fdisk -type=E -driveletter=B -name=Part3 -unit=K -size=200
#Elimina de forma rápida una partición llamada Partición 1
fdisk -delete=full -name="Particion1" -driveletter=A
#Elimina de forma completa una partición llamada Partición 1
fdisk -name=Particion1 -delete=full -driveletter=A
#Quitan 500 Kb de Partición 4 en disco D
#Ignora los demás parámetros (s)
#Se toma como válido el primero que aparezca, en este caso add
fdisk -add=-500 -size=10 -unit=K -driveletter=D -name=”Particion4”
#Agrega 1 Mb a la partición Partición 4 del disco D
#Se debe validar que haya espacio libre después de la partición
fdisk -add=1 -unit=M -driveletter=D -name="Particion4"

Verifica si el código seleccionado cumple con lo requerido