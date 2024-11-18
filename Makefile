build : 
		@go build -o  ./bin/fluentVal
run : build 
	    @./bin/fluentVal
