package main

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/genai"
)

const (
	model = "gemini-3.1-flash-lite"

	llmInstructions = `
	System: You are a clinical documentation assistant for a mobile crisis team. Your task is to read the provided crisis call transcript and extract the information into the exact JSON schema provided.

Instructions:

Output strictly in JSON matching the schema.

For boolean fields, map clinical denials to false and endorsements to true.

For the mental_health_symptoms_causing_impairment_summary, synthesize the caller's stated stressors, functional impairments (e.g., sleep, hygiene, motivation), and what they are seeking from the mobile crisis response into a concise, professional paragraph.

Do not infer information that is not explicitly stated in the transcript.

Extract the TEAMS Consult form data into strict JSON from the following transcript:
`
)

// TeamsConsultForm mirrors the JSON Schema we are asking the LLM to output
type TeamsConsultForm struct {
	VCCCaseNumber                      string `json:"vcc_case_number"`
	IsPICExperiencingSI                bool   `json:"is_pic_experiencing_si"`
	HasPICEngagedInSelfHarmToday       bool   `json:"has_pic_engaged_in_self_harm_today"`
	MentalHealthSymptomsSummary        string `json:"mental_health_symptoms_causing_impairment_summary"`
	NeedsImmediateIntervention         bool   `json:"needs_immediate_intervention"`
	HighLikelihoodHigherLevelOfCare    bool   `json:"high_likelihood_higher_level_of_care_needed"`
	UnderstandsAndWantsMCR             bool   `json:"understands_and_wants_mcr"`
	SubstanceUse                       bool   `json:"substance_use"`
	TypeOfSubstance                    string `json:"type_of_substance"`
	LastUseTime                        string `json:"last_use_time"`
	AmountUsed                         string `json:"amount_used"`
	ActiveWithdrawalSymptomsImpairment bool   `json:"active_withdrawal_symptoms_or_impairment"`
	ActivePsychosis                    bool   `json:"active_psychosis"`
}

// ExtractTeamsConsultForm sends a crisis call transcript to the model and parses
// the structured JSON response into a TeamsConsultForm. It returns the parsed
// form, the raw JSON text the model returned (handy for logging/debugging),
// and an error if the call or parse failed.
func ExtractTeamsConsultForm(ctx context.Context, client *genai.Client, transcriptText string) (*TeamsConsultForm, string, error) {
	prompt := fmt.Sprintf("%s %s", llmInstructions, transcriptText)

	config := &genai.GenerateContentConfig{
		Temperature:      genai.Ptr[float32](0.0), // Low temp for deterministic extraction
		ResponseMIMEType: "application/json",
		ResponseJsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"vcc_case_number":                                   map[string]any{"type": "string"},
				"is_pic_experiencing_si":                            map[string]any{"type": "boolean"},
				"has_pic_engaged_in_self_harm_today":                map[string]any{"type": "boolean"},
				"mental_health_symptoms_causing_impairment_summary": map[string]any{"type": "string"},
				"needs_immediate_intervention":                      map[string]any{"type": "boolean"},
				"high_likelihood_higher_level_of_care_needed":       map[string]any{"type": "boolean"},
				"understands_and_wants_mcr":                         map[string]any{"type": "boolean"},
				"substance_use":                                     map[string]any{"type": "boolean"},
				"type_of_substance":                                 map[string]any{"type": "string"},
				"last_use_time":                                     map[string]any{"type": "string"},
				"amount_used":                                       map[string]any{"type": "string"},
				"active_withdrawal_symptoms_or_impairment":          map[string]any{"type": "boolean"},
				"active_psychosis":                                  map[string]any{"type": "boolean"},
			},
			"required": []any{
				"vcc_case_number", "is_pic_experiencing_si",
				"has_pic_engaged_in_self_harm_today", "mental_health_symptoms_causing_impairment_summary",
				"needs_immediate_intervention", "high_likelihood_higher_level_of_care_needed",
				"understands_and_wants_mcr", "substance_use", "active_withdrawal_symptoms_or_impairment",
				"active_psychosis",
			},
		},
	}

	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, model, contents, config)
	if err != nil {
		return nil, "", fmt.Errorf("generate content: %w", err)
	}

	respText := resp.Text()

	var form TeamsConsultForm
	if err := json.Unmarshal([]byte(respText), &form); err != nil {
		return nil, respText, fmt.Errorf("parse JSON response: %w (raw: %s)", err, respText)
	}

	return &form, respText, nil
}