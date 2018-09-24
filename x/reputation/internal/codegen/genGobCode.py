#!/usr/bin/env python

import sys

# Though this file is called GobCodeGen, but we found some un-deterministic spots in
# gob, so we abandon it. As you can see, now it's a json marshal/unmarshal warpper.

template = """
func decode$TYPE(data []byte) *TYPE {
	if data == nil {
		return nil
	}
	rst := &TYPE{}
        err := json.Unmarshal(data, &rst)
	if err != nil {
		panic("error in json decode TYPE" + err.Error())
	}
	return rst
}

func encode$TYPE(dt *TYPE) []byte {
	if dt == nil {
		return nil
	}
        rst, err := json.Marshal(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return []byte(rst)
}
"""

def cap(word):
    return word[0].upper() + word[1:];

def main():
    if len(sys.argv) < 2:
        print "missing type name"
    else:
        type_name = sys.argv[1]
        tmp = template.replace("$TYPE", cap(type_name))
        print tmp.replace("TYPE", type_name)

if __name__ == '__main__':
    main()
