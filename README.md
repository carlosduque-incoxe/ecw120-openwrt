# OpenWrt para EnGenius ECW120 v1.0.0

Port comunitario experimental para el **ECW120 HW v1.0.0**. La prueba de banco partió del firmware EnGenius 1.3.75-1. Las dos imágenes factory usan OpenWrt 24.10.4 y el mismo kernel ARM gzip. La variante UniFi añade el agente comunitario de antoniocosta con una adaptación para leer perfiles de mesh.

| Imagen en images/ | Uso |
| --- | --- |
| OpenWrt-24.10.4-ECW120-base-factory.bin | OpenWrt independiente |
| OpenWrt-24.10.4-ECW120-unifi-factory.bin | OpenWrt con el agente UniFi instalado y desactivado inicialmente |

## Qué se ha comprobado

En un ECW120 físico arrancó la compilación de banco. El agente consiguió la adopción por Ethernet y recibió configuraciones de UniFi Network. El ejecutable de la variante pública coincide byte por byte con el instalado en ese equipo. Se verificaron el modelo, las sumas internas, la tabla de particiones, el kernel y la raíz de ambas imágenes públicas. Se retiraron los ajustes de acceso del banco antes de empaquetarlas. **Las dos imágenes públicas están pendientes de una prueba física de instalación y primer arranque.**

El kernel FIT ocupa 3 604 552 bytes dentro de una partición de 4 MiB. Quedan 589 752 bytes de margen. La imagen de fábrica contiene HLOS y rootfs. Este port solo genera imágenes factory para este modelo.

## Instalación y primer acceso

1. Comprueba en la etiqueta el modelo y la revisión. Descarga una sola variante y verifica su SHA-256 con el archivo SHA256SUMS. En Windows puedes usar Get-FileHash con el nombre del archivo.
2. Con el equipo alimentado y conectado por Ethernet, guarda tus ajustes de EnGenius. Desde la actualización manual del panel local original, selecciona la imagen factory correspondiente y espera a que termine el reinicio. Si el panel rechaza el archivo, detén la instalación. El procedimiento de recuperación y el retorno al firmware original están pendientes de validación.
3. Haz el primer arranque con un PC y el AP en un segmento aislado. Las imágenes recuperan los valores iniciales de OpenWrt. Prueba con el PC en 192.168.1.2/24 y accede mediante SSH a root@192.168.1.1. Establece una contraseña con passwd antes de conectar el AP a tu red.
4. Para usarlo como AP puente en una LAN que ya tiene router, configura una dirección de gestión adecuada y desactiva los servicios DHCP y RA del AP antes de conectarlo a esa LAN.

## Gestión opcional desde UniFi

En la variante con agente, introduce la URL inform de tu propia consola. Sustituye DIRECCION_DE_TU_CONSOLA antes de ejecutar estos comandos en el AP.

```sh
uci set unifi-openwrt.main.inform_url='http://DIRECCION_DE_TU_CONSOLA:8080/inform'
uci set unifi-openwrt.main.enabled='1'
uci commit unifi-openwrt
/etc/init.d/unifi-openwrt enable
/etc/init.d/unifi-openwrt restart
```

Completa la adopción por cable en UniFi Network. El agente emula el perfil de un UAP-AC-Pro. Las funciones que presenta la consola pueden superar las que implementa el agente.

- **Mesh con U7 Lite pendiente de validación.** La consola de prueba envió mesh habilitado y una interfaz de enlace, pero el perfil recibido carecía del SSID y la clave interna del mesh. El agente quedó esperando. La adaptación admite la clave explícita del sitio y una asociación inequívoca con el perfil de autenticación del enlace. La adopción por cable sí se verificó.
- **WPA-Enterprise con RADIUS.** El agente de esta versión deja desactivados los SSID WPA2-Enterprise y WPA3-Enterprise. El paquete instalado es wpad-mbedtls.
- **Recuperación, actualización desde OpenWrt, VLAN sobre mesh y estabilidad a largo plazo.** Esta publicación no proporciona un procedimiento probado para esas funciones.

## Compilar desde fuentes

En un entorno de compilación Linux para OpenWrt, clona openwrt/openwrt en la etiqueta v24.10.4. Desde la raíz de esa copia limpia, define PUB como la ruta absoluta a este repositorio y ejecuta:

```sh
patch -p1 < "$PUB/source/openwrt-24.10.4-ecw120.patch"
cp "$PUB/source/config.buildinfo" .config
make defconfig
make -j"$(nproc)"
```

Para la variante UniFi, prepara otra copia del árbol o añade el agente antes de compilar la imagen. Usa el repositorio antoniocosta/unifi-openwrt en el commit 19659bdbf07c5890e9ad3e7963c0ce85f2946a53. Desde su raíz, con Go 1.26 y OWRT apuntando al árbol OpenWrt:

```sh
git apply "$PUB/source/agent-mesh.patch"
cp "$PUB/source/mesh_profile_test.go" cmd/uow/
go test ./cmd/uow
mkdir -p "$OWRT/files/usr/bin" "$OWRT/files/etc/config" "$OWRT/files/etc/init.d"
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -buildvcs=false -trimpath -ldflags='-s -w' -o "$OWRT/files/usr/bin/uow" ./cmd/uow
install -m755 packaging/files/etc/init.d/unifi-openwrt "$OWRT/files/etc/init.d/unifi-openwrt"
sed 's|http://192.168.1.1:8080/inform||' packaging/files/etc/config/unifi-openwrt > "$OWRT/files/etc/config/unifi-openwrt"
```

Vuelve al árbol OpenWrt y ejecuta make. Las imágenes publicadas se reempaquetaron desde la compilación de banco tras retirar su configuración local. El procedimiento desde fuentes reconstruye el perfil y no se ha comprobado que reproduzca las sumas byte por byte.

Fuentes y licencias de los proyectos originales: [OpenWrt 24.10.4](https://github.com/openwrt/openwrt/tree/v24.10.4) y [unifi-openwrt en el commit usado](https://github.com/antoniocosta/unifi-openwrt/tree/19659bdbf07c5890e9ad3e7963c0ce85f2946a53). La licencia y los avisos del agente se incluyen en source/.
