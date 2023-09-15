#include "linux_wifi.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>




int execute_command(const char *command, char *const args[]) {
    pid_t pid, wpid;
    int status = 0;

    if ((pid = fork()) == 0) {
        // This block will be run by the child process
        if (execvp(command, args) == -1) {
            perror("Command execution failed");
            exit(EXIT_FAILURE);
        }
    } else if (pid < 0) {
        // Forking failed
        perror("Fork failed");
        return -1;
    } else {
        // Parent process waits for the child to terminate
        do {
            wpid = waitpid(pid, &status, WUNTRACED);
        } while (!WIFEXITED(status) && !WIFSIGNALED(status));
    }

    return WEXITSTATUS(status) == 0 ? 0 : -1;
}

int network_conn(const char* ssid, const char* password) {
    goSendToChannel("Starting WiFi connection");

    char *killargs[] = {"killall", "wpa_supplicant", NULL};
    char *rmargs[] = {"rm", "-f", "/var/run/wpa_supplicant/" WLAN_IFACE, NULL};

    if (execute_command("killall", killargs) != 0 ||
        execute_command("rm", rmargs) != 0) {
        return -1;
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

    char *wpaargs[] = {"wpa_supplicant", "-B", "-i", WLAN_IFACE, "-c", config_filename, NULL};

    if (execute_command("wpa_supplicant", wpaargs) != 0) {
        remove(config_filename);
        return -1;
    }

    char *dhclientargs[] = {"dhclient", WLAN_IFACE, NULL};

    if (execute_command("dhclient", dhclientargs) != 0) {
        execute_command("killall", killargs);
        remove(config_filename);
        return -1;
    }

    if (remove(config_filename) != 0) {
        perror("Failed to remove configuration file");
    }
    return 0;
}
