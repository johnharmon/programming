#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>

int main(){

	int seq1[10]={5,1,22,25,6,-1,8,10};
	int seq2[4]={1,6,-1,10};
	int lastMatchedIndex=0;
	int seq1Len=sizeof(seq1)/sizeof(seq1[0]);
	int seq2Len=sizeof(seq2)/sizeof(seq1[0]);
	_Bool isSubSeq=true;
	_Bool matchedThisIteration=false;
	int outer;
	int inner;
	for (outer=0; outer<seq2Len; outer++){
		matchedThisIteration=false;
		for (inner=lastMatchedIndex; inner<seq1Len; inner++){
			if (seq1[inner]==seq2[outer]){
				matchedThisIteration=true;
				lastMatchedIndex=inner;
				break;
			}
		}
		if (matchedThisIteration==false){
			isSubSeq=false;
			break;
		}
	}
	if (isSubSeq)
		printf("true\n");
	else
		printf("false\n");
}
			

