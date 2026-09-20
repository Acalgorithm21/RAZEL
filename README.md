# autoRTD

## Real-Time Audio Transcription & Structured Data Extraction

autoRTD is a real-time audio processing prototype that converts spoken conversations into structured information.

The project combines live speech-to-text transcription with modular parsing logic. Instead of requiring someone to manually listen to a conversation and document important information afterward, autoRTD is designed to identify relevant information as the conversation happens and convert it into structured data.

> **Audio → Transcript → Parser → Structured Data**

---

## The Problem

Important information is often communicated through conversations but still has to be manually documented.

This can create:

* Time-consuming documentation workflows
* Inconsistent data entry
* Missed information
* Repetitive manual work
* Unstructured information that is difficult for software to process

autoRTD explores a way to automate part of this workflow by turning live speech into machine-readable information.

---

## How It Works

The MVP uses a simple processing pipeline:

```text
🎤 Microphone
     │
     ▼
Audio Capture
     │
     ▼
Deepgram Nova-3
     │
     ▼
Live Transcript
     │
     ▼
Specialized Parser
     │
     ▼
Structured Data
```

The transcription layer and parsing layer are intentionally separated.

This allows the same live transcription system to support multiple specialized parsers.

---

## Current MVP

The current MVP focuses on proving the core technical concept:

1. Capture audio from a microphone.
2. Stream the audio to Deepgram.
3. Receive transcription results in real time.
4. Process the transcript with a specialized parser.
5. Return structured information.

The MVP does **not require a graphical user interface**.

The prototype can demonstrate the complete workflow directly through the terminal.

Example:

```text
🎤 Listening...

Transcript:
Caller reports SI with a plan to overdose on pills
but no access to means and denied SH.

Extracted:

{
    "si_status": "yes",
    "plan": "yes",
    "access_to_means": "no",
    "sh_status": "denied"
}
```

---

## Parser Architecture

autoRTD is designed around **specialized parsers** rather than one large parser attempting to understand every possible piece of information.

For example:

```text
                    Transcript
                         │
          ┌──────────────┼──────────────┐
          ▼              ▼              ▼
      SI Parser    Medication Parser   ...
          │              │
          ▼              ▼
     SI Fields      Medication Data
```

Each parser is responsible for a specific category of information and produces predictable structured output.

This architecture makes it possible to add new extraction capabilities without redesigning the entire transcription system.

---

## Example: SI Parser

One planned parser focuses on identifying:

* SI status
* Plan status
* Access to means
* SH status

The parser follows a consistent output structure:

```json
{
  "si_status": "yes",
  "plan": "yes",
  "access_to_means": "no",
  "sh_status": "denied"
}
```

The goal is not simply to produce another transcript. The goal is to transform relevant information from the transcript into data that another software system can use.

---

## Technology

### Language

* Python

### Speech-to-Text

* Deepgram API
* Nova-3
* Real-time streaming transcription

### Audio

* PyAudio
* PortAudio

### Transcription Parsing

* golang
* Gemini gemini-3.1-flash-lite

### Configuration

* python-dotenv
* Environment variables

---

## Project Structure

```text
ATVproject/
│
├── main.py
├── README.md
├── requirements.txt
├── .env.example
├── .gitignore
│
└── ...
```

The project intentionally keeps API credentials outside the repository.

---

## Setup

### Requirements

* Python 3.x
* Microphone
* Deepgram API key
* macOS/Linux/Windows environment capable of running PyAudio

### 1. Clone the repository

```bash
git clone <repository-url>
cd ATVproject
```

### 2. Create a virtual environment

```bash
python3 -m venv .venv
```

Activate it:

**macOS/Linux**

```bash
source .venv/bin/activate
```

**Windows**

```bash
.venv\Scripts\activate
```

### 3. Install dependencies

go v1.27

```bash
pip install -r requirements.txt
```

### 4. Configure your API key

Copy the environment template:

```bash
cp .env.example .env
```

Then add your Datagram and Google AI Studio API keys:

```env
DEEPGRAM_API_KEY=your_datagram_api_key_here
GOOGLE_API_KEY=your_google_api_key_here
```

The `.env` file is intentionally excluded from Git.

### 5. Run the prototype

```bash
python main.py
```

The application will begin listening through the microphone and stream audio for transcription.

---

## Environment Variables

The application currently requires:

```env
DEEPGRAM_API_KEY=your_api_key_here
GOOGLE_API_KEY=your_google_api_key_here
```

Never commit the actual API key to the repository.

---

## Roadmap

### Phase 1 — Real-Time Transcription

* [x] Python environment
* [x] Deepgram SDK integration
* [x] Microphone audio capture
* [ ] Validate live streaming transcription
* [ ] Improve connection/error handling

### Phase 2 — Structured Extraction

* [ ] Build first production parser
* [ ] Define parser input/output contracts
* [ ] Add parser tests
* [ ] Connect transcript events to parsers
* [ ] Support multiple specialized parsers

### Phase 3 — Application Layer

* [ ] API layer
* [ ] Structured data storage
* [ ] Real-time data updates
* [ ] Optional web interface

### Phase 4 — Production

* [ ] Authentication
* [ ] Security controls
* [ ] Logging and monitoring
* [ ] Automated testing
* [ ] Deployment
* [ ] Performance optimization

---

## MVP Goal

The goal of the MVP is to demonstrate that spoken information can be transformed into structured data in real time.

The core concept is:

```text
Audio
  ↓
Speech Recognition
  ↓
Transcript
  ↓
Specialized Parser
  ↓
Structured Data
```

A graphical interface is intentionally outside the initial MVP scope. The priority is proving the underlying technical pipeline before building additional application layers.

---

## Future Vision

autoRTD is designed to become a modular real-time information extraction platform.

Once the transcription and parser architecture are established, additional parsers can be developed for different types of conversations and workflows.

The long-term architecture could support:

```text
                    Live Audio
                        │
                        ▼
                Speech-to-Text
                        │
                        ▼
                 Transcript
                        │
          ┌─────────────┼─────────────┐
          ▼             ▼             ▼
       Parser A      Parser B      Parser C
          │             │             │
          └─────────────┼─────────────┘
                        ▼
                Structured Data
                        │
                        ▼
                 Applications
```

The purpose of the MVP is to prove the foundation of that system.

---

## Project Status

🚧 **Competition MVP — In Development**

The project is currently focused on establishing and demonstrating the real-time transcription and structured data extraction pipeline.
