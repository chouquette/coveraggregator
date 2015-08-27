# coveraggregator
Cover profile aggregator for golang

## What is this? ##
This is a really simple tool that aims to work around a current go test 
limitation, preventing you from computing a code coverage profile on multiple packages.

It is often suggested to simply append files one after another to achieve this, 
however the go tool cover utility won't aggregate results, which will make
your code coverage results invalid.

This tool is simply iterating over the results stored in the files you provide
to it, and sums things up, so that if a function is tested by a package, but 
not by another one, it will appear as "tested" in the aggregated result.

## How to use it? ##
Run the tool and provide these:
- a "-o <output file>" parameter, which is where to write the aggregated profile
- Individual cover profiles, which will be aggregated into one

## What's next? ##
- It would probably be cool to extend this tool to have it run the tests by invoking
go test itself, rather than post-process files.
- Whatever you can think of :)
