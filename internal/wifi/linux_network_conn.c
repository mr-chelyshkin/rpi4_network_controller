#ifdef linux
#include "linux_wifi.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>

// Function to execute a command using fork and execvp.
int execute_command(const char *command, char *const args[]) {
    pid_t pid, wpid;
    int status = 0;

    if ((pid = fork()) == 0) {
        if (execvp(command, args) == -1) {
            goSendToChannel("Network connection error, fork: command execution failed");
            exit(EXIT_FAILURE);
        }
    } else if (pid < 0) {
        goSendToChannel("Network connection error, fork failed");
        return -1;
    } else {
        do {
            wpid = waitpid(pid, &status, WUNTRACED);
        } while (!WIFEXITED(status) && !WIFSIGNALED(status));
    }
    return WEXITSTATUS(status) == 0 ? 0 : -1;
}

// Function to connect to a WiFi network given the SSID and password.
int network_conn(const char* ssid, const char* password) {
    goSendToChannel("Starting WiFi connection");

    char *killargs[] = {"killall", "wpa_supplicant", NULL};
    char *rmargs[] = {"rm", "-f", "/var/run/wpa_supplicant/" WLAN_IFACE, NULL};

    goSendToChannel("Stopping any running instances of wpa_supplicant and removing any existing configurations");
    if (execute_command("killall", killargs) != 0 ||
        execute_command("rm", rmargs) != 0) {
        return -1;
    }

    goSendToChannel("Creating a temporary wpa configuration file");
    char config_filename[256];
    snprintf(config_filename, sizeof(config_filename), "/tmp/wpa_conf_%s.conf", ssid);
    FILE *config_file = fopen(config_filename, "w");
    if (!config_file) {
        goSendToChannel("Network connection error, failed to open temporary wpa configuration file");
        execute_command("killall", killargs);
        return -1;
    }

    goSendToChannel("Writing the configuration for wpa_supplicant");
    fprintf(config_file,
            "ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev\n"
            "update_config=1\n"
            "country=US\n\n"
            "network={\n"
            "\tssid=\"%s\"\n"
            "\tpsk=\"%s\"\n"
            "}\n", ssid, password);
    fclose(config_file);

    goSendToChannel("Starting wpa_supplicant with the created configuration");
    char *wpaargs[] = {"wpa_supplicant", "-B", "-i", WLAN_IFACE, "-c", config_filename, NULL};
    if (execute_command("wpa_supplicant", wpaargs) != 0) {
        goSendToChannel("Failed to start wpa_supplicant with the created configuration");
        execute_command("killall", killargs);
        remove(config_filename);
        return -1;
    }

    goSendToChannel("Obtaining IP address using dhclient");
    char *dhclientargs[] = {"dhclient", WLAN_IFACE, NULL};
    if (execute_command("dhclient", dhclientargs) != 0) {
        goSendToChannel("Failed to obtain IP address using dhclient");
        execute_command("killall", wpaargs);
        execute_command("killall", killargs);
        remove(config_filename);
        return -1;
    }
    if (remove(config_filename) != 0) {
        goSendToChannel("Failed to remove configuration file");
    }
    return 0;
}
#endif