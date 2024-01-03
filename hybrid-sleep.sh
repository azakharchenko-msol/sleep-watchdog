#!/bin/bash

ALLOWED_IDLE_TIME=900
AC_SLEEP_TIME=60
BAT_SLEEP_TIME=10
while true; do
  # Get list of all logged-in users
  logged_in_users=$(who | cut -d' ' -f1 | sort -u)
  MAX_IDLE_TIME=0
  # Loop through each logged-in user
  for user in $logged_in_users; do
    DISPLAY=""
    SUDO_UID=$(id -u $user)
    display_info=$(w -h $user | awk '$3 ~ /:[0-9.]*/{print $3}')
    if [[ $display_info == *":"* ]]; then
      # This is a X11 display, export DISPLAY
      DISPLAY="DISPLAY=$display_info"
    elif [[ $display_info == *".0" ]]; then
      # This is a Wayland display, export WAYLAND_DISPLAY
      DISPLAY="WAYLAND_DISPLAY=$display_info"
    fi
    # Use su to run dbus-send as the user and store the idle time in a variable

    idle_time=$(sudo -u $user $DISPLAY DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/$SUDO_UID/bus dbus-send --print-reply --dest=org.gnome.Mutter.IdleMonitor /org/gnome/Mutter/IdleMonitor/Core org.gnome.Mutter.IdleMonitor.GetIdletime | awk '{print $NF}' | tail -n 1)
    if [ "$idle_time" -ge "$MAX_IDLE_TIME" ]; then
      MAX_IDLE_TIME="$idle_time"
    fi
    echo "Idle time for $user: $idle_time"
  done
  echo "system idle $MAX_IDLE_TIME"
  if [ "$ALLOWED_IDLE_TIME" -ge "$MAX_IDLE_TIME" ]; then
    if on_ac_power; then
      echo "On AC power. sleep $AC_SLEEP_TIME min"
      rtcwake -m mem -l 900 -t $(date +%s -d "+$AC_SLEEP_TIME minutes")
    else
      echo "On battery. sleep $BAT_SLEEP_TIME min"
      rtcwake -m mem -l 900 -t $(date +%s -d "+$BAT_SLEEP_TIME minutes")
    fi
    echo "Waked up"
    if [ "$?" -eq 0 ] && ! on_ac_power; then
      echo "Still on battery. hibernate..."
      systemctl hibernate
    fi

  fi
  sleep 60
done
