import sys,os,re,csv
from numpy import *
file =open('extracted_file1.txt','r')#extracted_file.txt
data= file.read()
data1=data.split()
#print data1[:-1]
count=-1
elem=[]
elem1=[]
''''
for i in range(0,228): #76 is number of basestations, 228 for sectoring
 elem1.append(str(i)) #"A"+str(i) 
''' 
with open('sectorloss.csv', 'w') as csvfile:
       writer = csv.writer(csvfile)
       #writer.writerow(elem1)
       for num in data1:
         count+=1
         if count%228==0 and count>0: #76 is number of basestations, 228 for sectoring
           writer.writerow(elem)
           elem=[]
         try:
        # if this succeeds it is float
           elem.append(float(num))
         except ValueError:
           try: 
             elem.append(float(num[1:]))#print float(num[1:]) #elem.append(float(num[1:])) 
           #except ValueError:
           except ValueError:
             print ""  
       writer.writerow(elem1)
      
'''
find=re.search('(-\d+\.*\d* \;*)+',data)
print find.group(0)
#with open('ConvertedData.csv', 'w') as csvfile:
#    writer = csv.writer(csvfile)
'''
