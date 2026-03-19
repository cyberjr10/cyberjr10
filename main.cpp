#include <iostream>
#include <cstring>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <fcntl.h>
#include <time.h>
#include <errno.h>

#define LISTEN_PORT 8053
#define BUFFER_SIZE 1400

using namespace std;

int main() {
    int local_fd = socket(AF_INET, SOCK_DGRAM, 0);
    int remote_fd = socket(AF_INET, SOCK_DGRAM, 0);

    if (local_fd < 0 || remote_fd < 0) return 1;

    // INTERACTIVE PRIORITY (96): High but smooth for ISP towers
    int tos = 96;
    setsockopt(remote_fd, IPPROTO_IP, IP_TOS, &tos, sizeof(tos));

    // JITTER ABSORPTION: 4KB buffer to smooth out the "shaking"
    int j_buf = 4096;
    setsockopt(remote_fd, SOL_SOCKET, SO_RCVBUF, &j_buf, sizeof(j_buf));
    setsockopt(local_fd, SOL_SOCKET, SO_SNDBUF, &j_buf, sizeof(j_buf));

    struct sockaddr_in servaddr{}, cliaddr{}, target_addr{};
    socklen_t addr_len = sizeof(cliaddr);
    bool linked = false;

    // Set Non-blocking
    fcntl(local_fd, F_SETFL, O_NONBLOCK);
    fcntl(remote_fd, F_SETFL, O_NONBLOCK);

    servaddr.sin_family = AF_INET;
    servaddr.sin_port = htons(LISTEN_PORT);
    servaddr.sin_addr.s_addr = INADDR_ANY;

    if (bind(local_fd, (struct sockaddr *)&servaddr, sizeof(servaddr)) < 0) return 1;

    target_addr.sin_family = AF_INET;
    target_addr.sin_port = htons(53);
    target_addr.sin_addr.s_addr = inet_addr("8.8.8.8");

    cout << "--- STARK ENGINE NATIVE BOOT ---" << endl;
    cout << "Listening on Port: " << LISTEN_PORT << endl;

    unsigned char buffer[BUFFER_SIZE];
    struct timespec req = {0, 500000}; // 0.5ms smoothing clock

    while (true) {
        // Handle Client -> ISP
        ssize_t n = recvfrom(local_fd, buffer, BUFFER_SIZE, 0, (struct sockaddr *)&cliaddr, &addr_len);
        if (n > 0) {
            if (!linked) {
                // Your custom success message
                cout << "\n\033[1;32m[+] CONNECTION SUCCESSFUL\033[0m" << endl;
                cout << "\033[1;36m[!] SMOOTH-RESUME ACTIVE | JITTER BUFFER 4KB | 0.5ms\033[0m\n" << endl;
                linked = true;
            }
            // PACED REDUNDANCY
            sendto(remote_fd, buffer, n, 0, (struct sockaddr *)&target_addr, sizeof(target_addr));
            usleep(100); 
            sendto(remote_fd, buffer, n, 0, (struct sockaddr *)&target_addr, sizeof(target_addr));
        }

        // Handle ISP -> Client
        struct sockaddr_in f_addr{};
        socklen_t f_len = sizeof(f_addr);
        ssize_t r = recvfrom(remote_fd, buffer, BUFFER_SIZE, 0, (struct sockaddr *)&f_addr, &f_len);
        if (r > 0 && cliaddr.sin_addr.s_addr != 0) {
            sendto(local_fd, buffer, r, 0, (struct sockaddr *)&cliaddr, addr_len);
        }

        nanosleep(&req, NULL);
    }
    return 0;
}
