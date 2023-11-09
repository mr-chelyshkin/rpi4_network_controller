#ifdef linux
#include "linux_wifi.h"

#include <stdio.h>
#include <errno.h>
#include <stdlib.h>
#include <string.h>

/**
 * Function to scan for available WiFi networks.
 */
wifi_info* network_scan(int* count) {
    wireless_scan_head head;
    wireless_scan *result;
    iwrange range;
    static wifi_info networks[MAX_NETWORKS];
    char error_msg[256];
    int sock;
    *count = 0;

    sock = iw_sockets_open();
    if (sock < 0) {
        snprintf(error_msg, sizeof(error_msg),
                 "Network scan error while opening socket: %s",
                 strerror(errno));
        goSendToChannel(error_msg);
        return NULL;
    }

    if (iw_get_range_info(sock, WLAN_IFACE, &range) < 0) {
        snprintf(error_msg, sizeof(error_msg),
                 "Network scan error while getting range info: %s",
                 strerror(errno));
        goSendToChannel(error_msg);
        iw_sockets_close(sock);
        return NULL;
    }

    if (iw_scan(sock, WLAN_IFACE, range.we_version_compiled, &head) < 0) {
        snprintf(error_msg, sizeof(error_msg),
                 "Network scan error: %s",
                 strerror(errno));
        goSendToChannel(error_msg);
        iw_sockets_close(sock);
        return NULL;
    }

    result = head.result;
    while (result != NULL && *count < MAX_NETWORKS) {
        if (result->b.has_essid && result->b.essid_on) {
            strncpy(networks[*count].ssid, result->b.essid, sizeof(networks[*count].ssid) - 1);
            networks[*count].ssid[sizeof(networks[*count].ssid) - 1] = '\0';
            networks[*count].quality = result->stats.qual.qual;
            networks[*count].level = result->stats.qual.level;
            networks[*count].freq = result->b.freq;
            (*count)++;
        }
        result = result->next;
    }

    iw_sockets_close(sock);
    return networks;
}
#endif
