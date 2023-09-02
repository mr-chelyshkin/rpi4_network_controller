#include "wifi.h"
#include <sys/ioctl.h>
#include <linux/wireless.h>

const char* active() {
    static char essid[IW_ESSID_MAX_SIZE + 1];
    struct iwreq wrq;

    memset(&wrq, 0, sizeof(wrq));
    strncpy(wrq.ifr_name, WLAN_IFACE, IFNAMSIZ);
    wrq.u.essid.pointer = essid;
    wrq.u.essid.length = sizeof(essid);

    int sock = iw_sockets_open();
    if (ioctl(sock, SIOCGIWESSID, &wrq) < 0) {
        iw_sockets_close(sock);
        return "Not connected";
    }

    essid[wrq.u.essid.length] = '\0';
    iw_sockets_close(sock);
    return essid;
}
