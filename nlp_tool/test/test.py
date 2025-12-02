import os, sys, csv

nlp_tool = os.path.join('..')
sys.path.insert(0, nlp_tool)

from naturalsyon import FilipinoCFGParser

parser = FilipinoCFGParser()

with open('data.csv', 'r') as f:
    reader = csv.reader(f)
    with open('res.txt', 'w') as out:
        ts = 0
        fs = 0
        for row in reader:
            inp, exp = row
            inp = inp.lower()
            exp = exp.lower()
            act = parser.apply_phonological_rules(inp)

            if exp == act:
                ts+=1
            else:
                fs+=1

            print(f"{inp:<20} -> {exp:<20} ? {act:<20} : {exp==act}")
            out.write(f"{inp:<20} -> {exp:<20} ? {act:<20} : {exp==act}")
            out.write("\n")


        print(f"Trues: {ts:>10}")
        print(f"Falses: {fs:>10}")
        out.write(f"Trues: {ts:>10}")
        out.write("\n")
        out.write(f"Falses: {fs:>10}")