CC=clang++

TARGET=toy

COMPILE_FLAGS=$(shell llvm-config --cxxflags)
LINK_FLAGS=$(shell llvm-config --ldflags  --system-libs --libs core)

$(TARGET): $(TARGET).cpp
	@$(CC) -g -c $(COMPILE_FLAGS) $(TARGET).cpp -o $(TARGET).o
	@$(CC) $(TARGET).o $(LINK_FLAGS) -lpthread -o $(TARGET)

.PHONY:clean
clean:
	@rm -rf *.out
	@rm -rf *.o
	@rm -rf $(TARGET)
