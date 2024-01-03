The systemd service for use in ubuntu/debian linux with gnome session to periodically checks system inactivity and put system into sleep or hibernate

The hibernate should be configured as described in https://ubuntuhandbook.org/index.php/2021/08/enable-hibernate-ubuntu-21-10/

dependencies: systemd, gnome, rtcwake

install:
1. clone the repository

2. Follow the instruction https://ubuntuhandbook.org/index.php/2021/08/enable-hibernate-ubuntu-21-10/ to enable hibernation

3. edit `hybrid-sleep.sh` to:

3.1. set `ALLOWED_IDLE_TIME` and `ALLOWED_IDLE_TIME_AC` in milliseconds. After the `ALLOWED_IDLE_TIME` the pc will be put into sleep

3.2 set `AC_SLEEP_TIME` and `BAT_SLEEP_TIME` in minutes. The `BAT_SLEEP_TIME` is the time period after which the pc will be hibernated after putting into sleep. the behavior is similar to systemd suspend-then-hibernate service. The `AC_SLEEP_TIME` is used to periodically wake up to check if the pc is still connected to AC

4. run `sudo ./install.sh`

call `sudo ./install.sh` each time after changing values described in p 3. 
