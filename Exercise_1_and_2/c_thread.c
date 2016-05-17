/*
  Øving 1/2 - C threads
*/

#include <pthread.h>
#include <stdio.h>
#include <unistd.h>


// Shared variables
static int i = 0;
pthread_mutex_t mymutex;


void* threadFunc1(){
  for(int j= 0; j<9998;j++){
    pthread_mutex_lock(&mymutex);
    i++;
    pthread_mutex_unlock(&mymutex);
  }
  pthread_exit(NULL);
}

void* threadFunc2(){
  for(int k= 0; k<10000;k++){
    pthread_mutex_lock(&mymutex);
    i--;
    pthread_mutex_unlock(&mymutex);
  }
  pthread_exit(NULL);
}


int main(){
    // Initilaize mutex
    pthread_mutex_init(&mymutex,NULL);

    pthread_t thread1;
    pthread_create(&thread1,NULL,threadFunc1,NULL);
    pthread_t thread2;
    pthread_create(&thread2,NULL,threadFunc2,NULL);

    pthread_join(thread1,NULL);
    pthread_join(thread2,NULL);

    printf("i = %d\n", i);

    pthread_mutex_destroy(&mymutex);
    pthread_exit(NULL);
}
