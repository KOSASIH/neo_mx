# streaming/__init__.py
from kafka import KafkaProducer
from flink import Flink

class NeoMXStreamer:
    def __init__(self):
        self.producer = KafkaProducer(bootstrap_servers='localhost:9092')
        self.flink = Flink()

    def stream_data(self, data):
        self.producer.send('neo_mx_topic', value=data)
        self.flink.add_sink(self.producer)

    def start_streaming(self):
        self.flink.start()
