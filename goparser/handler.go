package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/genai"
)

const (
	incomingDir  = "../transcripts/incoming"  // Python writes finished transcripts here (.txt)
	processedDir = "../transcripts/processed" // extraction results land here (.json)
	doneDir      = "../transcripts/done"      // original transcripts moved here once processed
	pollInterval = 2 * time.Second
)

func main() {
	ctx := context.Background()

	// Passing nil tells the SDK to automatically use the GOOGLE_API_KEY environment variable.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create GenAI client: %v", err)
	}

	for _, dir := range []string{incomingDir, processedDir, doneDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	log.Printf("Watching %s for finished transcripts...", incomingDir)

	for {
		if err := processNewTranscripts(ctx, client); err != nil {
			log.Printf("Error scanning for transcripts: %v", err)
		}
		time.Sleep(pollInterval)
	}
}

// processNewTranscripts looks for .txt files in incomingDir and hands each one
// to handleTranscriptFile. Errors on one file are logged and skipped so a single
// bad transcript doesn't block the rest of the queue.
func processNewTranscripts(ctx context.Context, client *genai.Client) error {
	entries, err := os.ReadDir(incomingDir)
	if err != nil {
		return fmt.Errorf("read incoming dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		path := filepath.Join(incomingDir, entry.Name())
		if err := handleTranscriptFile(ctx, client, path, entry.Name()); err != nil {
			log.Printf("Failed to process %s: %v", entry.Name(), err)
			continue
		}
	}

	return nil
}

// handleTranscriptFile reads one transcript, runs it through the LLM parser,
// writes the structured result as JSON, then moves the original transcript to
// doneDir so it won't be picked up again on the next poll.
func handleTranscriptFile(ctx context.Context, client *genai.Client, path, name string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	form, _, err := ExtractTeamsConsultForm(ctx, client, string(data))
	if err != nil {
		return fmt.Errorf("extract: %w", err)
	}

	// ADD THIS: print the extracted form to the terminal
	fmt.Println("✅ Extraction Successful!")
	fmt.Printf("Case Number: %s\n", form.VCCCaseNumber)
	fmt.Printf("Is Experiencing SI: %v\n", form.IsPICExperiencingSI)
	fmt.Printf("Summary: %s\n", form.MentalHealthSymptomsSummary)
	fmt.Printf("Substance Use: %v (Type: %s, Last Use: %s, Amount: %s)\n",
		form.SubstanceUse, form.TypeOfSubstance, form.LastUseTime, form.AmountUsed)
	fmt.Printf("Needs Immediate Intervention: %v\n", form.NeedsImmediateIntervention)
	fmt.Printf("Wants Mobile Crisis Response: %v\n", form.UnderstandsAndWantsMCR)
	fmt.Println()

	prettyJSON, err := json.MarshalIndent(form, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	baseName := strings.TrimSuffix(name, ".txt")
	outPath := filepath.Join(processedDir, baseName+".json")
	if err := os.WriteFile(outPath, prettyJSON, 0644); err != nil {
		return fmt.Errorf("write result: %w", err)
	}

	donePath := filepath.Join(doneDir, name)
	if err := os.Rename(path, donePath); err != nil {
		return fmt.Errorf("move processed transcript: %w", err)
	}

	log.Printf("Processed %s -> %s (case: %s)", name, outPath, form.VCCCaseNumber)
	return nil
}