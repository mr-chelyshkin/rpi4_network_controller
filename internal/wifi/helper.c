#include "linux_wifi.h"
#include <stdio.h>
#include <unistd.h>

static int orig_stdout_fd;
static int orig_stderr_fd;
static int pipe_fd[2];

void redirect_output() {
    orig_stdout_fd = dup(STDOUT_FILENO);
    orig_stderr_fd = dup(STDERR_FILENO);
    pipe(pipe_fd);

    dup2(pipe_fd[1], STDOUT_FILENO);
    dup2(pipe_fd[1], STDERR_FILENO);
    close(pipe_fd[1]);
}

void reset_output() {
    dup2(orig_stdout_fd, STDOUT_FILENO);
    dup2(orig_stderr_fd, STDERR_FILENO);

    close(orig_stdout_fd);
    close(orig_stderr_fd);
    close(pipe_fd[0]);
}
