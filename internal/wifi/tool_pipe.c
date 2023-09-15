#include "linux_wifi.h"

#include <stdio.h>
#include <unistd.h>

// Global variables to hold the original stdout and stderr file descriptors.
static int orig_stdout_fd;
static int orig_stderr_fd;

// Pipe used for redirecting output.
static int pipe_fd[2];

/**
 * Redirects the standard output and error streams to a custom channel.
 */
void redirect_output() {
    orig_stdout_fd = dup(STDOUT_FILENO);
    orig_stderr_fd = dup(STDERR_FILENO);

    pipe(pipe_fd);
    dup2(pipe_fd[1], STDOUT_FILENO);
    dup2(pipe_fd[1], STDERR_FILENO);
    close(pipe_fd[1]);
}

/**
 * Resets the standard output and error streams back to their original state.
 */
void reset_output() {
    dup2(orig_stdout_fd, STDOUT_FILENO);
    dup2(orig_stderr_fd, STDERR_FILENO);

    close(orig_stdout_fd);
    close(orig_stderr_fd);
    close(pipe_fd[0]);
}

/**
 * Writes the given buffer to a custom channel.
 *
 * @param fd File descriptor (unused in this context).
 * @param buf Pointer to the buffer containing data to be written.
 * @param count Number of bytes to write from the buffer.
 *
 * @return The number of bytes written, or -1 on error.
 */
int redirected_write(int fd, const void* buf, size_t count) {
    (void) fd;

    char* s = (char*) malloc(count + 1);
    if (!s) {
        perror("Memory allocation failed");
        return -1;
    }
    memcpy(s, buf, count);
    s[count] = '\0';
    goSendToChannel(s);

    free(s);
    return count;
}