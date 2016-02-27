from pylab import *
import random
a=linspace(0,100,100)
b=linspace(0,100,100)
for i in range(0,100):
 a[i]=random.random()
 b[i]=2*pi*random.random()
 if(abs(a[i])<0.3):
   a[i]+=0.4


c=array(a*cos(b))
d=array(a*sin(b))
e=1/(a*a)
f=80+40*log(a)	

c=20*c+30
d=20*d+30	
import csv
# write it
with open('test_file.csv', 'w') as csvfile:
    writer = csv.writer(csvfile)
    writer.writerow(linspace(1,100,100))
    writer.writerow(c) 
    writer.writerow(d)
    writer.writerow(e)
    writer.writerow(f)  
