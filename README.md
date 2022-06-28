# TERRITORIUM SYNC
Es un script para copiar recursivamente archivos de una o varias carpetas dentro de un mismo bucket S3
uso:
- local:
  - `territoriumsync awsdownload -b mybucket -f directory1,directory2 -l -s ~/Download/`
- Azure:
  - `territoriumsync awsdownload -b mybucket -f directory1,directory2 -a`

cambios:
agregar comando para descargar archivos especificos desde X origen y transmitirlos hacia X destino