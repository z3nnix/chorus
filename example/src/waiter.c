#include <stdio.h>
#include <string.h>

void waiter(char* input) {
    if (strcmp(input, "Waiter! I'll have a sample for Chorus!") == 0) {
        printf("— your example, sir.\n");
    } else {
        printf("— I'm sorry, but that's not..\n");
    }
}