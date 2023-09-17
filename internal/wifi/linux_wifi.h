#ifndef LINUX_WIFI_H
#define LINUX_WIFI_H

#include <iwlib.h>
#include <stddef.h>

#define WLAN_IFACE "wlan0"
#define MAX_NETWORKS 100

// Struct to hold information about WiFi networks.
typedef struct wifi_info {
    char   ssid[33];
    double freq;
    int    quality;
    int    level;
} wifi_info;

/**
 * Connect to a WiFi network.
 *
 * @param ssid SSID of the network.
 * @param password Password of the network.
 * @return A result code indicating success or failure.
 */
int network_conn(const char* ssid, const char* password);

/**
 * Scan for available WiFi networks.
 *
 * @param count Pointer to an integer to store the number of networks found.
 * @return A pointer to an array of wifi_info structs containing information about each network.
 */
wifi_info* network_scan(int* count);

/**
 * Get information about the currently connected WiFi network.
 *
 * @return The SSID of the current connection or NULL if not connected.
 */
const char* current_connection();

/**
 * Redirect the standard output and error streams.
 */
void redirect_output(void);

/**
 * Reset the standard output and error streams to their original state.
 */
void reset_output(void);

/**
 * Write data to a custom channel.
 *
 * @param fd File descriptor (unused in this context).
 * @param buf Pointer to the buffer containing data to be written.
 * @param count Number of bytes to write from the buffer.
 * @return The number of bytes written, or -1 on error.
 */
int redirected_write(int fd, const void* buf, size_t count);

/**
 * Send data to a custom channel (to be implemented in Go).
 *
 * @param s Pointer to a null-terminated string containing the data to be sent.
 */
extern void goSendToChannel(char* s);

#endif