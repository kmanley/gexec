#!/usr/bin/env python
import signal
import sys
import time
def signal_handler(signal, frame):
    sys.stdout.write("got SIGINT; exiting...\n")
    sys.stdout.flush()
    sys.exit(42)
signal.signal(signal.SIGINT, signal_handler)
while True:
	time.sleep(1.0)