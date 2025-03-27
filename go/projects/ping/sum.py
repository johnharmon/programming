#!/opt/homebrew/bin/python3 
import sys 

result_file = sys.argv[1] 

sum = 0
with open(result_file, 'r') as f:
    lines = f.readlines()
    seconds = 0
    for line in lines:
        columns = line.split() 
        time = columns[-2] 
        time = time.strip('s') 
        times = time.split(':')
        if len(times) > 1:
            seconds += float(60 * int(times[0]))
        seconds += float(times[-1])
    print(f'Average time: {round(seconds/len(lines), 2)}')


        
