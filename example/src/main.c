#include <stdio.h>

extern char waiter(char*);

int main() {
    char* request = "Waiter! I'll have a sample for Chorus!";
    printf("— %s\n", request);
    waiter(request);
}