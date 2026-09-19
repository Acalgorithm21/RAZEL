import os
import time
import threading
import pyaudio

from dotenv import load_dotenv
from deepgram import DeepgramClient
from deepgram.core.events import EventType

load_dotenv()

API_KEY = os.getenv("DEEPGRAM_API_KEY")

if not API_KEY:
    raise RuntimeError("DEEPGRAM_API_KEY is missing from .env")


def main():
    client = DeepgramClient(api_key=API_KEY)

    with client.listen.v1.connect(
        model="nova-3",
        language="en-US",
        smart_format="true",
        encoding="linear16",
        channels="1",
        sample_rate="16000",
        interim_results="true",
        endpointing="150",
        utterance_end_ms="1000",
        diarize="true",
    ) as connection:

        def on_message(message) -> None:
            msg_type = getattr(message, "type", None)

            if msg_type != "Results":
                return

            transcript = message.channel.alternatives[0].transcript

            if message.is_final:
                if transcript:
                    print(f"\r{transcript}{' ' * 20}")  # lock in the finished line
                else:
                    print(f"\r[inaudible]{' ' * 20}")
            else:
                if transcript:
                    print(f"\r{transcript}", end="", flush=True)  # live, updating in pla


        def on_error(error) -> None:
            print(f"\n[Error] {error}\n")

        connection.on(EventType.OPEN, lambda _: print("Connection opened"))
        connection.on(EventType.MESSAGE, on_message)
        connection.on(EventType.CLOSE, lambda _: print("Connection closed"))
        connection.on(EventType.ERROR, on_error)

        # Run start_listening() on a background thread so it doesn't block us
        listener_thread = threading.Thread(
            target=connection.start_listening,
            daemon=True,
        )
        listener_thread.start()

        audio = pyaudio.PyAudio()
        stream = audio.open(
            format=pyaudio.paInt16,
            channels=1,
            rate=16000,
            input=True,
            frames_per_buffer=800,
        )

        print("🎤 Listening...")
        print("Speak into your microphone.")
        print("Press Ctrl+C to stop.")

        try:
            while True:
                data = stream.read(800, exception_on_overflow=False)
                connection.send_media(data)
                time.sleep(0)
        except KeyboardInterrupt:
            print("\nStopping...")
        finally:
            stream.stop_stream()
            stream.close()
            audio.terminate()
            try:
                connection.send_close_stream()
            except Exception:
                pass


if __name__ == "__main__":
    main()