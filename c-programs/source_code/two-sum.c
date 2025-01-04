
#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>

int main(){

	int nums[10] = {1,2,3,4,5,6,7,8,9,10};
	int target=19;
	int length=sizeof(nums)/sizeof(nums[0]);

	int index;
	int nextIndex;
	int answer[2];
	for (index=0; index<length; index++){
		for ( nextIndex=index+1; nextIndex<length; nextIndex++){
			int sum = nums[index]+nums[nextIndex];
			if (sum == target){
				answer[0]=nums[index];
				answer[1]=nums[nextIndex];
			}
		}
	}
	printf("%d %d\n", answer[0], answer[1]);
}


