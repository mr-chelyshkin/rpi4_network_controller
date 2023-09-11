#include "wifi.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

wifi_info* scan(int* count) {
    wireless_scan_head head;
    wireless_scan *result;
    iwrange range;

    int sock;
    static wifi_info networks[MAX_NETWORKS];

    sock = iw_sockets_open();
    if (iw_get_range_info(sock, WLAN_IFACE, &range) < 0) {
        return NULL;
    }
    if (iw_scan(sock, WLAN_IFACE, range.we_version_compiled, &head) < 0) {
        return NULL;
    }

    int i = 0;
    result = head.result;
    while (result != NULL && i < MAX_NETWORKS) {
        if (result->b.has_essid && result->b.essid_on) {
            strncpy(networks[i].ssid, result->b.essid, 32);
            networks[i].ssid[32] = '\0';
            networks[i].freq = result->b.freq;
            networks[i].quality = result->stats.qual.qual;
            networks[i].level = result->stats.qual.level;
            i++;
        }
        result = result->next;
    }
    *count = i;

    iw_sockets_close(sock);
    return networks;
}
