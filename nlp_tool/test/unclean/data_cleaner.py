filename = "data_eng.csv"
import csv 

def clean(entry: str):
    entry = entry.lower()
    entry.replace('-', ' ')
    
    return entry

with open(filename, 'r') as f:
    cleaned = open("cleaned_"+filename, 'w')
    csvfile = csv.reader(f)
    for rows in csvfile:
        cleaned.write(clean(rows[0]))
        cleaned.write(",")
        cleaned.write(clean(rows[1]))
        cleaned.write("\n")