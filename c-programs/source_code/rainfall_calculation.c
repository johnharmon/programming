#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>

int main(){

	const char *months[12] = {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"};
	int startYear = 2015;
	double yearAvg[6];
	double monthAvg[12];
	double yearTotal;
	double monthTotal;
	bool secondLoop = false;
	double rainfallNumbers[5][12] ={
				{12.8,34.6,45.8,345.6,37.2,85.9,36.3,86.2,46.2,646.8,13.9,1.1},
				{2.8,34.6,45.8,35.6,37.2,85.9,36.,8.2,46.2,66.8,3.9,1.1},
				{12.8,4.6,45.8,35.6,37.2,8.9,3.3,86.2,46.2,66.8,3.9,1.1},
				{1.8,34,45.8,45.6,37.2,5.9,6.3,86.2,46.2,46.8,13.,.1},
				{12.8,3.6,5.8,345,7.2,8,36.,6.2,46.2,46.8,13.9,1.1}
				};
	int outer;
	int inner;
	for (outer=0; outer<5; outer++){
		for (inner=0; inner <12; inner++){
			yearTotal += rainfallNumbers[outer][inner];
			monthAvg[inner] += rainfallNumbers[outer][inner];
		}
		yearAvg[outer] = yearTotal;
		yearAvg[5]+=yearTotal;
		yearTotal = 0;
		}
	printf("Year:       Rainfall:\n");
	for (int year = 0; year<5; year++){
		printf("%d       %.2lf\n", startYear+year, yearAvg[year]);

	}
	printf("\nThe yearly average is %.2lf inches\n\n", yearAvg[5]/5.0);
	printf("Monthly Averages:\n\n");
	repeat : ;
	for (inner=0; inner<12; inner++){
	if (secondLoop == true){
		printf("%.1lf ", monthAvg[inner]/5.0);
	}
	else
	printf("%s  ", months[inner]);
	}
	if ( secondLoop == false ){
		printf("\n");
		secondLoop = true;
		goto repeat;
	}		
	printf("\n");
}
