cd ../kernel/ || exit
go build kernel.go

cd ../entradasalida/ || exit
go build entradasalida.go

cd ../memoria/ || exit
go build memoria.go

cd ../cpu/ || exit
go build cpu.go
