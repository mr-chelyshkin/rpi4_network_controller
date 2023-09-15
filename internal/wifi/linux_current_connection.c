#include "linux_wifi.h"

#include <linux/wireless.h>
#include <sys/ioctl.h>

// Function to check the currently connected WiFi network.
const char* current_connection() {
    static char essid[IW_ESSID_MAX_SIZE + 1];
    static char response_msg[256];
    static char error_msg[256];

    struct iwreq wrq;
    memset(&wrq, 0, sizeof(wrq));
    strncpy(wrq.ifr_name, WLAN_IFACE, IFNAMSIZ);
    wrq.u.essid.length = sizeof(essid);
    wrq.u.essid.pointer = essid;

    int sock = iw_sockets_open();
    if (sock < 0) {
        snprintf(
            error_msg,
            sizeof(error_msg),
            "Check current connection error while opening socket: %s",
            strerror(errno)
        );
        return error_msg;
    }
    if (ioctl(sock, SIOCGIWESSID, &wrq) < 0) {
        snprintf(
            error_msg,
            sizeof(error_msg),
            "Check current connection error in ioctl: %s",
            strerror(errno)
        );
        iw_sockets_close(sock);
        return error_msg;
    }
    essid[wrq.u.essid.length] = '\0';
    iw_sockets_close(sock);

    if(strlen(essid) == 0) {
        return "No active connection";
    }
    snprintf(
        response_msg,
        sizeof(response_msg),
        "Connected to: %s",
        essid
    );
    return response_msg;
}
