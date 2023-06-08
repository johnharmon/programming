#include <vector>
#include <map>
#include <string>
#include <iostream>
#include <math.h>
using namespace std;


int main(){
	
	vector<vector<string>> competitions{
						{"one", "two"},
						{"two", "three"},
						{"two", "one"}
						};
	string output=competitions[1][0];
	cout << output << endl;

						/*{"one", "two"},
						{"two", "three"},
						{"two", "one"}
						};*/
	//cout << output << endl;
	
	map <string, int> scores;

	scores.insert({"c++", 1});
	cout << scores.at("c++")  << endl;
	scores.at("c++")+=1;
	cout << scores.at("c++")  << endl;

	for (auto it = scores.cbegin(); it != scores.cend(); it++){
	cout << it -> first << endl;
	}
	int x = -4;
	int y = abs(x);
	cout << y << endl;
	//cout << scores.at("c++") << endl;
	//printf("%s\n", output);
}

