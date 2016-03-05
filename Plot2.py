from pylab import *
from numpy import *
import random
x=linspace(0,200,200)
y=linspace(0,200,200)
dist1=dist2=x
for i in range(0,200):
 x[i]=random.random()
 y[i]=random.random()
 #while(abs((x[i]-0.25)*(x[i]-0.25) + (y[i]-0.5)*(y[i]-0.5))<0.05  or abs((x[i]-0.75)*(x[i]-0.75) + (y[i]-0.5)*(y[i]-0.5))<0.05):
  #x[i]=random.random()
  #y[i]=random.random()

dist1=((x-0.25)*(x-0.25)+(y-0.5)*(y-0.5)) 
dist2=((x-0.75)*(x-0.75)+(y-0.5)*(y-0.5)) 

e=(0.2)*((1/dist1)+(1/dist2))
f=180+40*log(abs(dist1))+40*log(abs(dist2))

e1=e>50
x[e1]=random.random()
y[e1]=random.random()
dist1[e1]=((x[e1]-0.25)*(x[e1]-0.25)+(y[e1]-0.5)*(y[e1]-0.5)) 
dist2[e1]=((x[e1]-0.75)*(x[e1]-0.75)+(y[e1]-0.5)*(y[e1]-0.5))  
e[e1]=(0.2)*((1/dist1[e1])+(1/dist2[e1]))
f[e1]=220+80*log(abs(dist1[e1]))+80*log(abs(dist2[e1]))
e1=e>50
#print e1
#print min(f)
f=f+abs(min(f))+1
#print min(f)
x[0]=0.25
x[1]=0.75
y[0]=y[1]=0.5
e[0]=e[1]=max(e)+1
f[0]=200
f[1]=200
x=30*x
y=30*y-5	
import csv
# write it
with open('testfile2.csv', 'w') as csvfile:
    writer = csv.writer(csvfile)
    writer.writerow(linspace(1,200,200))
    writer.writerow(x) 
    writer.writerow(y)
    writer.writerow(e)
    writer.writerow(f)  
