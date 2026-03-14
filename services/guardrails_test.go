package services

import (
	"strings"
	"testing"
)

func TestApplyGuardrailsBlocksOutOfScopePrompts(t *testing.T) {
	t.Parallel()

	tests := []string{
		"what is 2+3?",
		"write a python function to print sum of 2 numbers",
		"can you write a python function?",
		"what is Go?",
	}

	for _, input := range tests {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			decision := ApplyGuardrails(input)
			if decision.DirectResponse != PortfolioOnlyResponse {
				t.Fatalf("expected out-of-scope response for %q, got %+v", input, decision)
			}
		})
	}
}

func TestApplyGuardrailsAllowsPortfolioQueries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "projects",
			input:  "What projects have you built?",
			expect: "PROJECTS:",
		},
		{
			name:   "skills",
			input:  "Tell me about your React experience",
			expect: "SKILLS:",
		},
		{
			name:   "education",
			input:  "Where did you study?",
			expect: "EDUCATION:",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			decision := ApplyGuardrails(test.input)
			if decision.DirectResponse != "" {
				t.Fatalf("expected LLM path for %q, got direct response %q", test.input, decision.DirectResponse)
			}
			if !strings.Contains(decision.SystemPrompt, test.expect) {
				t.Fatalf("expected system prompt to include %q, got %q", test.expect, decision.SystemPrompt)
			}
		})
	}
}

func TestApplyGuardrailsAllowsFollowUpWithinSession(t *testing.T) {
	t.Parallel()

	sessionID := "session-follow-up"

	firstDecision := ApplyGuardrailsWithSession(sessionID, "What's your favorite project?")
	if firstDecision.DirectResponse != "" {
		t.Fatalf("expected first portfolio query to use LLM path, got %+v", firstDecision)
	}

	followUpDecision := ApplyGuardrailsWithSession(sessionID, "Wow, can you tell me in detail?")
	if followUpDecision.DirectResponse != "" {
		t.Fatalf("expected follow-up query to use LLM path, got %+v", followUpDecision)
	}
	if !strings.Contains(followUpDecision.UserPrompt, "follow-up") {
		t.Fatalf("expected follow-up prompt instructions, got %q", followUpDecision.UserPrompt)
	}
	if !strings.Contains(followUpDecision.SystemPrompt, "PROJECTS:") {
		t.Fatalf("expected project context in follow-up prompt, got %q", followUpDecision.SystemPrompt)
	}
}

func TestApplyGuardrailsHandlesGreetingWithoutLLM(t *testing.T) {
	t.Parallel()

	decision := ApplyGuardrails("hello")
	if decision.DirectResponse != GreetingResponse {
		t.Fatalf("expected greeting response, got %+v", decision)
	}
}

func TestApplyGuardrailsHandlesAcknowledgmentsWithoutLLM(t *testing.T) {
	t.Parallel()

	tests := []string{
		"looks good",
		"impressive",
		"thanks",
		"ok",
		"great work",
		"well done",
		"nice",
		"good job",
		"awesome",
		"fantastic",
		"amazing",
		"excellent",
		"wonderful",
		"superb",
		"brilliant",
		"outstanding",
		"fine",
		"acceptable",
	}

	for _, input := range tests {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			decision := ApplyGuardrails(input)
			if decision.DirectResponse != AcknowledgmentResponse {
				t.Fatalf("expected acknowledgment response for %q, got %+v", input, decision)
			}
		})
	}
}
