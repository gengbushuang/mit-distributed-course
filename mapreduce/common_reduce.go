package mapreduce

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTask int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	fmt.Printf("doReduce action jobName=%v reduceTask=%v outFile=%v nMap=%v\n", jobName, reduceTask, outFile, nMap)

	mapKV := make(map[string][]string)

	for i := 0; i < nMap; i++ {
		mapFile := reduceName(jobName, i, reduceTask)
		fmt.Printf("doReduce mapFile=%v\n", mapFile)

		file, _ := os.Open(mapFile)

		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			err := dec.Decode(&kv)
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("reduce read file:%s can't open error:%v\n", mapFile, err)
				break
			}

			_, ok := mapKV[kv.Key]
			if !ok {
				mapKV[kv.Key] = make([]string, 1)

			}
			mapKV[kv.Key] = append(mapKV[kv.Key], kv.Value)
		}

		defer file.Close()

	}

	var keys []string
	for k, _ := range mapKV {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	file, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("reduce merge file:%s can't open error:%v\n", outFile, err)
		return
	}

	enc := json.NewEncoder(file)

	for _, k := range keys {
		enc.Encode(KeyValue{k, reduceF(k, mapKV[k])})
	}

	file.Close()

	//
	// doReduce manages one reduce task: it should read the intermediate
	// files for the task, sort the intermediate key/value pairs by key,
	// call the user-defined reduce function (reduceF) for each key, and
	// write reduceF's output to disk.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTask) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Your code here (Part I).
	//
}
