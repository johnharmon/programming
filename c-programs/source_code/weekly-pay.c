#include <stdio.h>
#include <stdlib.h>
#include <math.h>

int main () {

	float payRate;
	float hours;
	float grossPay;
	float netPay;
	float taxes;
	float normalHours;
	float overTime;
	float net15=0;
	float net20=0;
	float net25=0;

	printf("Enter the hourly pay:\n");
	scanf("%f", &payRate);
	getchar();
	printf("Enter hours worked for the week:\n");
	scanf("%f", &hours);
	getchar();

	if (hours >= 40){
		normalHours=40;
		overTime = (float) hours - 40;
	}
	else {
		normalHours=hours;
		overTime = 0;
	}

	grossPay = (normalHours * payRate) + (overTime * payRate * 1.5);

	if ( grossPay >= 300) {
		net15 = 300 * .85;
		if (grossPay >= 450){
			net20 = 150*.80;
			net25 = (grossPay-450)*.75;
		}
		else 
			net20 = (grossPay - 300 ) * .80;
	}
	else{
		net15 = grossPay * .85;
	}
	
	netPay = net15 + net20 + net25;
	
	printf("Net15: %.2f\n", net15 );
	printf("Net20: %.2f\n", net20 );
	printf("Net25: %.2f\n", net25 );
	printf("Normal Hours: %.2f\n", normalHours);
	printf("OT Hours: %.2f\n", overTime);
	printf("Gross pay is : %.2f\n", grossPay);
	printf("Net Pay is : %.2f\n", netPay);
}



