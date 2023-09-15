#include "wifi.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>

extern void goSendToChannel(const char* s);

int custom_write1(int fd, const void* buf, size_t count) {
  (void) fd;

  char* s = (char*) malloc(sizeof(count + 1));
  if (s == NULL) return -1;

  memcpy(s, buf, count);
  s[count] = '\0';
  goSendToChannel(s);
  free(s);
  return count;
}

int conn(const char* ssid, const char* password) {
    int retval;
    goSendToChannel("hello");
    retval = system("killall wpa_supplicant");
    if (retval == -1) {
        perror("Failed to run 'killall wpa_supplicant'");
    }

    retval = system("rm -f /var/run/wpa_supplicant/" WLAN_IFACE);
    if (retval == -1) {
        perror("Failed to run 'rm' command");
    }

    char config_filename[256];
    snprintf(config_filename, sizeof(config_filename), "/tmp/wpa_conf_%s.conf", ssid);

    FILE *config_file = fopen(config_filename, "w");
    if (!config_file) {
        perror("Failed to open wpa configuration file");
        return -1;
    }
    fprintf(config_file,
            "ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev\n"
            "update_config=1\n"
            "country=US\n\n"
            "network={\n"
            "\tssid=\"%s\"\n"
            "\tpsk=\"%s\"\n"
            "}\n", ssid, password);
    fclose(config_file);

    char cmd[512];
    snprintf(cmd, sizeof(cmd), "wpa_supplicant -B -i " WLAN_IFACE " -c %s", config_filename);
    int status = system(cmd);
    if (status == -1) {
        perror("Failed to execute wpa_supplicant");
        remove(config_filename);
        return -1;
    }
    if (WEXITSTATUS(status) != 0) {
        perror("wpa_supplicant execution failed");
        remove(config_filename);
        return -1;
    }

    status = system("dhclient " WLAN_IFACE);
    if (status == -1 || WEXITSTATUS(status) != 0) {
        perror("Failed to get IP address using dhclient");
        retval = system("killall wpa_supplicant");
        if (retval == -1) {
            perror("Failed to run 'killall wpa_supplicant' after dhclient failure");
        }
        remove(config_filename);
        return -1;
    }

    remove(config_filename);
    return 0;
}
