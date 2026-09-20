package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/genai"
)

const (
	model           = "gemini-3.1-flash-lite"
	llmInstructions = `
	System: You are a clinical documentation assistant for a mobile crisis team. Your task is to read the provided crisis call transcript and extract the information into the exact JSON schema provided.

Instructions:

Output strictly in JSON matching the schema.

For boolean fields, map clinical denials to false and endorsements to true.

For the mental_health_symptoms_causing_impairment_summary, synthesize the caller's stated stressors, functional impairments (e.g., sleep, hygiene, motivation), and what they are seeking from the mobile crisis response into a concise, professional paragraph.

Do not infer information that is not explicitly stated in the transcript.

Extract the TEAMS Consult form data into strict JSON from the following transcript: 
`

	transcript = `
	CW: Hey, I have a caller on hold — case code [FAKE-1234]. Male, 40s, laid off last week, endorsing passive SI, no plan. He's asking about a mobile crisis response. Bringing him on.
MCS: Hi, this is the mobile crisis specialist. We'll confirm some of your information and talk a little more about what's been going on. The crisis worker mentioned the layoff has been hitting you hard — I hate that you're going through that.
Caller: Yeah. Fifteen years at the same warehouse and it just ended Friday.
MCS: That's a huge loss. Is there anything in addition going on that's making things more stressful?
Caller: Rent's due in two weeks. That's most of it.
MCS: And the crisis worker noted some thoughts of not wanting to be here — can you tell me more about those?
Caller: It's like "what's the point" popping in. Every night since Friday. No plan or anything, I wouldn't do it.
MCS: I'm glad you told me. How has all this been impacting your day-to-day functioning?
Caller: I don't know. Stressed, mostly.
MCS: Sometimes when people are going through something like this, it hits sleep, appetite, motivation — anything like that for you?
Caller: Sleep's shot — three hours maybe. Haven't left the apartment since Saturday. I can't even start the unemployment paperwork.
MCS: Okay. How do you feel a mobile crisis response could help you tonight?
Caller: I just don't want to sit alone with my head like this. Talking to someone in person... yeah.
MCS: Alright. Let me verify some information. Can you spell your first and last name for me?
Caller: [Spells fictional name.]
MCS: I have an address here — [fictional address]. It's showing me an apartment complex on Google, is that right?
Caller: Yeah, unit B.
MCS: Date of birth? And just so the team can identify you, your race?
Caller: [Fictional DOB.] White.
MCS: Thank you. Give me one moment, I'd like to do some documentation. (pause) Okay — I'm going to loop in one of our clinicians to review your case, so I have a few more questions. Are you under the influence of any drugs or alcohol right now?
Caller: No. I've been drinking a little more at night — few beers — but nothing today.
MCS: When was the last time, and how much?
Caller: Last night, four beers.
MCS: Any withdrawal symptoms — shakiness, sweating, nausea when you don't drink?
Caller: No.
MCS: And I ask everyone this — any experiences like hearing or seeing things others don't, or feeling paranoid, like people are watching or after you?
Caller: No, nothing like that.
MCS: Okay. Give me a few minutes to pull your information together for the clinician. I'll place you on a brief hold.
(hold)
MCS: Are you still there?
Caller: ...Yeah, still here.
MCS: Thanks for your patience. Now let's go through some safety questions while we wait. Any legal history or history of violence?
Caller: No.
MCS: Any mental health diagnoses?
Caller: Depression, back in my twenties. Nothing current.
MCS: Any weapons in the home?
Caller: No guns. Just my regular blood pressure medication.
MCS: Any illness we should know about — COVID, anything contagious?
Caller: No.
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

func main() {
	ctx := context.Background()

	// Initialize the client.
	// Passing nil tells the SDK to automatically use the GOOGLE_API_KEY environment variable.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create GenAI client: %v", err)
	}

	// This is where you would load your caller transcript
	transcriptText := transcript

	prompt := fmt.Sprintf("%s %s", llmInstructions, transcriptText)

	// Define the exact JSON schema via GenerateContentConfig
	// Setting ResponseMIMEType and ResponseJsonSchema forces the model into Structured Output mode.
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

	// 3. Construct the message content payload
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}

	// Call the GenerateContent API
	resp, err := client.Models.GenerateContent(ctx, model, contents, config)
	if err != nil {
		log.Fatalf("Error generating content: %v", err)
	}

	// Extract the generated text and unmarshal it into the Go struct
	var form TeamsConsultForm
	respText := resp.Text()

	if err := json.Unmarshal([]byte(respText), &form); err != nil {
		log.Fatalf("Failed to parse JSON response: %v\nRaw response: %s", err, respText)
	}

	// 6. Output the successful parse
	fmt.Println("✅ Extraction Successful!")
	fmt.Printf("Case Number: %s\n", form.VCCCaseNumber)
	fmt.Printf("Summary: %s\n", form.MentalHealthSymptomsSummary)
	fmt.Printf("Substances Endorsed: %v (Type: %s, Amount: %s)\n",
		form.SubstanceUse, form.TypeOfSubstance, form.AmountUsed)
	fmt.Println()
	fmt.Printf("%s\n", respText)
	fmt.Println()
	fmt.Printf("%+v\n", form)
}

func listModels(ctx context.Context, client *genai.Client) {
	modelIter := client.Models.All(ctx)
	for model := range modelIter {
		fmt.Println(model)
		fmt.Println()
	}

	if ctx != nil {
		return
	}
}
