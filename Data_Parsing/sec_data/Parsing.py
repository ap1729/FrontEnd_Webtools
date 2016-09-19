import sys,os,re,csv

#get operator for users
'''
file =open('operator.txt')
oper={}
lvl0bs={}
line= file.readline() 
while line:
 data = line.split()
 if len(data)>0:
  for i in range(0,12): 
  #  print data[i]
    oper[str(data[i])]=i/3
    lvl0bs[str(data[i])]=i%3 +(int(i/3))*57
  line =file.readline()
 else:
  break
file.close()
#operators are got

#lvl 1
file =open('operatorlev1.txt')
lvl1bs={}
line= file.readline() 
while line:
 data = line.split()
 if len(data)>0:
  for i in range(0,12): 
    lvl1bs[str(data[i])]=i%3 +(int(i/3))*57
  line =file.readline()
 else:
  break
#lvl1 oper got
print lvl1bs
'''
file =open('Node_loc1.txt')
from pylab import *
import matplotlib.pyplot as plt

line= file.readline() #reads first line which has headers
count=-1 #row count
linedata=[] #to store linedata
elem =[] #list to store data to write into csv file 
elem.append("node")
elem.append("x")
elem.append("y")
elem.append("OPERATOR")
elem.append("LEVEL0 BS")
elem.append("LEVEL1 BS")
xdata=[]
ydata=[]	
uedatax=[]
uedatay=[]
uecount=0
count=0
#print linedata
with open('Nodelocations.csv', 'w') as csvfile:
  writer = csv.writer(csvfile)
  writer.writerow(elem)
  elem=[]
  line = file.readline()
  linedata=line.split()
  while count<228:#bs 
    elem=[]
    elem.append('BS'+str(count/57))
    elem.append(linedata[0])
    elem.append(linedata[1])
    line = file.readline()
    linedata=line.split()
    writer.writerow(elem)
  while line[0]!=' ': #ue
   print line[0]
   if count==12000:
    break
   count+=1
   elem.append('UE'+(count-228))
   elem.append(linedata[0])
   elem.append(linedata[1])
   line = file.readline()
   linedata=line.split()
   writer.writerow(elem)


plot(xdata,ydata,'ko')
show()


