# Tarea 1 - Sistemas Distribuidos
Integrantes: Felipe Lillo y Natalia Romero
- API: https://developers.google.com/books
- Lenguaje: GO 1.20.3
---
### **Instalación:**
Verificar si se tiene docker instalado
```
docker -v
```
> De lo contrario, instalar con https://docs.docker.com/get-docker/
> 
Verificar si se tiene GO instalado
```
go -v
```
> De lo contrario, instalar con https://go.dev/doc/install
>
Clonar repositorio en su local
```
git clone https://github.com/natalia-romero/t1_sd.git
cd t1_sd
```
### **Configuración de contenedores:**
Ejecutar el siguiente comando para crear los contenedores de redis:
```
sudo docker compose up -d --build
```
Luego, ejecutar el siguiente comando para correr contenedores de redis:
```
sudo docker compose up
```
Buscar ID de cada contenedor redis:
```
sudo docker ps -a
```
Con la ID, se debe sacar la IP de cada contenedor redis:
```
sudo docker inspect <CONTAINER_ID >
```
Ya teniendo la IP de cada contenedor, se debe cambiar la línea 86, 91 y 96 en main.go con la IP que corresponda.

### **Correr programa:**
Ejecutar el siguiente comando:
```
go run main.go
```
Ahora ya puede acceder a la API y podrá buscar N libros (este parámetro es configurable en la línea 15 de main.go)
