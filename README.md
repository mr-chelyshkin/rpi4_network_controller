# Network

> This project is in MVP stage.  
> Problems or errors may occur.
> Pull Requests or issues are welcome

This is a console tool for manage wireless network connection on Raspberry.
# [project was moved](https://github.com/Mad-Pixels/wf)
### menu
![menu](https://github.com/mr-chelyshkin/rpi4_network_controller/blob/main/.img/main_menu.jpeg)

### scan
![scan](https://github.com/mr-chelyshkin/rpi4_network_controller/blob/main/.img/scan_wifi.jpeg)

### connect
![connect](https://github.com/mr-chelyshkin/rpi4_network_controller/blob/main/.img/connect.jpeg)

## Requirements:
 - Raspberry Pi
 - Raspbian OS (Debian)

## Install:
> binary require root privileges as default 
> because it works with iwlib directly. 

```shell
sudo apt-get install -y libiw-dev libcap2-bin
sudo rfkill unblock wifi

# Download binary from project releases

sudo mv ./network /usr/bin/network
# grant privileges to `/usr/bin/network`
sudo setcap cap_net_raw,cap_net_admin=eip /usr/bin/network

# run
network
```
