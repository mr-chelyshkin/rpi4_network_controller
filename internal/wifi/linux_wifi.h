#ifndef LINUX_WIFI_H
#define LINUX_WIFI_H

#include <iwlib.h>
#include <stddef.h>

#define WLAN_IFACE "wlan0"
#define MAX_NETWORKS 100

typedef struct wifi_info {
    char ssid[33];
    double freq;
    int quality;
    int level;
} wifi_info;

int         network_conn(const char* ssid, const char* password);
wifi_info*  network_scan(int* count);
const char* current_connection();

void redirect_output(void);
void reset_output(void);

#endif
