#include <unistd.h>
#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>

#include <sys/types.h>

int main(){

	int uid;
	uid = geteuid();

	printf("%d", uid);
	return uid;
}
