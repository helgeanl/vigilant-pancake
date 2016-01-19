/*
  Øving 1 - C threads
*/

#include <pthread.h>
#include <stdio.h>

void* someThreadFunction(){
  printf("Hello from a thread!\n");
  return NULL;
}

int main(){
    pthread_t someThread;
    pthread_create(&someThread, NULL,someThreadFunction,NULL);
    printf("Hello from main\n");
    return 0;
}
