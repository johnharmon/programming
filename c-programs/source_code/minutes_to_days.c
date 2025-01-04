#include <stdio.h>
#include <stdlib.h>
#include <math.h>


int main(){

	int minutes;
	int days=0;
	int years=0;
	int hours=0;
	int minutes_remaining = 0;
	printf("Enter the amount of minutes you would like to convert: ");
	scanf("%d", &minutes);
	getchar();

//#########################Two ways to calculate this:######################################################
//
	int totalhours = minutes/60;
	int totaldays = totalhours/24;
	int totalyears = totaldays/365;
	int remaininghours = totalhours%24;
	int remainingdays = totaldays%365;
	int remainingminutes = minutes%60;
	printf("Total years: %d\n", totalyears);
	printf("Remaining days: %d\n", remainingdays);
	printf("Remaining hours: %d\n", remaininghours);
	printf("Remaining minutes: %d\n", remainingminutes);

//####################################################### Second Way ####################################

	int minutesInDay = 60*24;
	days = minutes/minutesInDay;
	years = days / 365;
	days = days % 365;	
	minutes_remaining = minutes%60;
	hours = minutes / 60 % 24;
	printf("You entered %d minutes\n", minutes);
	printf("This converts to %d years, %d days, %d hours, and %d minutes\n", years, days, hours, minutes_remaining);

}

