#!/usr/bin/env python

import sys

# This file was called GobCodeGen, but we found some un-deterministic spots in
# gob, so we abandon it. As you can see, now it's a amigo(cdc)
# marshal/unmarshal warpper.

template = """
func decode$TYPE(data []byte) *TYPE {
	if data == nil {
		return nil
	}
	rst := &TYPE{}
	cdc.MustUnmarshalBinaryBare(data, rst)
	return rst
}

func encode$TYPE(dt *TYPE) []byte {
	if dt == nil {
		return nil
	}
	rst := cdc.MustMarshalBinaryBare(dt)
	return []byte(rst)
}
"""

structs = ['userMeta', 'roundMeta', 'roundPostMeta', 'gameMeta']


def cap(word):
    return word[0].upper() + word[1:]


def main():
    if len(sys.argv) < 2:
        for type_name in structs:
            tmp = template.replace("$TYPE", cap(type_name))
            print tmp.replace("TYPE", type_name)
    else:
        type_name = sys.argv[1]
        tmp = template.replace("$TYPE", cap(type_name))
        print tmp.replace("TYPE", type_name)


if __name__ == '__main__':
    main()
