#!/usr/bin/env python
import sys
template = """
func decode$TYPE(data []byte) *TYPE {
	if data == nil {
		return nil
	}
	rst := &TYPE{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(rst)
	if err != nil {
		panic("error in gob decode TYPE" + err.Error())
	}
	return rst
}

func encode$TYPE(dt *TYPE) []byte {
	if dt == nil {
		return nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(dt)
	if err != nil {
		panic("error in encoding: " + err.Error())
	}
	return buf.Bytes()
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
