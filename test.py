import pyaudio

print("Opening PyAudio...")
audio = pyaudio.PyAudio()
print("PyAudio initialized")

stream = audio.open(
    format=pyaudio.paInt16,
    channels=1,
    rate=16000,
    input=True,
    frames_per_buffer=800,
)
print("Stream opened — reading 5 chunks")

for i in range(5):
    data = stream.read(800, exception_on_overflow=False)
    print(f"Chunk {i}: {len(data)} bytes")

stream.stop_stream()
stream.close()
audio.terminate()
print("Done")