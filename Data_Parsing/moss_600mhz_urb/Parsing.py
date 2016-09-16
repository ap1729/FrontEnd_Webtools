import sys,os,re,csv

#get operator for users
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

file =open('Nodelocations1.txt')
from pylab import *
import matplotlib.pyplot as plt

line= file.readline() #reads first line which has headers
count=-1 #row count
linedata=[] #to store linedata
elem =[] #list to store data to write into csv file 
elem.append("col1")
elem.append("node")
elem.append("x")
elem.append("y")
elem.append("OPERATOR")
elem.append("LEVEL0 BS")
elem.append("LEVEL1 BS")
xdata=[]
ydata=[]	
bs1x=[]
bs2x=[]
bs3x=[]
bs4x=[]
ue1x=[]
ue2x=[]
ue3x=[]
ue4x=[]
bs1y=[]
bs2y=[]
bs3y=[]
bs4y=[]
ue1y=[]
ue2y=[]
ue3y=[]
ue4y=[]
uecount=0
#print linedata
with open('Nodelocations.csv', 'w') as csvfile:
  writer = csv.writer(csvfile)
  while line and (line[0]=="B" or line[0]=="N" or  line[0]=="U"): #for sectoring 2184,else 2032
    #print line[0]
    writer.writerow(elem)
    elem=[]
    linedata=[]
    line = file.readline()
    linedata=line.split()
   # print linedata
    count+=1
    elem.append(count)  #can use node id as well , but redundancy 
     
    if linedata[0][0]=='B':  #means it is basestation
      if linedata[0][2]=='1':
       bs1x.append(linedata[2])
       bs1y.append(linedata[3])
      elif linedata[0][2]=='2':
       bs2x.append(linedata[2])
       bs2y.append(linedata[3])
      elif linedata[0][2]=='3':
       bs3x.append(linedata[2])
       bs3y.append(linedata[3])
      else:
        bs4x.append(linedata[2])
        bs4y.append(linedata[3])
      circle1=plt.Circle((linedata[2],linedata[3]),10,color='grey',clip_on=False)
      fig = plt.gcf()
      fig.gca().add_artist(circle1)
    else :#means it is UE
      if linedata[0][2]=='1':
       ue1x.append(linedata[2])
       ue1y.append(linedata[3])
      elif linedata[0][2]=='2':
       ue2x.append(linedata[2])
       ue2y.append(linedata[3])
      elif linedata[0][2]=='3':
       ue3x.append(linedata[2])
       ue3y.append(linedata[3])
      else:
        ue4x.append(linedata[2])
        ue4y.append(linedata[3])
    
    try:
     elem.append(linedata[0])
    except IndexError:
      print "Indexerror"
    elem.append(linedata[2])
    elem.append(linedata[3])
    xdata.append(linedata[2])
    ydata.append(linedata[3])

    #Add Operator
    if linedata[0][0]=="U":
     uecount+=1
     try:
       elem.append(oper[str(uecount)]) 
       elem.append(lvl0bs[str(uecount)])
       elem.append(lvl1bs[str(uecount)])    
      # print "UE",oper[str(uecount)]
     except KeyError:
       elem.append(-1)
       elem.append(-1)
       elem.append(-1)
    else: 
       #it is basestation
      elem.append(int(linedata[0][2])-1)	
      elem.append(-1)
      elem.append(-1)
  writer.writerow(elem) #For last line 
'''
plot(xdata,ydata,'ko')
show()
'''
plt.plot(bs1x,bs1y,'ro',markersize=20)
plt.plot(bs2x,bs2y,'go',markersize=10)
plt.plot(bs3x,bs3y,'bo',markersize=10)
plt.plot(bs4x,bs4y,'yo',markersize=10)
plt.plot(ue1x,ue1y,'ro')
plt.plot(ue2x,ue2y,'go')
plt.plot(ue3x,ue3y,'bo')
plt.plot(ue4x,ue4y,'yo')
show()
