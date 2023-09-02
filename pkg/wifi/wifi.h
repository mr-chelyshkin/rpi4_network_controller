#ifndef WIFI_H
#define WIFI_H

#include <iwlib.h>

#define WLAN_IFACE "wlan0"
#define MAX_NETWORKS 100

typedef struct wifi_info {
    char ssid[33];
    double freq;
    int quality;
    int level;
} wifi_info;

wifi_info* scan(int* count);
const char* active();
int conn(const char* ssid, const char* password);

#endif
