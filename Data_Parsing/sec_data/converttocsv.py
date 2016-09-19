import sys,os,re,csv
import numpy as np
file =open('iteration1.txt','r')#extracted_file.txt
#data= file.read()
#data1=data.split()
#print data1[:-1]
count=-1
w, h =  228,11400
elem = [[0 for x in range(w)] for y in range(h)]
elem1=[]
''''
for i in range(0,228): #76 is number of basestations, 228 for sectoring
 elem1.append(str(i)) #"A"+str(i) 
''' 
with open('sectorloss1.csv', 'w') as csvfile:
      while count<227:   
       line=file.readline()
       linedata=line.split()
       print len(linedata)
       writer = csv.writer(csvfile)
       #writer.writerow(elem1)
       count+=1
       for i in range(0,11400):
         #print i,count
         elem[i][count]=linedata[i]
#       print elem[227][count],linedata[227]
      count=0
      while count<11400:
         elem1=[]
         for i in range(0,228):
          elem1.append(elem[count][i])
         count+=1
         writer.writerow(elem1)
       
      
'''
find=re.search('(-\d+\.*\d* \;*)+',data)
print find.group(0)
#with open('ConvertedData.csv', 'w') as csvfile:
#    writer = csv.writer(csvfile)
'''
