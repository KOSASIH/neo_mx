# neural_networks/__init__.py
from tensorflow.keras.models import Model
from tensorflow.keras.layers import Input, Dense, LSTM, Conv1D

class NeoMXModel(Model):
    def __init__(self):
        super(NeoMXModel, self).__init__()
        self.input_layer = Input(shape=(100, 10))
        self.lstm_layer = LSTM(64, return_sequences=True)
        self.conv_layer = Conv1D(32, kernel_size=3, activation='relu')
        self.dense_layer = Dense(10, activation='softmax')
        self.output_layer = self.dense_layer(self.conv_layer(self.lstm_layer(self.input_layer)))

    def call(self, inputs):
        return self.output_layer
