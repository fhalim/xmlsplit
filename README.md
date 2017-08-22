# Introduction
Throwaway utility to split out large XML file into chunks

# Building

```sh
go get github.com/fhalim/xmlsplit
```

# Using

```sh 
./xmlsplit -infile=WorkOrder.00.xml -outfileprefix=split-wo-00 -tagname=WorkOrder
```
