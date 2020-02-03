cc = g++
prom = helloworld.out
deps = $(shell find ./include -name "*.h")
src = $(shell find ./src -name "*.cpp")
obj = $(src:%.c=%.o) 

$(prom): $(obj)
	$(cc) -o $(prom) $(obj)

clean:
	rm -rf $(prom)