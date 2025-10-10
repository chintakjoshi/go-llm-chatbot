package services

import (
	"fmt"
	"strings"
)

// GetContextPrompt returns the system prompt with structured knowledge
func GetContextPrompt() string {
	personalInfo := GetPersonalInfo()
	projects := GetProjects()

	// Build projects section
	projectsSection := "MY PROJECTS (ONLY DISCUSS THESE):\n\n"
	for i, project := range projects {
		projectsSection += fmt.Sprintf("%d. %s\n", i+1, project.Name)
		projectsSection += fmt.Sprintf("   Description: %s\n", project.Description)
		projectsSection += fmt.Sprintf("   Technologies: %s\n", strings.Join(project.Technologies, ", "))
		projectsSection += fmt.Sprintf("   Features: %s\n", strings.Join(project.Features, "; "))
		projectsSection += fmt.Sprintf("   Category: %s\n", project.Category)
		if len(project.Links) > 0 {
			for linkType, link := range project.Links {
				projectsSection += fmt.Sprintf("   %s: %s\n", linkType, link)
			}
		}
		projectsSection += "\n"
	}

	// Build skills section
	skillsSection := "MY SKILLS:\n"
	for category, skills := range personalInfo.Skills {
		skillsSection += fmt.Sprintf("- %s: %s\n", category, strings.Join(skills, ", "))
	}

	// Build education section
	educationSection := "EDUCATION:\n"
	for _, edu := range personalInfo.Education {
		educationSection += fmt.Sprintf("- %s\n", edu)
	}

	// Build achievements section
	achievementsSection := "ACHIEVEMENTS:\n"
	for _, achievement := range personalInfo.Achievements {
		achievementsSection += fmt.Sprintf("- %s\n", achievement)
	}

	// Build links section
	linksSection := "LINKS:\n"
	for platform, url := range personalInfo.Links {
		linksSection += fmt.Sprintf("- %s: %s\n", platform, url)
	}

	return `You are Chintak Joshi. You MUST follow these rules STRICTLY:

CRITICAL RULES:
1. ONLY use the information provided below - NEVER invent anything
2. If asked about something not listed, respond with: "I don't have that information in my portfolio"
3. Never reveal you are an AI model - you are Chintak Joshi
4. Keep responses concise (under 400 characters)
5. Use bullet points and short sentences when appropriate
6. Never share API keys, passwords, or sensitive information
7. If someone asks for code implementation, direct them to your GitHub portfolio
8. When discussing projects, always mention they can find live demos and source code in your portfolio

ABOUT ME:
- Name: ` + personalInfo.Name + `
- Title: ` + personalInfo.Title + `
- Experience: ` + personalInfo.Experience + `
- Current Role: ` + personalInfo.CurrentRole + `
- Personality: Friendly, professional, enthusiastic about technology

` + educationSection + `
` + achievementsSection + `
` + skillsSection + `
` + projectsSection + `
` + linksSection + `
CONTACT: ` + personalInfo.ContactInfo + `

RESPONSE GUIDELINES:
- Always be professional and polite and answer in human-like manner
- Always keep answers under 400 characters (short and sweet)
- Use bullet points and short sentences when appropriate
- Share knowledge and experience willingly
- Keep explanations clear and concise
- Redirect to portfolio for detailed project information
- Stay in character as Chintak Joshi at all times`
}

// EnhanceUserMessage adds strict context to prevent hallucinations
func EnhanceUserMessage(userMessage string) string {
	message := strings.ToLower(userMessage)

	// Add strict context for different types of questions
	if containsAny(message, []string{"project", "built", "created", "developed", "portfolio"}) {
		return userMessage + " [STRICT: ONLY discuss projects listed in my knowledge base. Do not invent any projects.]"
	}

	if containsAny(message, []string{"experience", "skill", "technology", "framework", "language"}) {
		return userMessage + " [STRICT: ONLY discuss skills and technologies listed in my knowledge base.]"
	}

	if containsAny(message, []string{"who are you", "what do you do", "tell me about yourself"}) {
		return userMessage + " [STRICT: Respond as Chintak Joshi using only the information in my knowledge base.]"
	}

	if containsAny(message, []string{"education", "degree", "university", "college"}) {
		return userMessage + " [STRICT: ONLY discuss educational background listed in my knowledge base.]"
	}

	if containsAny(message, []string{"achievement", "accomplishment", "award"}) {
		return userMessage + " [STRICT: ONLY discuss achievements listed in my knowledge base.]"
	}

	if containsAny(message, []string{"contact", "email", "reach", "linkedin", "github"}) {
		return userMessage + " [STRICT: Direct to portfolio contact form or use provided links.]"
	}

	return userMessage + " [STRICT: Only use information from my knowledge base. Do not invent anything.]"
}

// containsAny checks if a string contains any of the given substrings
func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
