package services

import (
	"strings"
)

// ValidateResponse checks if the response contains known information
func ValidateResponse(response string, userMessage string) (bool, string) {
	personalInfo := GetPersonalInfo()
	projects := GetProjects()

	responseLower := strings.ToLower(response)
	userMessageLower := strings.ToLower(userMessage)

	// Check for common hallucination patterns
	forbiddenPhrases := []string{
		"as an ai",
		"as a language model",
		"i cannot",
		"i'm sorry, i",
		"openai",
		"chatgpt",
		"i am an ai",
		"as an artificial intelligence",
	}

	for _, phrase := range forbiddenPhrases {
		if strings.Contains(responseLower, phrase) {
			return false, "Response contains AI identity revelation"
		}
	}

	// Check for known projects
	knownProjectNames := make(map[string]bool)
	for _, project := range projects {
		knownProjectNames[strings.ToLower(project.Name)] = true
		// Also check for partial matches of project names
		words := strings.Fields(strings.ToLower(project.Name))
		for _, word := range words {
			if len(word) > 4 { // Only consider words longer than 4 characters
				knownProjectNames[word] = true
			}
		}
	}

	// Check for known skills and technologies
	knownSkills := make(map[string]bool)
	for _, skills := range personalInfo.Skills {
		for _, skill := range skills {
			knownSkills[strings.ToLower(skill)] = true
		}
	}

	// Validation based on message type
	if strings.Contains(userMessageLower, "project") {
		// Count how many known projects are mentioned
		mentionedKnown := 0
		for projectName := range knownProjectNames {
			if strings.Contains(responseLower, projectName) {
				mentionedKnown++
			}
		}

		// If talking about projects but not mentioning any known ones, might be hallucinating
		if mentionedKnown == 0 && len(response) > 100 {
			return false, "Response about projects but no known projects mentioned"
		}
	}

	// Check if response mentions unknown technologies when asked about skills
	if strings.Contains(userMessageLower, "skill") || strings.Contains(userMessageLower, "technology") || strings.Contains(userMessageLower, "framework") {
		// Simple check: if the response contains many technology words but few known ones
		techWords := []string{"react", "node", "python", "java", "docker", "aws", "machine learning", "ai", "ml", "database", "server"}
		mentionedKnownTech := 0
		totalTechMentions := 0

		// Count known technologies mentioned
		for tech := range knownSkills {
			if strings.Contains(responseLower, tech) {
				mentionedKnownTech++
			}
		}

		// Count total technology-related words mentioned
		for _, techWord := range techWords {
			if strings.Contains(responseLower, techWord) {
				totalTechMentions++
			}
		}

		// If talking about technologies but very few known ones are mentioned relative to total tech words, might be hallucinating
		if mentionedKnownTech < 2 && totalTechMentions > 3 && len(response) > 150 {
			return false, "Response about skills but very few known technologies mentioned compared to total tech words"
		}
	}

	// Check for educational background questions
	if strings.Contains(userMessageLower, "education") || strings.Contains(userMessageLower, "degree") || strings.Contains(userMessageLower, "university") {
		hasEducationInfo := strings.Contains(responseLower, "saint louis") || strings.Contains(responseLower, "gujarat") || strings.Contains(responseLower, "computer science")
		if !hasEducationInfo && len(response) > 100 {
			return false, "Response about education but no known educational background mentioned"
		}
	}

	return true, ""
}

// SimpleResponseValidation is a lighter version for production use
func SimpleResponseValidation(response string) bool {
	responseLower := strings.ToLower(response)

	// Quick checks for obvious AI identity leaks
	forbiddenPhrases := []string{
		"as an ai",
		"as a language model",
		"i cannot",
		"openai",
		"chatgpt",
	}

	for _, phrase := range forbiddenPhrases {
		if strings.Contains(responseLower, phrase) {
			return false
		}
	}

	return true
}
