
EXECUTABLE = ttk4145demoelevator

CC = gcc
CFLAGS = -Wall -Wextra -g -std=gnu11
LDFLAGS = -lcomedi -lm -lpthread

SOURCES := $(patsubst ./%, %, $(shell find . -name '*.c'))
OBJECTS = $(addprefix $(OBJDIR)/, $(SOURCES:.c=.o))

DC = dmd
SIMSRC = driver/sim_backend.d driver/timer_event.d
SIMLIB = $(OBJDIR)/driver/simelev.a
PHOBOS = /usr/lib/x86_64-linux-gnu/libphobos2.a

OBJDIR = obj

all: $(EXECUTABLE)

rebuild:	clean all
	
clean:
	rm -f $(EXECUTABLE)
	rm -rf $(OBJDIR)

$(EXECUTABLE): $(OBJECTS) $(DRIVEROBJ) $(SIMLIB) $(PHOBOS)
	$(CC) $^ -o $@ $(LDFLAGS)
	
$(OBJDIR)/%.o: %.c
	@mkdir -p $(@D)
	$(CC) -o $@ -c $(CFLAGS) $<	
    
$(SIMLIB): $(SIMSRC)
	$(DC) -lib -of$@ $^
	
.PHONY: all rebuild clean

