#!/usr/bin/env python
import signal
import sys
import time
def signal_handler(signal, frame):
    sys.stdout.write("got SIGINT; ignoring it...\n")
    sys.stdout.flush()
signal.signal(signal.SIGINT, signal_handler)
while True:
	time.sleep(1.0)