# streaming/stream.py
import time
from streaming import NeoMXStreamer

streamer = NeoMXStreamer()

while True:
    data = ...  # Generate or fetch data
    streamer.stream_data(data)
    time.sleep(1)
